package repository

import (
	"context"
	"fmt"
	"time"

	"shipt-route-optimizer/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OptimizationRunRecord struct {
	ID             uuid.UUID `json:"id"`
	Algorithm      string    `json:"algorithm"`
	TotalOrders    int       `json:"totalOrders"`
	TotalShoppers  int       `json:"totalShoppers"`
	DistanceBefore float64   `json:"distanceBefore"`
	DistanceAfter  float64   `json:"distanceAfter"`
	ImprovementPct float64   `json:"improvementPct"`
	DurationMs     int       `json:"durationMs"`
	CreatedAt      time.Time `json:"createdAt"`
}

type AssignmentRecord struct {
	ID          uuid.UUID `json:"id"`
	RunID       uuid.UUID `json:"runId"`
	ShopperID   string    `json:"shopperId"`
	OrderID     string    `json:"orderId"`
	SequenceNum int       `json:"sequenceNum"`
	Distance    float64   `json:"distance"`
}

func SaveOptimizationRun(ctx context.Context, algorithm string, totalOrders, totalShoppers int,
	distanceBefore, distanceAfter, improvementPct float64, durationMs int) (*OptimizationRunRecord, error) {

	run := &OptimizationRunRecord{}
	err := database.Pool.QueryRow(ctx,
		`INSERT INTO optimization_runs (algorithm, total_orders, total_shoppers, distance_before, distance_after, improvement_pct, duration_ms)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, algorithm, total_orders, total_shoppers, distance_before, distance_after, improvement_pct, duration_ms, created_at`,
		algorithm, totalOrders, totalShoppers, distanceBefore, distanceAfter, improvementPct, durationMs,
	).Scan(&run.ID, &run.Algorithm, &run.TotalOrders, &run.TotalShoppers,
		&run.DistanceBefore, &run.DistanceAfter, &run.ImprovementPct, &run.DurationMs, &run.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert optimization run: %w", err)
	}
	return run, nil
}

func SaveAssignment(ctx context.Context, runID uuid.UUID, shopperID, orderID string, sequenceNum int, distance float64) error {
	_, err := database.Pool.Exec(ctx,
		`INSERT INTO assignments (run_id, shopper_id, order_id, sequence_num, distance)
		 VALUES ($1, $2, $3, $4, $5)`,
		runID, shopperID, orderID, sequenceNum, distance,
	)
	if err != nil {
		return fmt.Errorf("insert assignment: %w", err)
	}
	return nil
}

func GetOptimizationRun(ctx context.Context, id uuid.UUID) (*OptimizationRunRecord, error) {
	run := &OptimizationRunRecord{}
	err := database.Pool.QueryRow(ctx,
		`SELECT id, algorithm, total_orders, total_shoppers, distance_before, distance_after, improvement_pct, duration_ms, created_at
		 FROM optimization_runs WHERE id = $1`, id,
	).Scan(&run.ID, &run.Algorithm, &run.TotalOrders, &run.TotalShoppers,
		&run.DistanceBefore, &run.DistanceAfter, &run.ImprovementPct, &run.DurationMs, &run.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get optimization run: %w", err)
	}
	return run, nil
}

func ListOptimizationRuns(ctx context.Context, limit, offset int) ([]OptimizationRunRecord, int, error) {
	var total int
	if err := database.Pool.QueryRow(ctx, `SELECT count(*) FROM optimization_runs`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count optimization runs: %w", err)
	}

	rows, err := database.Pool.Query(ctx,
		`SELECT id, algorithm, total_orders, total_shoppers, distance_before, distance_after, improvement_pct, duration_ms, created_at
		 FROM optimization_runs
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list optimization runs: %w", err)
	}
	defer rows.Close()

	runs := []OptimizationRunRecord{}
	for rows.Next() {
		var r OptimizationRunRecord
		if err := rows.Scan(&r.ID, &r.Algorithm, &r.TotalOrders, &r.TotalShoppers,
			&r.DistanceBefore, &r.DistanceAfter, &r.ImprovementPct, &r.DurationMs, &r.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan optimization run: %w", err)
		}
		runs = append(runs, r)
	}

	return runs, total, nil
}

func GetRunAssignments(ctx context.Context, runID uuid.UUID) ([]AssignmentRecord, error) {
	rows, err := database.Pool.Query(ctx,
		`SELECT id, run_id, shopper_id, order_id, sequence_num, distance
		 FROM assignments
		 WHERE run_id = $1
		 ORDER BY shopper_id, sequence_num`, runID,
	)
	if err != nil {
		return nil, fmt.Errorf("get assignments: %w", err)
	}
	defer rows.Close()

	assignments := []AssignmentRecord{}
	for rows.Next() {
		var a AssignmentRecord
		if err := rows.Scan(&a.ID, &a.RunID, &a.ShopperID, &a.OrderID,
			&a.SequenceNum, &a.Distance); err != nil {
			return nil, fmt.Errorf("scan assignment: %w", err)
		}
		assignments = append(assignments, a)
	}

	return assignments, nil
}
