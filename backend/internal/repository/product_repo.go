package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"rwiratama.com/m/internal/models"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

// Create inserts a new product into the database
func (r *ProductRepository) Create(ctx context.Context, productID, productName string, productQuantity int, productPrices, productType string, createdBy uuid.UUID, imagePath string) (*models.Product, error) {
	// Convert string price to decimal
	price, err := decimal.NewFromString(productPrices)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %w", err)
	}

	query := `
		INSERT INTO products (product_id, product_name, product_quantity, product_prices, product_type, created_by, image_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
	`

	var product models.Product
	err = r.pool.QueryRow(ctx, query, productID, productName, productQuantity, price, productType, createdBy, imagePath).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.ProductQuantity,
		&product.ProductPrices,
		&product.ProductType,
		&product.CreatedAt,
		&product.CreatedBy,
		&product.ImagePath,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &product, nil
}

// GetByID retrieves a product by product_id
func (r *ProductRepository) GetByID(ctx context.Context, productID string) (*models.Product, error) {
	query := `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		WHERE product_id = $1
	`

	var product models.Product
	err := r.pool.QueryRow(ctx, query, productID).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.ProductQuantity,
		&product.ProductPrices,
		&product.ProductType,
		&product.CreatedAt,
		&product.CreatedBy,
		&product.ImagePath,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// GetAll retrieves all products from the database
func (r *ProductRepository) GetAll(ctx context.Context) ([]*models.Product, error) {
	query := `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ProductQuantity,
			&product.ProductPrices,
			&product.ProductType,
			&product.CreatedAt,
			&product.CreatedBy,
			&product.ImagePath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, productID, productName string, productQuantity int, productPrices, productType, imagePath string) (*models.Product, error) {
	// Convert string price to decimal
	price, err := decimal.NewFromString(productPrices)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %w", err)
	}

	query := `
		UPDATE products
		SET product_name = $1, product_quantity = $2, product_prices = $3, product_type = $4, image_path = $5
		WHERE product_id = $6
		RETURNING product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
	`

	var product models.Product
	err = r.pool.QueryRow(ctx, query, productName, productQuantity, price, productType, imagePath, productID).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.ProductQuantity,
		&product.ProductPrices,
		&product.ProductType,
		&product.CreatedAt,
		&product.CreatedBy,
		&product.ImagePath,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &product, nil
}

// Delete removes a product from the database
func (r *ProductRepository) Delete(ctx context.Context, productID string) error {
	query := `DELETE FROM products WHERE product_id = $1`
	commandTag, err := r.pool.Exec(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// GetByType retrieves all products of a specific type
func (r *ProductRepository) GetByType(ctx context.Context, productType string) ([]*models.Product, error) {
	query := `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		WHERE product_type = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, productType)
	if err != nil {
		return nil, fmt.Errorf("failed to query products by type: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ProductQuantity,
			&product.ProductPrices,
			&product.ProductType,
			&product.CreatedAt,
			&product.CreatedBy,
			&product.ImagePath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetLowStockProducts retrieves products with quantity below threshold (uses idx_products_low_quantity)
func (r *ProductRepository) GetLowStockProducts(ctx context.Context, threshold int) ([]*models.Product, error) {
	query := `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		WHERE product_quantity < $1
		ORDER BY product_quantity ASC, created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to query low stock products: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ProductQuantity,
			&product.ProductPrices,
			&product.ProductType,
			&product.CreatedAt,
			&product.CreatedBy,
			&product.ImagePath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetByCreatedBy retrieves all products created by a specific user (uses idx_products_created_by_created_at)
func (r *ProductRepository) GetByCreatedBy(ctx context.Context, createdBy uuid.UUID) ([]*models.Product, error) {
	query := `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		WHERE created_by = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to query products by creator: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ProductQuantity,
			&product.ProductPrices,
			&product.ProductType,
			&product.CreatedAt,
			&product.CreatedBy,
			&product.ImagePath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetByTypeAndQuantity retrieves products of specific type with quantity threshold (uses idx_products_type_quantity)
func (r *ProductRepository) GetByTypeAndQuantity(ctx context.Context, productType string, minQuantity int) ([]*models.Product, error) {
	query := `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		WHERE product_type = $1 AND product_quantity >= $2
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, productType, minQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to query products by type and quantity: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ProductQuantity,
			&product.ProductPrices,
			&product.ProductType,
			&product.CreatedAt,
			&product.CreatedBy,
			&product.ImagePath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetByNameLike retrieves products matching a name pattern (uses idx_products_name)
func (r *ProductRepository) GetByNameLike(ctx context.Context, namePattern string) ([]*models.Product, error) {
	query := `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		WHERE product_name ILIKE $1
		ORDER BY created_at DESC
		LIMIT 100
	`

	pattern := "%" + namePattern + "%"
	rows, err := r.pool.Query(ctx, query, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to query products by name: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ProductQuantity,
			&product.ProductPrices,
			&product.ProductType,
			&product.CreatedAt,
			&product.CreatedBy,
			&product.ImagePath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetAllWithPagination retrieves products with pagination (uses idx_products_created_at)
func (r *ProductRepository) GetAllWithPagination(ctx context.Context, limit, offset int) ([]*models.Product, int, error) {
	// Get total count
	var totalCount int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	query := `
		SELECT product_id, product_name, product_quantity, product_prices, product_type, created_at, created_by, image_path
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query products with pagination: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ProductQuantity,
			&product.ProductPrices,
			&product.ProductType,
			&product.CreatedAt,
			&product.CreatedBy,
			&product.ImagePath,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating products: %w", err)
	}

	return products, totalCount, nil
}
