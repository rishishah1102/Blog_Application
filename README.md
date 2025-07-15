# Blog Application API

A RESTful API for managing blog posts built with Go, Gin, and MongoDB.

## Features

- Create, read, update, and delete blog posts
- Structured logging with Zap
- MongoDB persistence
- Unit tests with high coverage
- Swagger documentation

## API Endpoints

| Method | Endpoint           | Description                |
|--------|--------------------|----------------------------|
| POST   | /api/blog-post     | Create a new blog post     |
| GET    | /api/blog-post     | Get all blog posts         |
| GET    | /api/blog-post/:id | Get a single blog post     |
| PATCH  | /api/blog-post/:id | Update a blog post         |
| DELETE | /api/blog-post/:id | Delete a blog post         |

## Development

### Prerequisites

- Go (Gin)
- MongoDB
- Swag (for documentation)

### Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Set up environment variables in `.env` file
4. Generate Swagger docs: `swag init`
5. Run the application: `go run main.go`

### Testing

Run tests with coverage:
```bash
go test -coverprofile=coverage ./...
go tool cover -html=coverage