-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_name VARCHAR(50) NOT NULL UNIQUE,
    telegram_id INT NOT NULL UNIQUE,
    points INT DEFAULT 100,
    is_subscribed BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS subscribed_users(
    id UUID PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL REFERENCES users(user_name),
    telegram_id INT NOT NULL REFERENCES users(telegram_id),
    telegram_charge_id TEXT NOT NULL UNIQUE,
    provider_charge_id TEXT NOT NULL UNIQUE,
    subscription_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS subscribed_users;
DROP TABLE IF EXISTS users;
