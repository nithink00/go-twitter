# 🐦 Twitter Backend API

A complete, production-ready RESTful API backend for a Twitter-like social media platform built with **Go**, **Gin**, and **MySQL**.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## ✨ Features

### Core Functionality

- 🔐 **Complete Authentication System** - Register, Login, JWT tokens, Refresh tokens, Logout
- 📝 **Post Management** - Create, Read, Update, Delete posts with pagination
- 💬 **Comment System** - Comment on posts with full CRUD operations
- ❤️ **Like System** - Like/unlike posts and comments
- 👤 **User Profiles** - View user information and their posts
- 🔒 **Security** - Password hashing, JWT authentication, protected routes
- 📄 **Pagination** - Efficient data fetching for posts and comments
- 🗑️ **Soft Deletes** - Posts and comments are soft deleted for data integrity

### Technical Highlights

- ⚡ **Clean Architecture** - Hexagonal architecture with clear separation of concerns
- 🎯 **RESTful Design** - Standard HTTP methods and status codes
- 🛡️ **Security Best Practices** - bcrypt hashing, JWT tokens, input validation
- 📊 **Scalable Structure** - Easy to extend and maintain
- 🔄 **Dependency Injection** - Loosely coupled components
- ✅ **Production Ready** - Builds successfully with no errors

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [API Endpoints](#-api-endpoints)
- [Project Structure](#-project-structure)
- [Tech Stack](#-tech-stack)
- [Database Schema](#-database-schema)
- [Authentication](#-authentication)
- [Documentation](#-documentation)
- [Contributing](#-contributing)
- [License](#-license)

## 🚀 Quick Start

### Prerequisites

- Go 1.25 or higher
- MySQL 8.0 or Docker
- Git

### Installation

1. **Clone the repository**

```bash
cd "/Users/nithinkatla00/My-projects/Go Sample Projects/Go Twitter"
```

2. **Set up environment variables**

```bash
# Edit .env file with your configuration
PORT=8080
JWT_SECRET=your-secret-key
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=twitter_db
```

3. **Start MySQL database**

```bash
docker-compose up -d
```

4. **Run database migrations**

```bash
cat db/migrations/*.sql | mysql -u root -p twitter_db
```

5. **Install dependencies**

```bash
go mod download
```

6. **Run the application**

```bash
go run cmd/main.go
```

The server will start at `http://localhost:8080` 🎉

### Quick Test

```bash
# Register a user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123",
    "password_confirm": "password123"
  }'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

For detailed setup instructions, see [QUICKSTART.md](./QUICKSTART.md)

## 🔌 API Endpoints

### Authentication

| Method | Endpoint         | Description          | Auth |
| ------ | ---------------- | -------------------- | ---- |
| POST   | `/auth/register` | Register new user    | No   |
| POST   | `/auth/login`    | Login user           | No   |
| POST   | `/auth/refresh`  | Refresh access token | No   |
| POST   | `/auth/logout`   | Logout user          | No   |

### Users

| Method | Endpoint     | Description      | Auth |
| ------ | ------------ | ---------------- | ---- |
| GET    | `/users/:id` | Get user profile | No   |

### Posts

| Method | Endpoint     | Description               | Auth |
| ------ | ------------ | ------------------------- | ---- |
| POST   | `/posts`     | Create new post           | Yes  |
| GET    | `/posts`     | Get all posts (paginated) | No   |
| GET    | `/posts/:id` | Get single post           | No   |
| PUT    | `/posts/:id` | Update post (owner only)  | Yes  |
| DELETE | `/posts/:id` | Delete post (owner only)  | Yes  |

### Comments

| Method | Endpoint                   | Description                 | Auth |
| ------ | -------------------------- | --------------------------- | ---- |
| POST   | `/posts/:post_id/comments` | Create comment              | Yes  |
| GET    | `/posts/:post_id/comments` | Get comments (paginated)    | No   |
| GET    | `/comments/:id`            | Get single comment          | No   |
| PUT    | `/comments/:id`            | Update comment (owner only) | Yes  |
| DELETE | `/comments/:id`            | Delete comment (owner only) | Yes  |

### Likes

| Method | Endpoint                            | Description             | Auth |
| ------ | ----------------------------------- | ----------------------- | ---- |
| POST   | `/posts/:post_id/likes`             | Like a post             | Yes  |
| DELETE | `/posts/:post_id/likes`             | Unlike a post           | Yes  |
| GET    | `/posts/:post_id/likes/count`       | Get post likes count    | Yes  |
| POST   | `/comments/:comment_id/likes`       | Like a comment          | Yes  |
| DELETE | `/comments/:comment_id/likes`       | Unlike a comment        | Yes  |
| GET    | `/comments/:comment_id/likes/count` | Get comment likes count | Yes  |

**Total: 21 API Endpoints**

For detailed API documentation with request/response examples, see [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

## 📁 Project Structure

```
go-twitter/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── dto/                    # Data Transfer Objects
│   ├── handler/                # HTTP handlers
│   │   ├── user/              # User endpoints
│   │   ├── post/              # Post endpoints
│   │   ├── comment/           # Comment endpoints
│   │   └── like/              # Like endpoints
│   ├── middleware/             # JWT auth middleware
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
│   └── migrations/             # Database migrations (6 files)
├── docker-compose.yml          # Docker configuration
├── go.mod                      # Go modules
└── .env                        # Environment variables
```
