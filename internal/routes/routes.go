package routes

import (
	"net/http"

	"github.com/llinder/svtauthpoc/internal/grant"
	"github.com/llinder/svtauthpoc/internal/model"

	"github.com/gin-gonic/gin"
)

func GetHealth() gin.HandlerFunc {
	healthOk := map[string]string{"status": "ok"}
	return func(c *gin.Context) {
		c.JSON(200, healthOk)
	}
}

func PostGrant(client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var grantRequest model.GrantRequest
		if err := c.ShouldBind(&grantRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		token, err := grant.DoGrant(client, &grantRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(200, gin.H{"token": token})
	}
}
