package errors

import (
	apperrors "github.com/joshuarp/kubera/internal/pkg/error"
)

var (
	ErrNoStockAvailable = apperrors.BadRequest("no stock available", nil)
	ErrAlreadyClaimed   = apperrors.Conflict("coupon already claimed by this user", nil)
	ErrInvalidPayload   = apperrors.BadRequest("invalid request payload", nil)
)
