CREATE TABLE auto_payments(
                             payment_id BIGSERIAL PRIMARY KEY,
                             name TEXT NOT NULL,
                             number TEXT NOT NULL,
                             timecreated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             timeupdated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);