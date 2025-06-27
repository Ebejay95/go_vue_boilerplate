package main

import (
    "log"
    "myapp/internal/database"
    "myapp/internal/handlers"
    "os"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    db, err := database.Connect()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"*"},
        AllowCredentials: true,
    }))

    userHandler := handlers.NewUserHandler(db)
    api := r.Group("/api")
    {
        api.GET("/users", userHandler.GetUsers)
        api.POST("/users", userHandler.CreateUser)
        api.GET("/users/:id", userHandler.GetUser)
        api.PUT("/users/:id", userHandler.UpdateUser)
        api.DELETE("/users/:id", userHandler.DeleteUser)
    }

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    r.Run(":" + port)
}