package room

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	db "github.com/imide/debashare-go/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// Room represents a file sharing room
type Room struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Files     []File    `json:"files"`
}

// File represents a file in a room
type File struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	UploadedAt time.Time `json:"uploadedAt"`
}

// DBFile represents a file as stored in the database
type DBFile struct {
	ID         string           `json:"id"`
	RoomID     string           `json:"room_id"`
	Name       string           `json:"name"`
	Size       int64            `json:"size"`
	UploadedAt pgtype.Timestamp `json:"uploaded_at"`
}

// Manager handles room operations
type Manager struct {
	queries *db.Queries
}

// NewManager creates a new room manager
func NewManager(queries *db.Queries) *Manager {
	return &Manager{
		queries: queries,
	}
}

// CreateRoom creates a new room and returns it
func (m *Manager) CreateRoom(ctx context.Context) (*Room, error) {
	roomID := uuid.New().String()
	now := time.Now()

	err := m.queries.CreateRoom(ctx, db.CreateRoomParams{
		ID: pgtype.Text{
			String: roomID,
			Valid:  true,
		}.String,
		CreatedAt: pgtype.Timestamp{
			Time:  now,
			Valid: true,
		},
	})

	if err != nil {
		return nil, err
	}

	return &Room{
		ID:        roomID,
		CreatedAt: now,
		Files:     make([]File, 0),
	}, nil
}

// GetRoom retrieves a room by ID
func (m *Manager) GetRoom(ctx context.Context, id string) (*Room, error) {
	roomWithFiles, err := m.queries.GetRoomWithFiles(ctx, id)
	if err != nil {
		return nil, err
	}

	var dbFiles []DBFile
	if roomWithFiles.Files != nil {
		if err := json.Unmarshal(roomWithFiles.Files, &dbFiles); err != nil {
			return nil, err
		}
	}

	files := make([]File, len(dbFiles))
	for i, f := range dbFiles {
		files[i] = File{
			ID:         f.ID,
			Name:       f.Name,
			Size:       f.Size,
			UploadedAt: f.UploadedAt.Time,
		}
	}

	return &Room{
		ID:        roomWithFiles.ID,
		CreatedAt: roomWithFiles.CreatedAt.Time,
		Files:     files,
	}, nil
}

// AddFile adds a file to a room
func (m *Manager) AddFile(ctx context.Context, roomID string, name string, size int64) (*File, error) {
	fileID := uuid.New().String()
	now := time.Now()

	err := m.queries.AddFile(ctx, db.AddFileParams{
		ID: pgtype.Text{
			String: fileID,
			Valid:  true,
		}.String,
		RoomID: pgtype.Text{
			String: roomID,
			Valid:  true,
		}.String,
		Name: pgtype.Text{
			String: name,
			Valid:  true,
		}.String,
		Size: pgtype.Int8{
			Int64: size,
			Valid: true,
		}.Int64,
		UploadedAt: pgtype.Timestamp{
			Time:  now,
			Valid: true,
		},
	})

	if err != nil {
		return nil, err
	}

	return &File{
		ID:         fileID,
		Name:       name,
		Size:       size,
		UploadedAt: now,
	}, nil
}

// RemoveFile removes a file from a room
func (m *Manager) RemoveFile(ctx context.Context, roomID string, fileID string) error {
	return m.queries.DeleteFile(ctx, db.DeleteFileParams{
		ID:     fileID,
		RoomID: roomID,
	})
}

// GetFiles returns all files in a room
func (m *Manager) GetFiles(ctx context.Context, roomID string) ([]File, error) {
	dbFiles, err := m.queries.ListFiles(ctx, roomID)
	if err != nil {
		return nil, err
	}

	files := make([]File, len(dbFiles))
	for i, f := range dbFiles {
		files[i] = File{
			ID:         f.ID,
			Name:       f.Name,
			Size:       f.Size,
			UploadedAt: f.UploadedAt.Time,
		}
	}

	return files, nil
}
