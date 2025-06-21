CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id    uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    login text NOT NULL UNIQUE,
    hash  text NOT NULL
);

CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders (
    id          text PRIMARY KEY,
    user_id     uuid REFERENCES users NOT NULL,
    status      order_status NOT NULL DEFAULT 'NEW',
    uploaded_at timestamp NOT NULL DEFAULT current_timestamp
);

CREATE INDEX idx_orders_user_id ON orders (user_id);

CREATE TABLE accrual_flow (
    id       uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id text REFERENCES orders NOT NULL,
    amount   numeric(15, 2) NOT NULL DEFAULT 0,
    processed_at timestamp NOT NULL DEFAULT current_timestamp
);

CREATE INDEX idx_accrual_flow_order_id ON accrual_flow (order_id);

CREATE TABLE withdrawal_flow (
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id     text NOT NULL,
    user_id      uuid REFERENCES users NOT NULL,
    amount       numeric(15, 2) NOT NULL DEFAULT 0,
    processed_at timestamp NOT NULL DEFAULT current_timestamp
);
