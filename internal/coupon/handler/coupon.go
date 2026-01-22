package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	claimerrors "github.com/joshuarp/kubera/internal/claim/errors"
	claimservice "github.com/joshuarp/kubera/internal/claim/service"
	"github.com/joshuarp/kubera/internal/coupon/entity"
	"github.com/joshuarp/kubera/internal/coupon/service"
	apperrors "github.com/joshuarp/kubera/internal/pkg/error"
)

// Handler handles coupon-related HTTP requests
type Handler struct {
	service      service.Service
	claimService claimservice.Service
}

// New creates a new coupon handler
func New(service service.Service, claimService claimservice.Service) *Handler {
	return &Handler{
		service:      service,
		claimService: claimService,
	}
}

// CreateCoupon handles POST /api/coupons
func (h *Handler) CreateCoupon(c *gin.Context) {
	var req entity.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": claimerrors.ErrInvalidPayload.Message})
		return
	}

	coupon, err := h.service.CreateCoupon(c.Request.Context(), &req)
	if err != nil {
		switch e := err.(type) {
		case *apperrors.AppError:
			c.JSON(e.Code, gin.H{"error": e.Message})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, coupon)
}

// GetCoupon handles GET /api/coupons/:name
func (h *Handler) GetCoupon(c *gin.Context) {
	name := c.Param("name")

	coupon, err := h.service.GetCoupon(c.Request.Context(), name)
	if err != nil {
		switch e := err.(type) {
		case *apperrors.AppError:
			c.JSON(e.Code, gin.H{"error": e.Message})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Get claimed by list
	var claimedBy []string
	if h.claimService != nil {
		claimedBy, err = h.claimService.GetClaimedByCoupon(c.Request.Context(), name)
		if err != nil {
			// If there's an error getting claimed by, just return empty list
			claimedBy = []string{}
		}
	} else {
		claimedBy = []string{}
	}

	response := &entity.CouponResponse{
		Name:            coupon.Name,
		Amount:          coupon.Amount,
		RemainingAmount: coupon.RemainingAmount,
		ClaimedBy:       claimedBy,
	}

	c.JSON(http.StatusOK, response)
}

// ListCoupons handles GET /api/coupons
func (h *Handler) ListCoupons(c *gin.Context) {
	coupons, err := h.service.ListCoupons(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coupons"})
		return
	}

	c.JSON(http.StatusOK, coupons)
}

// RegisterRoutes registers coupon routes
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	coupons := r.Group("/coupons")
	{
		coupons.POST("", h.CreateCoupon)
		coupons.GET("/:name", h.GetCoupon)
		coupons.GET("", h.ListCoupons)
	}
}
