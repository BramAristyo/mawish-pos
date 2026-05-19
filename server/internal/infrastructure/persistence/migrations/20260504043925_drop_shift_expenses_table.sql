-- +goose Up
DROP TABLE IF EXISTS shift_expenses;

-- +goose Down
CREATE TABLE shift_expenses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shift_id    UUID NOT NULL REFERENCES shifts(id),
    coa_id      UUID NOT NULL REFERENCES chart_of_accounts(id),
    amount      DECIMAL(12, 2) NOT NULL DEFAULT 0,
    description TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shift_expenses_shift_id ON shift_expenses(shift_id);
CREATE INDEX idx_shift_expenses_coa_id ON shift_expenses(coa_id);
