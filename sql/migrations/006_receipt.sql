-- +goose Up
CREATE TABLE receipt(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id       UUID UNIQUE NOT NULL REFERENCES invoice(id) ON DELETE CASCADE, 
    customer_name    TEXT NOT NULL,
    customer_email   TEXT,
    customer_phone   TEXT NOT NULL,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount_paid      NUMERIC(15,2) NOT NULL,
    payment_method   TEXT NOT NULL,           -- bank transfer, cash, POS, etc.
    payment_date     TIMESTAMPTZ NOT NULL,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE receipt;