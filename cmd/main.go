package main

import (
	"fmt"
	"go-twitter/internal/config"
	commentHandler "go-twitter/internal/handler/comment"
	likeHandler "go-twitter/internal/handler/like"
	postHandler "go-twitter/internal/handler/post"
	userHandler "go-twitter/internal/handler/user"
	"go-twitter/internal/middleware"
	commentRepo "go-twitter/internal/repository/comment"
	likeRepo "go-twitter/internal/repository/like"
	postRepo "go-twitter/internal/repository/post"
	userRepo "go-twitter/internal/repository/user"
	commentService "go-twitter/internal/service/comment"
	likeService "go-twitter/internal/service/like"
	postService "go-twitter/internal/service/post"
	"go-twitter/internal/service/user"
	"go-twitter/pkg/internalsql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Initialize Gin router
	r := gin.Default()
	validate := validator.New()
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := internalsql.ConnectMySQL(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.SecreetJwt)

	// Initialize repositories
	userRepository := userRepo.NewRepository(db)
	postRepository := postRepo.NewRepository(db)
	commentRepository := commentRepo.NewRepository(db)
	likeRepository := likeRepo.NewRepository(db)

	// Initialize services
	userService := user.NewService(cfg, userRepository)
	postSvc := postService.NewService(cfg, postRepository, db)
	commentSvc := commentService.NewService(commentRepository, userRepository)
	likeSvc := likeService.NewService(likeRepository)

	// Initialize handlers
	userHandlerInstance := userHandler.NewHandler(r, validate, userService)
	postHandlerInstance := postHandler.NewHandler(r, validate, postSvc, authMiddleware)
	commentHandlerInstance := commentHandler.NewHandler(r, validate, commentSvc, authMiddleware)
	likeHandlerInstance := likeHandler.NewHandler(r, likeSvc, authMiddleware)

	// Register routes
	userHandlerInstance.RouteList()
	postHandlerInstance.RouteList()
	commentHandlerInstance.RouteList()
	likeHandlerInstance.RouteList()

	server := fmt.Sprintf("127.0.0.1:%s", cfg.Port)
	fmt.Printf("Server starting on %s\n", server)
	r.Run(server)
}		