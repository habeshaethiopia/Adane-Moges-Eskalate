# Eskalate Movie API

A RESTful API for managing a personal movie collection. This project allows users to create accounts, authenticate, and manage their movie collections with features like creating, updating, deleting, and searching movies.

## Technology Choices

- **Go (Golang)**: Chosen for its excellent performance, strong standard library, and great support for building concurrent web services.
- **Gin Framework**: A high-performance HTTP web framework that provides a great balance between performance and developer productivity.
- **GORM**: The most widely-used ORM for Go, providing excellent database abstraction and features like auto-migrations.
- **PostgreSQL**: A robust, open-source relational database with excellent support for UUID, JSON, and full-text search capabilities.
- **JWT**: For secure, stateless authentication.
- **Swagger/OpenAPI**: For comprehensive API documentation.
- **Cloudinary**: For efficient cloud-based image storage and management.

## Local Setup

### Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Git

### Installation Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/habeshaethiopia/eskalate-movie-api.git
   cd eskalate-movie-api
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Set up environment variables by creating a `.env` file in the project root:

   ```env
   # Server Configuration
   PORT=:8080

   # Database Configuration
   DATABASE_URL=postgresql://username:password@localhost:5432/movie_db?sslmode=disable

   # JWT Configuration
   JWT_SECRET=your_jwt_secret_key
   JWT_EXPIRATION_HOURS=24

   # Cloudinary Configuration
   CLOUDINARY_CLOUD_NAME=your_cloud_name
   CLOUDINARY_API_KEY=your_api_key
   CLOUDINARY_API_SECRET=your_api_secret
   ```

4. Create the database:

   ```bash
   createdb movie_db
   ```

5. Run the application:
   ```bash
   go run main.go
   ```

The server will start at `http://localhost:8080`

## API Documentation

Once the server is running, you can access:

- API Landing Page: `http://localhost:8080`
- Swagger Documentation: `http://localhost:8080/swagger/index.html`

## Environment Variables

| Variable              | Description                   | Required | Default |
| --------------------- | ----------------------------- | -------- | ------- |
| PORT                  | Server port                   | No       | :8080   |
| DATABASE_URL          | PostgreSQL connection string  | Yes      | -       |
| JWT_SECRET            | Secret key for JWT signing    | Yes      | -       |
| JWT_EXPIRATION_HOURS  | JWT token expiration in hours | No       | 24      |
| CLOUDINARY_CLOUD_NAME | Cloudinary cloud name         | Yes      | -       |
| CLOUDINARY_API_KEY    | Cloudinary API key            | Yes      | -       |
| CLOUDINARY_API_SECRET | Cloudinary API secret         | Yes      | -       |

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   ├── models/             # Database models
│   ├── repository/         # Database operations
│   ├── routes/             # Route definitions
│   └── utils/              # Utility functions
├── static/                 # Static files for landing page
├── docs/                   # Swagger documentation
├── go.mod                  # Go module file
├── go.sum                  # Go module checksum
└── README.md              # Project documentation
```

## Features

- User authentication (signup/login) with JWT
- CRUD operations for movies
- Movie search functionality
- Image upload for movie posters
- Pagination for movie listings
- Swagger documentation
- Beautiful landing page with API documentation

## API Endpoints

### Authentication

- `POST /api/auth/signup` - Register a new user
- `POST /api/auth/login` - Login user

### Movies

- `GET /api/movies` - Get paginated list of movies
- `POST /api/movies` - Create a new movie (auth required)
- `GET /api/movies/search` - Search movies by title
- `GET /api/movies/{id}` - Get movie details
- `PUT /api/movies/{id}` - Update movie (auth required)
- `DELETE /api/movies/{id}` - Delete movie (auth required)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
