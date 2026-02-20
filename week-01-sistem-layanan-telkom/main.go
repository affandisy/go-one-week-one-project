package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Ticket struct {
	TicketID      string    `json:"ticket_id" gorm:"primaryKey;column:ticket_id"`
	CustomerID    string    `json:"customer_id" gorm:"column:customer_id"`
	Issue         string    `json:"issue" gorm:"column:issue"`
	Priority      string    `json:"priority" gorm:"column:priority"`
	Status        string    `json:"status" gorm:"column:status;default:open"`
	AssignedAgent string    `json:"assigned_agent" gorm:"column:assigned_agent"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

var (
	db          *gorm.DB
	rmqConn     *amqp.Connection
	rmqChannel  *amqp.Channel
	redisClient *redis.Client
	ctx         = context.Background()
)

const QueueName = "customer.tickets"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default env variables")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	log.Println("Database connection successful!")

	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	rmqConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}
	defer rmqConn.Close()

	rmqChannel, err = rmqConn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel: ", err)
	}
	defer rmqChannel.Close()

	_, err = rmqChannel.QueueDeclare(QueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue: ", err)
	}

	app := fiber.New()

	app.Post("/tickets", CreateTicket)
	app.Get("/tickets", GetTicket)
	app.Put("/ticket/:id", UpdateTicketStatus)

	log.Println("Server is running on port 3000")
	app.Listen(":3000")
}

func CreateTicketProducer(c *fiber.Ctx) error {
	var ticket Ticket
	if err := c.BodyParser(&ticket); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	ticket.TicketID = "TKT-" + uuid.New().String()[:8]
	ticket.CreatedAt = time.Now()
	ticket.Status = "open"

	body, _ := json.Marshal(ticket)

	err := rmqChannel.PublishWithContext(ctx, "", QueueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to publish ticket"})
	}

	return c.Status(201).JSON(fiber.Map{"ticket_id": ticket.TicketID})
}

func StartTicketConsumer() {
	msgs, err := rmqChannel.Consume(QueueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register a consumer: ", err)
	}

	log.Println("[*] Waiting for tickets in RabbitMQ. To exit press CTRL+C")
	for msg := range msgs {
		var ticket Ticket
		json.Unmarshal(msg.Body, &ticket)

		if err := db.Create(&ticket).Error; err != nil {
			log.Printf("Error saving ticket %s to DB: %v \n", ticket.TicketID, err)
		} else {
			log.Printf("Successfully processed and saved ticket: %s\n", ticket.TicketID)
		}
	}
}

func GetTicketsCached(c *fiber.Ctx) error {
	agentID := c.Query("agent_id")
	status := c.Query("status")

	if agentID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "agent_id is required"})
	}

	cacheKey := fmt.Sprintf("tickets:active:%s", agentID)

	cachedTickets, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		log.Println("Cache hit! Fetching from Redis")
		var tickets []Ticket
		json.Unmarshal([]byte(cachedTickets), &tickets)
		return c.JSON(tickets)
	}

	log.Println("Cache miss! Fetching from DB")
	var tickets []Ticket
	query := db.Where("assigned_agent = ?", agentID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tickets).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database Error"})
	}

	ticketJSON, _ := json.Marshal(tickets)
	redisClient.Set(ctx, cacheKey, ticketJSON, 5*time.Minute)

	return c.JSON(tickets)
}

func CreateTicket(c *fiber.Ctx) error {
	var ticket Ticket
	if err := c.BodyParser(ticket); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := db.Create(&ticket).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save ticket"})
	}

	return c.Status(201).JSON(fiber.Map{"ticket_id": ticket.TicketID})
}

func GetTicket(c *fiber.Ctx) error {
	agentID := c.Query("agent_id")
	status := c.Query("status")

	var tickets []Ticket

	query := db.Where("assigned_agent = ?", agentID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tickets).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	return c.JSON(tickets)
}

func UpdateTicketStatus(c *fiber.Ctx) error {
	ticketID := c.Params("id")

	var reqBody struct {
		Status        string `json:"status"`
		AssignedAgent string `json:"assigned_agent"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	updateData := map[string]interface{}{"status": reqBody.Status}
	if reqBody.AssignedAgent != "" {
		updateData["assigned_agent"] = reqBody.AssignedAgent
	}

	result := db.Model(&Ticket{}).Where("ticket_id = ?", ticketID).Updates(updateData)
	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update ticket"})
	}

	if reqBody.AssignedAgent != "" {
		cacheKey := fmt.Sprintf("tickets:active:%s", reqBody.AssignedAgent)
		redisClient.Del(ctx, cacheKey)
	}

	return c.JSON(fiber.Map{"message": "Ticket status updated successfully"})
}
