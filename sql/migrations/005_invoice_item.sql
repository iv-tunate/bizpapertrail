-- +goose Up
CREATE TABLE invoice_item(
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id       UUID NOT NULL REFERENCES invoice(id) ON DELETE CASCADE,
    name             TEXT NOT NULL,
    quantity         NUMERIC(15,2) NOT NULL,
    unit             TEXT,                     -- kg, cm, litres, pcs, hrs, etc.
    unit_price       NUMERIC(15,2) NOT NULL,
    total            NUMERIC(15,2) NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE invoice_item