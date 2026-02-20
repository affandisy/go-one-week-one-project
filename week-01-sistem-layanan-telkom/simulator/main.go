package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:3000/tickets"

	for i := 1; i <= 5; i++ {
		payload := map[string]string{
			"customer_id": fmt.Sprintf("CUST-100%d", i),
			"issue":       fmt.Sprintf("Simulasi error jaringan ke-%d", i),
			"priority":    "high",
		}

		jsonData, _ := json.Marshal(payload)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Gagal mengirim tiket: ", err)
			continue
		}

		log.Printf("Tiket %d terkirim. Status HTTP: %s \n", i, resp.Status)
		time.Sleep(1 * time.Second)
	}

	log.Println("Simulasi selesai!")
}
