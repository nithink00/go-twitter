# Twitter Backend API Documentation

A complete RESTful API backend for a Twitter-like social media platform built with Go, Gin, and MySQL.

## Table of Contents
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Database Schema](#database-schema)

## Features

### Core Functionality
- ✅ User authentication (register, login, refresh token, logout)
- ✅ JWT-based authorization
- ✅ Post creation, reading, updating, and deletion (CRUD)
- ✅ Comment system on posts
- ✅ Like system for posts and comments
- ✅ Pagination for posts and comments
- ✅ User profile management
- ✅ Soft deletes for posts and comments

### Security
- Password hashing with bcrypt
- JWT token authentication (60-minute expiration)
- Refresh tokens (7-day expiration)
- Protected routes with middleware
- User ownership validation for updates/deletes

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: Gin Web Framework
- **Database**: MySQL 8.0
- **Authentication**: JWT (golang-jwt/jwt)
- **Validation**: go-playground/validator
- **Password Hashing**: bcrypt

## Getting Started

### Prerequisites
- Go 1.25 or higher
- MySQL 8.0 or higher
- Docker (optional, for MySQL container)

### Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables in `.env`:
```env
PORT=8080
JWT_SECRET=your-secret-key-here
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=twitter_db
DATABASE_URL=mysql://root:password@tcp(localhost:3306)/twitter_db
```

4. Start MySQL (using Docker):
```bash
docker-compose up -d
```

5. Run database migrations:
```bash
# Migrations are located in db/migrations/
# Run them in order: 001 -> 002 -> 003 -> 004 -> 005 -> 006
```

6. Start the application:
```bash
go run cmd/main.go
```

The server will start on `http://127.0.0.1:8080`

## API Endpoints

### Authentication Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepassword",
  "password_confirm": "securepassword"
}
```

**Response**:
```json
{
  "id": 1
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "a1b2c3d4e5f6..."
}
```

#### Refresh Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "a1b2c3d4e5f6..."
}
```

**Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "g7h8i9j0k1l2..."
}
```

#### Logout
```http
POST /auth/logout
Content-Type: application/json

{
  "refresh_token": "a1b2c3d4e5f6..."
}
```

### User Endpoints

#### Get User Profile
```http
GET /users/:id
```

**Response**:
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "user@example.com",
  "created_at": "2024-01-15 10:30:00"
}
```

### Post Endpoints

#### Create Post (Protected)
```http
POST /posts
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "My First Post",
  "content": "This is the content of my post..."
}
```

**Response**:
```json
{
  "id": 1
}
```

#### Get All Posts (with pagination)
```http
GET /posts?page=1&page_size=10
```

**Query Parameters**:
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 10, max: 100) - Items per page
- `user_id` (optional) - Filter posts by user ID

**Response**:
```json
{
  "posts": [
    {
      "id": 1,
      "user_id": 1,
      "username": "johndoe",
      "title": "My First Post",
      "content": "This is the content...",
      "likes_count": 5,
      "comments_count": 3,
      "created_at": "2024-01-15 10:30:00",
      "updated_at": "2024-01-15 10:30:00"
    }
  ],
  "total_count": 100,
  "page": 1,
  "page_size": 10,
  "total_pages": 10
}
```

#### Get Single Post
```http
GET /posts/:id
```

**Response**:
```json
{
  "id": 1,
  "user_id": 1,
  "username": "johndoe",
  "title": "My First Post",
  "content": "This is the content...",
  "likes_count": 5,
  "comments_count": 3,
  "created_at": "2024-01-15 10:30:00",
  "updated_at": "2024-01-15 10:30:00"
}
```

#### Update Post (Protected)
```http
PUT /posts/:id
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content..."
}
```

**Response**:
```json
{
  "id": 1
}
```

#### Delete Post (Protected)
```http
DELETE /posts/:id
Authorization: Bearer {token}
```

**Response**:
```json
{
  "message": "post deleted successfully"
}
```

### Comment Endpoints

#### Create Comment (Protected)
```http
POST /posts/:post_id/comments
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "Great post!"
}
```

**Response**:
```json
{
  "id": 1
}
```

#### Get Comments for Post (with pagination)
```http
GET /posts/:post_id/comments?page=1&page_size=10
```

**Response**:
```json
{
  "comments": [
    {
      "id": 1,
      "post_id": 1,
      "user_id": 2,
      "username": "janedoe",
      "content": "Great post!",
      "likes_count": 2,
      "created_at": "2024-01-15 11:00:00",
      "updated_at": "2024-01-15 11:00:00"
    }
  ],
  "total_count": 50,
  "page": 1,
  "page_size": 10,
  "total_pages": 5
}
```

#### Get Single Comment
```http
GET /comments/:id
```

#### Update Comment (Protected)
```http
PUT /comments/:id
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "Updated comment text"
}
```

**Response**:
```json
{
  "id": 1
}
```

#### Delete Comment (Protected)
```http
DELETE /comments/:id
Authorization: Bearer {token}
```

**Response**:
```json
{
  "message": "comment deleted successfully"
}
```

### Like Endpoints

#### Like Post (Protected)
```http
POST /posts/:post_id/likes
Authorization: Bearer {token}
```

**Response**:
```json
{
  "message": "post liked successfully"
}
```

#### Unlike Post (Protected)
```http
DELETE /posts/:post_id/likes
Authorization: Bearer {token}
```

**Response**:
```json
{
  "message": "post unliked successfully"
}
```

#### Get Post Likes Count (Protected)
```http
GET /posts/:post_id/likes/count
Authorization: Bearer {token}
```

**Response**:
```json
{
  "count": 42
}
```

#### Like Comment (Protected)
```http
POST /comments/:comment_id/likes
Authorization: Bearer {token}
```

**Response**:
```json
{
  "message": "comment liked successfully"
}
```

#### Unlike Comment (Protected)
```http
DELETE /comments/:comment_id/likes
Authorization: Bearer {token}
```

**Response**:
```json
{
  "message": "comment unliked successfully"
}
```

#### Get Comment Likes Count (Protected)
```http
GET /comments/:comment_id/likes/count
Authorization: Bearer {token}
```

**Response**:
```json
{
  "count": 15
}
```

## Authentication

### Protected Routes
Routes marked as "Protected" require a valid JWT token in the Authorization header:

```
Authorization: Bearer {your-jwt-token}
```

### Token Flow
1. Register or login to receive a JWT token and refresh token
2. Use the JWT token in the Authorization header for protected routes
3. JWT tokens expire after 60 minutes
4. Use the refresh token endpoint to get a new JWT token
5. Refresh tokens expire after 7 days
6. Logout to invalidate the refresh token

## Database Schema

### Users Table
- `id` - INT, PRIMARY KEY, AUTO_INCREMENT
- `username` - VARCHAR(50), UNIQUE
- `email` - VARCHAR(100), UNIQUE
- `password` - VARCHAR(500), bcrypt hashed
- `created_at` - TIMESTAMP
- `updated_at` - TIMESTAMP

### Posts Table
- `id` - INT, PRIMARY KEY, AUTO_INCREMENT
- `user_id` - INT, FOREIGN KEY -> users(id)
- `title` - VARCHAR(100)
- `content` - LONGTEXT
- `deleted_at` - TIMESTAMP (soft delete)
- `created_at` - TIMESTAMP
- `updated_at` - TIMESTAMP

### Comments Table
- `id` - INT, PRIMARY KEY, AUTO_INCREMENT
- `post_id` - INT, FOREIGN KEY -> posts(id)
- `user_id` - INT, FOREIGN KEY -> users(id)
- `content` - LONGTEXT
- `deleted_at` - TIMESTAMP (soft delete)
- `created_at` - TIMESTAMP
- `updated_at` - TIMESTAMP

### Post Likes Table
- `id` - INT, PRIMARY KEY, AUTO_INCREMENT
- `post_id` - INT, FOREIGN KEY -> posts(id)
- `user_id` - INT, FOREIGN KEY -> users(id)
- `created_at` - TIMESTAMP
- `updated_at` - TIMESTAMP

### Comment Likes Table
- `id` - INT, PRIMARY KEY, AUTO_INCREMENT
- `comment_id` - INT, FOREIGN KEY -> comments(id)
- `user_id` - INT, FOREIGN KEY -> users(id)
- `created_at` - TIMESTAMP
- `updated_at` - TIMESTAMP

### Refresh Tokens Table
- `id` - INT, PRIMARY KEY, AUTO_INCREMENT
- `user_id` - INT, FOREIGN KEY -> users(id)
- `refresh_token` - TEXT
- `expires_at` - TIMESTAMP
- `created_at` - TIMESTAMP
- `updated_at` - TIMESTAMP

## Error Responses

All error responses follow this format:
```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., already liked)
- `500 Internal Server Error` - Server error

## Project Structure

```
go-twitter/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── dto/                    # Data Transfer Objects
│   ├── handler/                # HTTP request handlers
│   │   ├── user/
│   │   ├── post/
│   │   ├── comment/
│   │   └── like/
│   ├── middleware/             # Middleware (auth, etc.)
│   ├── model/                  # Domain models
│   ├── repository/             # Database access layer
│   │   ├── user/
│   │   ├── post/
│   │   ├── comment/
│   │   └── like/
│   └── service/                # Business logic layer
│       ├── user/
│       ├── post/
│       ├── comment/
│       └── like/
├── pkg/
│   ├── internalsql/            # MySQL utilities
│   ├── jwt/                    # JWT token generation
│   └── refreshtoken/           # Refresh token generation
├── db/
│   └── migrations/             # Database migrations
├── docker-compose.yml
├── go.mod
├── go.sum
└── .env
```

## Architecture

The application follows a **layered hexagonal architecture**:

1. **Handler Layer**: Handles HTTP requests/responses, validation
2. **Service Layer**: Contains business logic
3. **Repository Layer**: Handles database operations
4. **Model Layer**: Defines data structures

### Benefits:
- Clean separation of concerns
- Easy to test individual layers
- Flexible and maintainable
- Follows SOLID principles

## Development Notes

### Adding New Features

1. Create model in `internal/model/`
2. Create DTOs in `internal/dto/`
3. Create repository interface and implementation in `internal/repository/`
4. Create service interface and implementation in `internal/service/`
5. Create handlers in `internal/handler/`
6. Wire dependencies in `cmd/main.go`
7. Add routes in handler's `RouteList()` method

### Running Migrations

Migrations are located in `db/migrations/` and should be run in order:
1. `001_create_users_table.sql`
2. `002_create_refresh_tokens_table.sql`
3. `003_create_posts_table.sql`
4. `004_create_comments_table.sql`
5. `005_create_post_likes_table.sql`
6. `006_create_comment_likes_table.sql`

## Future Enhancements

Potential features to add:
- [ ] Follow/unfollow users
- [ ] User feed (posts from followed users)
- [ ] Search functionality
- [ ] Trending posts/hashtags
- [ ] User profile updates
- [ ] Image uploads for posts
- [ ] Rate limiting
- [ ] Comprehensive error handling
- [ ] API documentation with Swagger
- [ ] Unit and integration tests
- [ ] Logging middleware
- [ ] Request validation middleware
- [ ] CORS configuration

## License

MIT License
