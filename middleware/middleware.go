package middleware

import (
	"net/http"

	"github.com/DilipR14/E-products/tokens" // Update the import path
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		ClientToken := c.Request.Header.Get("token") // Corrected typo in Request

		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		claims, err := tokens.ValidateToken(ClientToken) // Update the tokens package name
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Corrected error handling
			c.Abort()
			return
		}

		c.Set("email", claims.Email) // Corrected Set method name
		c.Set("uid", claims.UID)    // Corrected Set method name
		c.Next()
	}
}
