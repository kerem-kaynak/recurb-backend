# Recurb

Recurb is a simple web application that helps you track your recurring expenses and set reminders to ease the burden of subscription hell. You can create subscriptions on the application, which will automatically generate payment plans. You can then track your monthly expenses, most expensive categories and set reminders on your subscriptions.

## Technology

- Golang / Gin REST API
- Google Cloud SQL PostgreSQL database
- GORM
- Google OAuth + Goth for authentication
- Docker for containerization

## Requirements
- Go v1.2+
- Docker (optional)
- .env file (example in the repo)

## Local development

Install dependencies
```

go mod tidy

```

Build the project
```

go build -o recurb

```

Run the project
```

go run cmd/recurb-api/main.go

```

Run the project in a Docker container
```

docker build -t <image-name> .
docker run -p 8080:8080 --env-file .env <image-name>

```
