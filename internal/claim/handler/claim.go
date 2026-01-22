package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshuarp/kubera/internal/claim/entity"
	claimerrors "github.com/joshuarp/kubera/internal/claim/errors"
	claimservice "github.com/joshuarp/kubera/internal/claim/service"
	appresponse "github.com/joshuarp/kubera/internal/pkg/response"
)

// Handler handles claim-related HTTP requests
type Handler struct {
	service *claimservice.Service
}

// New creates a new claim handler
func New(service *claimservice.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ClaimCoupon handles POST /api/coupons/claim
func (h *Handler) ClaimCoupon(c *gin.Context) {
	var req entity.ClaimCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appresponse.HandleValidationError(c, err, claimerrors.ErrInvalidPayload.Message)
		return
	}

	result, err := h.service.ClaimCoupon(c.Request.Context(), &req)
	if appresponse.HandleError(c, err) {
		return
	}

	c.JSON(http.StatusCreated, result)
}

// RegisterRoutes registers claim routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/coupons/claim", h.ClaimCoupon)
}
