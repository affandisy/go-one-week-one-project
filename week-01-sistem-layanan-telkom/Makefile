migrate create -ext sql -dir db/migrations -seq create_tickets_table

migrate -path db/migrations -database "postgresql://telkom_user:telkom_password@localhost:5432/telkom_tickets?sslmode=disable" -verbose up

