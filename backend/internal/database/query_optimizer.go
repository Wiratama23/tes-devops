package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// QueryOptimizer provides index-aware query optimization utilities
type QueryOptimizer struct {
	pool *pgxpool.Pool
}

// NewQueryOptimizer creates a new query optimizer
func NewQueryOptimizer(pool *pgxpool.Pool) *QueryOptimizer {
	return &QueryOptimizer{pool: pool}
}

// ExplainQuery returns the query execution plan
func (qo *QueryOptimizer) ExplainQuery(ctx context.Context, query string, args ...interface{}) ([]string, error) {
	explainQuery := "EXPLAIN " + query
	rows, err := qo.pool.Query(ctx, explainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to explain query: %w", err)
	}
	defer rows.Close()

	var plans []string
	for rows.Next() {
		var plan string
		if err := rows.Scan(&plan); err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}
		plans = append(plans, plan)
	}

	return plans, rows.Err()
}

// ExplainAnalyzeQuery returns the query execution plan with actual execution statistics
func (qo *QueryOptimizer) ExplainAnalyzeQuery(ctx context.Context, query string, args ...interface{}) ([]string, error) {
	explainQuery := "EXPLAIN ANALYZE " + query
	rows, err := qo.pool.Query(ctx, explainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to explain analyze query: %w", err)
	}
	defer rows.Close()

	var plans []string
	for rows.Next() {
		var plan string
		if err := rows.Scan(&plan); err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}
		plans = append(plans, plan)
	}

	return plans, rows.Err()
}

// GetTableStats returns table statistics
func (qo *QueryOptimizer) GetTableStats(ctx context.Context, tableName string) (map[string]interface{}, error) {
	query := `
		SELECT
			schemaname,
			tablename,
			seq_scan,
			seq_tup_read,
			idx_scan,
			idx_tup_fetch,
			n_tup_ins,
			n_tup_upd,
			n_tup_del,
			n_live_tup,
			n_dead_tup
		FROM pg_stat_user_tables
		WHERE tablename = $1
	`

	var schemaName, tblName string
	var seqScan, seqTupRead, idxScan, idxTupFetch int64
	var nTupIns, nTupUpd, nTupDel, nLiveTup, nDeadTup int64

	err := qo.pool.QueryRow(ctx, query, tableName).Scan(
		&schemaName, &tblName, &seqScan, &seqTupRead, &idxScan, &idxTupFetch,
		&nTupIns, &nTupUpd, &nTupDel, &nLiveTup, &nDeadTup,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get table stats: %w", err)
	}

	return map[string]interface{}{
		"schema_name":   schemaName,
		"table_name":    tblName,
		"seq_scan":      seqScan,
		"seq_tup_read":  seqTupRead,
		"idx_scan":      idxScan,
		"idx_tup_fetch": idxTupFetch,
		"n_tup_ins":     nTupIns,
		"n_tup_upd":     nTupUpd,
		"n_tup_del":     nTupDel,
		"n_live_tup":    nLiveTup,
		"n_dead_tup":    nDeadTup,
	}, nil
}

// GetMissingIndexes returns suggested indexes for inefficient queries
func (qo *QueryOptimizer) GetMissingIndexes(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT
			schemaname,
			tablename,
			attname,
			n_distinct,
			correlation
		FROM pg_stats
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
		ORDER BY abs(correlation) DESC, schemaname, tablename
	`

	rows, err := qo.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query missing indexes: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var schemaName, tableName, attName string
		var nDistinct, correlation float64

		if err := rows.Scan(&schemaName, &tableName, &attName, &nDistinct, &correlation); err != nil {
			return nil, fmt.Errorf("failed to scan missing index: %w", err)
		}

		results = append(results, map[string]interface{}{
			"schema_name": schemaName,
			"table_name":  tableName,
			"column_name": attName,
			"n_distinct":  nDistinct,
			"correlation": correlation,
		})
	}

	return results, rows.Err()
}

// GetSlowQueries returns information about slow queries (requires pg_stat_statements extension)
func (qo *QueryOptimizer) GetSlowQueries(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT
			query,
			calls,
			total_time,
			mean_time,
			max_time,
			rows
		FROM pg_stat_statements
		ORDER BY mean_time DESC
		LIMIT $1
	`

	rows, err := qo.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query slow queries: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var queryStr string
		var calls int64
		var totalTime, meanTime, maxTime float64
		var rowsAffected int64

		if err := rows.Scan(&queryStr, &calls, &totalTime, &meanTime, &maxTime, &rowsAffected); err != nil {
			return nil, fmt.Errorf("failed to scan slow query: %w", err)
		}

		results = append(results, map[string]interface{}{
			"query":      queryStr,
			"calls":      calls,
			"total_time": totalTime,
			"mean_time":  meanTime,
			"max_time":   maxTime,
			"rows":       rowsAffected,
		})
	}

	return results, rows.Err()
}

// CacheInfo returns information about table caches and hits
func (qo *QueryOptimizer) CacheInfo(ctx context.Context, tableName string) (map[string]interface{}, error) {
	query := `
		SELECT
			heap_blks_read,
			heap_blks_hit,
			idx_blks_read,
			idx_blks_hit
		FROM pg_statio_user_tables
		WHERE relname = $1
	`

	var heapBlksRead, heapBlksHit, idxBlksRead, idxBlksHit int64

	err := qo.pool.QueryRow(ctx, query, tableName).Scan(
		&heapBlksRead, &heapBlksHit, &idxBlksRead, &idxBlksHit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache info: %w", err)
	}

	cacheHitRatio := float64(0)
	if (heapBlksRead + heapBlksHit) > 0 {
		cacheHitRatio = float64(heapBlksHit) / float64(heapBlksRead+heapBlksHit) * 100
	}

	return map[string]interface{}{
		"heap_blks_read":  heapBlksRead,
		"heap_blks_hit":   heapBlksHit,
		"idx_blks_read":   idxBlksRead,
		"idx_blks_hit":    idxBlksHit,
		"cache_hit_ratio": cacheHitRatio,
	}, nil
}
