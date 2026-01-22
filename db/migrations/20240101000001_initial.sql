-- +goose Up
-- +goose StatementBegin

-- Create coupons table
CREATE TABLE IF NOT EXISTS coupons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    amount INT NOT NULL CHECK (amount > 0),
    remaining_amount INT NOT NULL CHECK (remaining_amount >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create claims table with unique constraint on (user_id, coupon_name)
CREATE TABLE IF NOT EXISTS claims (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    coupon_name VARCHAR(255) NOT NULL,
    claimed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_coupon UNIQUE (user_id, coupon_name)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_coupons_name ON coupons(name);
CREATE INDEX IF NOT EXISTS idx_claims_user_coupon ON claims(user_id, coupon_name);
CREATE INDEX IF NOT EXISTS idx_claims_coupon_name ON claims(coupon_name);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for coupons table
DROP TRIGGER IF EXISTS update_coupons_updated_at ON coupons;
CREATE TRIGGER update_coupons_updated_at
    BEFORE UPDATE ON coupons
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop trigger
DROP TRIGGER IF EXISTS update_coupons_updated_at ON coupons;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_claims_coupon_name;
DROP INDEX IF EXISTS idx_claims_user_coupon;
DROP INDEX IF EXISTS idx_coupons_name;

-- Drop tables (claims must be dropped first due to foreign key)
DROP TABLE IF EXISTS claims;
DROP TABLE IF EXISTS coupons;

-- +goose StatementEnd
