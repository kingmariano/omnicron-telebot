-- name: CreateUser :one
INSERT INTO users (id, user_name, telegram_id)
VALUES (
  $1,
  $2,
  $3
)
RETURNING *;
--

-- name: GetUserByUsername :one
SELECT * FROM users 
WHERE user_name = $1;
--
-- name: UpdateUserPoints :one
UPDATE users SET points = $1
WHERE telegram_id = $1
RETURNING *;
--

-- name: UpdateUserSubscriptionStatus :exec
UPDATE users SET is_subscribed = TRUE, points = $1 
WHERE  telegram_id = $2;
--
-- name: AddUserToSubscription :exec
INSERT INTO subscribed_users(id, user_name, telegram_id, telegram_charge_id, provider_charge_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);
--