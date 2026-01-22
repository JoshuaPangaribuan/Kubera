package entity

import "time"

// Claim represents the claim domain entity
type Claim struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	CouponName string    `json:"coupon_name"`
	ClaimedAt  time.Time `json:"claimed_at"`
}

// ClaimCouponRequest represents the request to claim a coupon
type ClaimCouponRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	CouponName string `json:"coupon_name" binding:"required"`
}
