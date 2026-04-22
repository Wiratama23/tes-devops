package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// IndexManager handles database index operations
type IndexManager struct {
	pool *pgxpool.Pool
}

// NewIndexManager creates a new index manager
func NewIndexManager(pool *pgxpool.Pool) *IndexManager {
	return &IndexManager{pool: pool}
}

// CreateProductIndexes ensures all product table indexes exist
func (im *IndexManager) CreateProductIndexes(ctx context.Context) error {
	indexes := []struct {
		name  string
		query string
	}{
		{
			name:  "idx_products_type",
			query: "CREATE INDEX IF NOT EXISTS idx_products_type ON products(product_type);",
		},
		{
			name:  "idx_products_created_by",
			query: "CREATE INDEX IF NOT EXISTS idx_products_created_by ON products(created_by);",
		},
		{
			name:  "idx_products_created_at",
			query: "CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at DESC);",
		},
		{
			name:  "idx_products_type_created_at",
			query: "CREATE INDEX IF NOT EXISTS idx_products_type_created_at ON products(product_type, created_at DESC);",
		},
		{
			name:  "idx_products_created_by_created_at",
			query: "CREATE INDEX IF NOT EXISTS idx_products_created_by_created_at ON products(created_by, created_at DESC);",
		},
		{
			name:  "idx_products_name",
			query: "CREATE INDEX IF NOT EXISTS idx_products_name ON products(product_name);",
		},
		{
			name:  "idx_products_low_quantity",
			query: "CREATE INDEX IF NOT EXISTS idx_products_low_quantity ON products(product_quantity) WHERE product_quantity < 10;",
		},
		{
			name:  "idx_products_type_quantity",
			query: "CREATE INDEX IF NOT EXISTS idx_products_type_quantity ON products(product_type, product_quantity);",
		},
	}

	for _, idx := range indexes {
		if err := im.createIndex(ctx, idx.name, idx.query); err != nil {
			return fmt.Errorf("failed to create index %s: %w", idx.name, err)
		}
	}

	return nil
}

// createIndex creates a single index and logs the result
func (im *IndexManager) createIndex(ctx context.Context, indexName, query string) error {
	_, err := im.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute index query: %w", err)
	}
	return nil
}

// DropProductIndexes drops all product table indexes
func (im *IndexManager) DropProductIndexes(ctx context.Context) error {
	indexNames := []string{
		"idx_products_type",
		"idx_products_created_by",
		"idx_products_created_at",
		"idx_products_type_created_at",
		"idx_products_created_by_created_at",
		"idx_products_name",
		"idx_products_low_quantity",
		"idx_products_type_quantity",
	}

	for _, indexName := range indexNames {
		query := fmt.Sprintf("DROP INDEX IF EXISTS %s;", indexName)
		if _, err := im.pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to drop index %s: %w", indexName, err)
		}
	}

	return nil
}

// GetIndexInfo retrieves information about indexes on a table
func (im *IndexManager) GetIndexInfo(ctx context.Context, tableName string) ([]map[string]interface{}, error) {
	query := `
		SELECT
			indexname,
			indexdef,
			tablename
		FROM pg_indexes
		WHERE tablename = $1
		ORDER BY indexname;
	`

	rows, err := im.pool.Query(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query indexes: %w", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var indexName, indexDef, tblName string
		if err := rows.Scan(&indexName, &indexDef, &tblName); err != nil {
			return nil, fmt.Errorf("failed to scan index info: %w", err)
		}

		result = append(result, map[string]interface{}{
			"index_name":       indexName,
			"index_definition": indexDef,
			"table_name":       tblName,
		})
	}

	return result, rows.Err()
}

// AnalyzeIndexUsage returns index usage statistics
func (im *IndexManager) AnalyzeIndexUsage(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT
			schemaname,
			tablename,
			indexname,
			idx_scan as scan_count,
			idx_tup_read as tuples_read,
			idx_tup_fetch as tuples_fetched
		FROM pg_stat_user_indexes
		ORDER BY idx_scan DESC;
	`

	rows, err := im.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query index usage: %w", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var schemaName, tableName, indexName string
		var scanCount, tuplesRead, tuplesFetched int64

		if err := rows.Scan(&schemaName, &tableName, &indexName, &scanCount, &tuplesRead, &tuplesFetched); err != nil {
			return nil, fmt.Errorf("failed to scan index usage: %w", err)
		}

		result = append(result, map[string]interface{}{
			"schema_name":    schemaName,
			"table_name":     tableName,
			"index_name":     indexName,
			"scan_count":     scanCount,
			"tuples_read":    tuplesRead,
			"tuples_fetched": tuplesFetched,
		})
	}

	return result, rows.Err()
}

// ReindexTable reindexes a specific table
func (im *IndexManager) ReindexTable(ctx context.Context, tableName string) error {
	query := fmt.Sprintf("REINDEX TABLE %s;", tableName)
	_, err := im.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to reindex table %s: %w", tableName, err)
	}
	return nil
}
