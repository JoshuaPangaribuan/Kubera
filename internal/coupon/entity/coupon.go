package entity

import "time"

// Coupon represents the coupon domain entity
type Coupon struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Amount         int       `json:"amount"`
	RemainingAmount int      `json:"remaining_amount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CouponResponse represents the API response for a coupon
type CouponResponse struct {
	Name            string   `json:"name"`
	Amount          int      `json:"amount"`
	RemainingAmount int      `json:"remaining_amount"`
	ClaimedBy       []string `json:"claimed_by"`
}

// CreateCouponRequest represents the request to create a new coupon
type CreateCouponRequest struct {
	Name   string `json:"name" binding:"required"`
	Amount int    `json:"amount" binding:"required,gt=0"`
}

// ToEntity converts from a database model to domain entity
func ToEntity(dbCoupon interface{}) *Coupon {
	// This will be implemented with actual SQLC types
	return &Coupon{}
}
