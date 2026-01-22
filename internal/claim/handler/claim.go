package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshuarp/kubera/internal/claim/entity"
	claimerrors "github.com/joshuarp/kubera/internal/claim/errors"
	"github.com/joshuarp/kubera/internal/claim/service"
	apperrors "github.com/joshuarp/kubera/internal/pkg/error"
)

// Handler handles claim-related HTTP requests
type Handler struct {
	service service.Service
}

// New creates a new claim handler
func New(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ClaimCoupon handles POST /api/coupons/claim
func (h *Handler) ClaimCoupon(c *gin.Context) {
	var req entity.ClaimCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": claimerrors.ErrInvalidPayload.Message})
		return
	}

	result, err := h.service.ClaimCoupon(c.Request.Context(), &req)
	if err != nil {
		switch e := err.(type) {
		case *apperrors.AppError:
			c.JSON(e.Code, gin.H{"error": e.Message})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, result)
}

// RegisterRoutes registers claim routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/coupons/claim", h.ClaimCoupon)
}
