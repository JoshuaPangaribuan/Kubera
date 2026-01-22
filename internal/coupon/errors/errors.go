package errors

import (
	apperrors "github.com/joshuarp/kubera/internal/pkg/error"
)

var (
	ErrCouponNotFound      = apperrors.NotFound("coupon not found")
	ErrCouponAlreadyExists = apperrors.Conflict("coupon already exists", nil)
)
