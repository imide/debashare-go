-- name: CreateRoom :exec
INSERT INTO rooms (id, created_at) VALUES ($1, $2);

-- name: GetRoom :one
SELECT * FROM rooms WHERE id = $1;

-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = $1;

-- name: AddFile :exec
INSERT INTO files (id, room_id, name, size, uploaded_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetFile :one
SELECT * FROM files
WHERE room_id = $1 AND id = $2;

-- name: DeleteFile :exec
DELETE FROM files
WHERE room_id = $1 AND id = $2;

-- name: ListFiles :many
SELECT * FROM files
WHERE room_id = $1
ORDER BY uploaded_at DESC;

-- name: GetRoomWithFiles :one
SELECT r.*, json_agg(f.*) as files
FROM rooms r
         LEFT JOIN files f ON r.id = f.room_id
WHERE r.id = $1
GROUP BY r.id;

-- name: DeleteOldRooms :exec
DELETE FROM rooms
WHERE created_at < $1;

-- name: GetRoomStats :one
SELECT
    COUNT(*) as total_rooms,
    COUNT(DISTINCT f.room_id) as rooms_with_files,
    COUNT(f.id) as total_files,
    COALESCE(SUM(f.size), 0) as total_size
FROM rooms r
         LEFT JOIN files f ON r.id = f.room_id;
