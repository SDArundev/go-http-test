-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: SearchUsers :many
SELECT * FROM users
WHERE username ILIKE '%' || sqlc.arg(username) || '%';

-- name: Register :exec
INSERT INTO users (username, password)
VALUES (sqlc.arg(username), crypt(sqlc.arg(password), gen_salt('bf')));

-- name: Login :one
SELECT * FROM users
WHERE username = sqlc.arg(username) AND password = crypt(sqlc.arg(password), password);

-- name: VerifyIfUserIsAdmin :one
SELECT is_admin FROM users
WHERE username = sqlc.arg(username) AND is_admin = TRUE;
