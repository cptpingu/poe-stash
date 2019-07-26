package page

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MainPageHandler serves the main page of this website.
func MainPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "main", "my_url")
}
