-- name: CreateChat :one
INSERT INTO chats (name, is_group)
VALUES ($1, $2)
RETURNING *;
-- name: GetChatByID :one
SELECT * FROM chats
WHERE id = $1
LIMIT 1;
-- name: AddChatMember :one
INSERT INTO chat_members (chat_id, user_id)
VALUES ($1, $2)
RETURNING *;
-- name: IsChatMember :one
SELECT EXISTS (
    SELECT 1 FROM chat_members
    WHERE chat_id = $1 AND user_id = $2
);
-- name: CreateMessage :one
INSERT INTO messages (chat_id, sender_id, content)
VALUES ($1, $2, $3)
RETURNING *;
-- name: ListChatMessages :many
SELECT * FROM messages
WHERE chat_id = $1
ORDER BY created_at ASC;