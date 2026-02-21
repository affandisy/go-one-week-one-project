migrate create -ext sql -dir db/migrations -seq create_tickets_table

migrate -path db/migrations -database "postgresql://telkom_user:telkom_password@localhost:5433/telkom_tickets?sslmode=disable" -verbose up

migrate -path migrations -database "postgres://telkom_user:telkom_password@localhost:5433/telkom_tickets?sslmode=disable" -verbose up

go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest

docker logs -f week-01-sistem-layanan-telkom-app-1

docker restart week-01-sistem-layanan-telkom-app-1

{
  "customer_id": "CUST-001",
  "issue": "Internet mati total sejak pagi",
  "priority": "high"
}