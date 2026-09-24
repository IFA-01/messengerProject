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
-- name: FindDirectChatsBetweenUsers :one
SELECT c.*
FROM chats c
WHERE c.is_group = false
  AND EXISTS (
    SELECT 1 FROM chat_members cm1
    WHERE cm1.chat_id = c.id AND cm1.user_id = $1
  )
  AND EXISTS (
    SELECT 1 FROM chat_members cm2
    WHERE cm2.chat_id = c.id AND cm2.user_id = $2
  )
  AND (
    SELECT COUNT(*) FROM chat_members cm3
    WHERE cm3.chat_id = c.id
  ) = 2
LIMIT 1;

-- name: ListUsersChats :many
SELECT c.*
FROM chats c
JOIN chat_members cm ON cm.chat_id = c.id
WHERE cm.user_id = $1
ORDER BY c.updated_at DESC;