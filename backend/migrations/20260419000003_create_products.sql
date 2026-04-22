-- +goose Up
-- This file should be renamed to 20260419000003_create_products.sql
-- to maintain proper migration sequence (after articles)
CREATE TABLE products (
    product_id VARCHAR(20) PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    product_quantity INT NOT NULL DEFAULT 0,
    product_prices DECIMAL(10, 2) NOT NULL,
    product_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
    image_path VARCHAR(255) DEFAULT 'assets/default_image.jpg'
);

-- Basic indexes
CREATE INDEX idx_products_created_by ON products(created_by);
CREATE INDEX idx_products_type ON products(product_type);

-- +goose Down
DROP TABLE IF EXISTS products;
