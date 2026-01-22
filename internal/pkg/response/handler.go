package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/joshuarp/kubera/internal/pkg/error"
)

// HandleError handles errors in HTTP handlers.
// Returns true if an error was handled (and the caller should return),
// false if there was no error (and the caller should continue).
func HandleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	switch e := err.(type) {
	case *apperrors.AppError:
		c.JSON(e.Code, gin.H{"error": e.Message})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
	return true
}

// HandleValidationError handles JSON binding validation errors.
func HandleValidationError(c *gin.Context, err error, message string) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
	}
}
