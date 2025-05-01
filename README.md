# Healthcare Appointment System API

A RESTful API for managing healthcare appointments built with Go, using Gin framework, GORM, and SQLite database.

## Features

- Patient management (CRUD operations)
- Appointment scheduling and management
- SQLite database for data persistence
- Docker support for easy deployment
- Automated tests
- Sample data seeding

## Prerequisites

- Docker and Docker Compose
- curl (for testing)

## Quick Start with Docker

1. Clone the repository:
```bash
git clone https://github.com/Olt-Kondirolli91/healthcare-api.git
cd healthcare-api
```

2. Start the application:
```bash
docker compose up --build -d
```

The API will be available at `http://localhost:8080`

## API Testing Guide

### Base URL
All endpoints are accessible at `http://localhost:8080`

### Patient Endpoints

#### List all patients
```bash
curl -X GET http://localhost:8080/patients
```

#### Create a new patient
```bash
curl -X POST http://localhost:8080/patients \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "555-123-4567"
  }'
```

#### Get patient by ID
```bash
curl -X GET http://localhost:8080/patients/1
```

#### Update patient
```bash
curl -X PUT http://localhost:8080/patients/1 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Updated",
    "email": "john.updated@example.com",
    "phone": "555-999-8888"
  }'
```

#### Delete patient
```bash
curl -X DELETE http://localhost:8080/patients/1
```

### Appointment Endpoints

#### List all appointments
```bash
curl -X GET http://localhost:8080/appointments
```

#### Create new appointment
```bash
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 1,
    "date_time": "2025-05-15T14:30:00Z",
    "notes": "Regular checkup"
  }'
```

#### Get appointment by ID
```bash
curl -X GET http://localhost:8080/appointments/1
```

#### Update appointment
```bash
curl -X PUT http://localhost:8080/appointments/1 \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 1,
    "date_time": "2025-05-16T15:00:00Z",
    "notes": "Rescheduled checkup"
  }'
```

#### Delete appointment
```bash
curl -X DELETE http://localhost:8080/appointments/1
```

## Development Setup

### Project Structure
```
healthcare-api/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── database/           # Database operations
│   ├── handlers/           # HTTP handlers
│   └── models/             # Data models
├── tests/                  # Unit tests
├── Dockerfile             
├── docker-compose.yml      
└── README.md              
```

### Building from Source

1. Install Go 1.24 or later
2. Clone the repository
3. Install dependencies:
```bash
go mod download
```
4. Run the application:
```bash
go run cmd/server/main.go
```

### Running Tests
```bash
go test ./tests/...
```

## Docker Configuration

### Dockerfile
The provided Dockerfile creates a minimal production image:
- Uses multi-stage build
- Builds with CGO enabled for SQLite support
- Creates final image from Alpine Linux
- Exposes port 8080

### Docker Compose
The docker-compose.yml file provides:
- Automatic build from Dockerfile
- Port mapping (8080:8080)
- Volume mounting for database persistence
- Environment variable configuration

### Volume Management
Data is persisted in a Docker volume named `healthcare_data`. To manage:

List volumes:
```bash
docker volume ls
```

Clean up volume:
```bash
docker compose down -v
```

## Sample Data

The application automatically seeds sample data on first run, including:
- 3 sample patients
- 4 sample appointments

You can start testing the API immediately using the provided curl commands.

## Troubleshooting

### Common Issues

1. Port already in use:
```bash
docker compose down
docker compose up --build -d
```

2. Database reset:
```bash
docker compose down -v
docker compose up --build -d
```

3. Container logs:
```bash
docker compose logs
```

## API Response Formats

### Success Response
```json
{
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response
```json
{
  "error": "Error message description"
}
```