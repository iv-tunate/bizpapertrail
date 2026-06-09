-- +goose Up

CREATE TABLE invoice(
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id               UUID REFERENCES users(id) ON DELETE CASCADE,
    invoice_number        TEXT UNIQUE NOT NULL,
    business_name         TEXT NOT NULL,
    customer_name         TEXT NOT NULL,
    customer_email        TEXT,
    customer_phone        TEXT NOT NULL,
    due_date              DATE NOT NULL,
    discount              NUMERIC(15,2),
    tax                   NUMERIC(15,2),
    notes                 TEXT,
    delivery_address      TEXT,
    po_number             TEXT,
    status                TEXT DEFAULT 'draft',  -- draft, sent, paid, cancelled
    receipt_generated     BOOLEAN DEFAULT false,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE invoice;