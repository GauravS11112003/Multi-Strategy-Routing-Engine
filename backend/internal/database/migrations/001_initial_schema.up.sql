CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE shoppers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,
    capacity INT NOT NULL DEFAULT 10,
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,
    item_count INT NOT NULL,
    delivery_window VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE optimization_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    algorithm VARCHAR(50) NOT NULL,
    total_orders INT NOT NULL,
    total_shoppers INT NOT NULL,
    distance_before DOUBLE PRECISION NOT NULL,
    distance_after DOUBLE PRECISION NOT NULL,
    improvement_pct DOUBLE PRECISION NOT NULL,
    duration_ms INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL REFERENCES optimization_runs(id) ON DELETE CASCADE,
    shopper_id UUID NOT NULL,
    order_id UUID NOT NULL,
    sequence_num INT NOT NULL,
    distance DOUBLE PRECISION NOT NULL
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_shoppers_status ON shoppers(status);
CREATE INDEX idx_assignments_run_id ON assignments(run_id);
CREATE INDEX idx_optimization_runs_created_at ON optimization_runs(created_at DESC);
