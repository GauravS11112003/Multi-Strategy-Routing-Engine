package repository

import (
	"context"
	"fmt"
	"time"

	"shipt-route-optimizer/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ShopperRecord struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Capacity  int       `json:"capacity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func CreateShopper(ctx context.Context, name string, lat, lng float64, capacity int) (*ShopperRecord, error) {
	shopper := &ShopperRecord{}
	err := database.Pool.QueryRow(ctx,
		`INSERT INTO shoppers (name, lat, lng, capacity)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, lat, lng, capacity, status, created_at, updated_at`,
		name, lat, lng, capacity,
	).Scan(&shopper.ID, &shopper.Name, &shopper.Lat, &shopper.Lng, &shopper.Capacity,
		&shopper.Status, &shopper.CreatedAt, &shopper.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert shopper: %w", err)
	}
	return shopper, nil
}

func GetShopperByID(ctx context.Context, id uuid.UUID) (*ShopperRecord, error) {
	shopper := &ShopperRecord{}
	err := database.Pool.QueryRow(ctx,
		`SELECT id, name, lat, lng, capacity, status, created_at, updated_at
		 FROM shoppers WHERE id = $1`, id,
	).Scan(&shopper.ID, &shopper.Name, &shopper.Lat, &shopper.Lng, &shopper.Capacity,
		&shopper.Status, &shopper.CreatedAt, &shopper.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get shopper: %w", err)
	}
	return shopper, nil
}

func ListShoppers(ctx context.Context, status string, limit, offset int) ([]ShopperRecord, int, error) {
	var total int
	query := `SELECT count(*) FROM shoppers`
	args := []interface{}{}

	if status != "" {
		query += ` WHERE status = $1`
		args = append(args, status)
	}

	if err := database.Pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count shoppers: %w", err)
	}

	selectQuery := `SELECT id, name, lat, lng, capacity, status, created_at, updated_at
		FROM shoppers`
	selectArgs := []interface{}{}
	argIdx := 1

	if status != "" {
		selectQuery += fmt.Sprintf(` WHERE status = $%d`, argIdx)
		selectArgs = append(selectArgs, status)
		argIdx++
	}

	selectQuery += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	selectArgs = append(selectArgs, limit, offset)

	rows, err := database.Pool.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list shoppers: %w", err)
	}
	defer rows.Close()

	shoppers := []ShopperRecord{}
	for rows.Next() {
		var s ShopperRecord
		if err := rows.Scan(&s.ID, &s.Name, &s.Lat, &s.Lng, &s.Capacity,
			&s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan shopper: %w", err)
		}
		shoppers = append(shoppers, s)
	}

	return shoppers, total, nil
}

func UpdateShopperStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := database.Pool.Exec(ctx,
		`UPDATE shoppers SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("update shopper status: %w", err)
	}
	return nil
}

func DeleteShopper(ctx context.Context, id uuid.UUID) error {
	_, err := database.Pool.Exec(ctx, `DELETE FROM shoppers WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete shopper: %w", err)
	}
	return nil
}
