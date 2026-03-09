package repository

import (
	"context"
	"fmt"
	"time"

	"shipt-route-optimizer/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrderRecord struct {
	ID             uuid.UUID `json:"id"`
	Lat            float64   `json:"lat"`
	Lng            float64   `json:"lng"`
	ItemCount      int       `json:"itemCount"`
	DeliveryWindow string    `json:"deliveryWindow"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func CreateOrder(ctx context.Context, lat, lng float64, itemCount int, deliveryWindow string) (*OrderRecord, error) {
	order := &OrderRecord{}
	err := database.Pool.QueryRow(ctx,
		`INSERT INTO orders (lat, lng, item_count, delivery_window)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, lat, lng, item_count, delivery_window, status, created_at, updated_at`,
		lat, lng, itemCount, deliveryWindow,
	).Scan(&order.ID, &order.Lat, &order.Lng, &order.ItemCount, &order.DeliveryWindow,
		&order.Status, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert order: %w", err)
	}
	return order, nil
}

func GetOrderByID(ctx context.Context, id uuid.UUID) (*OrderRecord, error) {
	order := &OrderRecord{}
	err := database.Pool.QueryRow(ctx,
		`SELECT id, lat, lng, item_count, delivery_window, status, created_at, updated_at
		 FROM orders WHERE id = $1`, id,
	).Scan(&order.ID, &order.Lat, &order.Lng, &order.ItemCount, &order.DeliveryWindow,
		&order.Status, &order.CreatedAt, &order.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	return order, nil
}

func ListOrders(ctx context.Context, status string, limit, offset int) ([]OrderRecord, int, error) {
	var total int
	query := `SELECT count(*) FROM orders`
	args := []interface{}{}

	if status != "" {
		query += ` WHERE status = $1`
		args = append(args, status)
	}

	if err := database.Pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count orders: %w", err)
	}

	selectQuery := `SELECT id, lat, lng, item_count, delivery_window, status, created_at, updated_at
		FROM orders`
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
		return nil, 0, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	orders := []OrderRecord{}
	for rows.Next() {
		var o OrderRecord
		if err := rows.Scan(&o.ID, &o.Lat, &o.Lng, &o.ItemCount, &o.DeliveryWindow,
			&o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

func UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := database.Pool.Exec(ctx,
		`UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	return nil
}

func DeleteOrder(ctx context.Context, id uuid.UUID) error {
	_, err := database.Pool.Exec(ctx, `DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete order: %w", err)
	}
	return nil
}
