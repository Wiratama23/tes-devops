-- +goose Up
-- Create comprehensive indexes for the products table

-- Index on product_type for filtering products by category
CREATE INDEX IF NOT EXISTS idx_products_type ON products(product_type);

-- Index on created_by for querying user's products
CREATE INDEX IF NOT EXISTS idx_products_created_by ON products(created_by);

-- Index on created_at for time-based queries and sorting
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at DESC);

-- Composite index on (product_type, created_at) for efficient filtering and sorting by type
CREATE INDEX IF NOT EXISTS idx_products_type_created_at ON products(product_type, created_at DESC);

-- Composite index on (created_by, created_at) for user-specific queries with ordering
CREATE INDEX IF NOT EXISTS idx_products_created_by_created_at ON products(created_by, created_at DESC);

-- Index on product_name for search/filter operations
CREATE INDEX IF NOT EXISTS idx_products_name ON products(product_name);

-- Partial index on product_quantity to quickly find low-stock items
CREATE INDEX IF NOT EXISTS idx_products_low_quantity ON products(product_quantity) WHERE product_quantity < 10;

-- Composite index on (product_type, product_quantity) for filtering by type and stock
CREATE INDEX IF NOT EXISTS idx_products_type_quantity ON products(product_type, product_quantity);

-- +goose Down
-- Drop all indexes created in this migration
DROP INDEX IF EXISTS idx_products_type;
DROP INDEX IF EXISTS idx_products_created_by;
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_type_created_at;
DROP INDEX IF EXISTS idx_products_created_by_created_at;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_low_quantity;
DROP INDEX IF EXISTS idx_products_type_quantity;
