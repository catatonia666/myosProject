package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const isAuthenticatedContextKey string = "isAuthenticated"

// authenticateMiddleware checks if the session presists for the user and sets key for authentication.
func (s *server) authenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Get the session ID from cookie.
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.Set(isAuthenticatedContextKey, false)
			c.Next()
			return
		}

		//Get user ID related to the session from Redis.
		var sessionData map[string]any
		jsonData, _ := s.redisClient.Get(c, sessionID).Result()
		deserialize(jsonData, &sessionData)
		idFloat, _ := sessionData["userID"].(float64)
		id := int(idFloat)

		if id == 0 {
			c.Set(isAuthenticatedContextKey, false)
			c.Next()
			return
		}

		//Check if the user with gathered ID persists in database.
		_, err = s.store.User().FindByID(id)
		if err != nil {
			c.Set(isAuthenticatedContextKey, false)
			return
		}
		c.Set(isAuthenticatedContextKey, true)
	}
}

func (s *server) methodOverrideMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			if override := c.PostForm("_method"); override == http.MethodPut {
				c.Request.Method = http.MethodPut
			}
		}
		c.Next()
	}
}
