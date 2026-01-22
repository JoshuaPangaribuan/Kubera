-- name: CreateCoupon :one
INSERT INTO coupons (name, amount, remaining_amount)
VALUES ($1, $2, $2)
RETURNING *;

-- name: GetCouponByName :one
SELECT * FROM coupons
WHERE name = $1
LIMIT 1;

-- name: GetCouponForUpdate :one
SELECT * FROM coupons
WHERE name = $1
FOR UPDATE;

-- name: UpdateRemainingAmount :one
UPDATE coupons
SET remaining_amount = $1
WHERE name = $2
RETURNING *;

-- name: ListAllCoupons :many
SELECT * FROM coupons
ORDER BY created_at DESC;

-- name: DeleteCoupon :exec
DELETE FROM coupons
WHERE name = $1;
