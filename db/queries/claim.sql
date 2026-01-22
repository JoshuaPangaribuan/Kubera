-- name: CreateClaim :one
INSERT INTO claims (user_id, coupon_name)
VALUES ($1, $2)
RETURNING *;

-- name: GetClaimByUserAndCoupon :one
SELECT * FROM claims
WHERE user_id = $1 AND coupon_name = $2
LIMIT 1;

-- name: ListClaimsByCoupon :many
SELECT * FROM claims
WHERE coupon_name = $1
ORDER BY claimed_at ASC;

-- name: CountClaimsByCoupon :one
SELECT COUNT(*) FROM claims
WHERE coupon_name = $1;

-- name: ListClaimedByCoupon :many
SELECT user_id FROM claims
WHERE coupon_name = $1
ORDER BY claimed_at ASC;
