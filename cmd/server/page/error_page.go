package page

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CustomErrorHandler serves any http error.
func CustomErrorHandler(c *gin.Context) {
	c.HTML(http.StatusNotFound, "error", "404 page not found!")
}
