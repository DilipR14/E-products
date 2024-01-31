package main

import (
    "github.com/DilipR14/E-Products/eproducts/database"
    "github.com/DilipR14/E-Products/eproducts/routes"
    "github.com/gin-gonic/gin"
    "log"
)

func main() {
    // Connect to MongoDB
    client := database.DataBaseSet()
    defer client.Disconnect()

    // Initialize the router
    router := gin.Default()

    // Set up routes
    routes.UserRoutes(router)

    // Run the server
    err := router.Run(":8080")
    if err != nil {
        log.Fatal(err)
    }
}
