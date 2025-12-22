-- Database schema for the signature service

CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY,
    label VARCHAR(255) NOT NULL,
    algorithm VARCHAR(10) NOT NULL CHECK (algorithm IN ('RSA', 'ECC')),
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,
    signature_counter INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'deactivated')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_devices_created_at ON devices(created_at);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    device_id UUID NOT NULL REFERENCES devices(id),
    counter INTEGER NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    data TEXT NOT NULL,
    signature TEXT NOT NULL,
    signed_data TEXT NOT NULL,
    previous_signature TEXT NOT NULL
);

CREATE INDEX idx_transactions_device_id ON transactions(device_id);
CREATE INDEX idx_transactions_timestamp ON transactions(timestamp);
CREATE INDEX idx_transactions_counter ON transactions(device_id, counter);
