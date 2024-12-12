package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/imide/debashare-go/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	queries *db.Queries
}

func NewHandler(queries *db.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

func (h *Handler) CreateRoom(c *fiber.Ctx) error {
	// Create room with UUID
	err := h.queries.CreateRoom(c.Context(), db.CreateRoomParams{
		ID: pgtype.Text{
			String: uuid.New().String(),
			Valid:  true,
		}.String,
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create room",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) GetRoom(c *fiber.Ctx) error {
	roomID := c.Params("id")

	room, err := h.queries.GetRoomWithFiles(c.Context(), roomID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Room not found",
		})
	}

	return c.JSON(room)
}

func (h *Handler) AddFile(c *fiber.Ctx) error {
	roomID := c.Params("id")

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get file",
		})
	}

	err = h.queries.AddFile(c.Context(), db.AddFileParams{
		ID: pgtype.Text{
			String: uuid.New().String(),
			Valid:  true,
		}.String,
		RoomID: pgtype.Text{
			String: roomID,
			Valid:  true,
		}.String,
		Name: pgtype.Text{
			String: file.Filename,
			Valid:  true,
		}.String,
		Size: pgtype.Int8{
			Int64: file.Size,
			Valid: true,
		}.Int64,
		UploadedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add file",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) DeleteFile(c *fiber.Ctx) error {
	roomID := c.Params("roomId")
	fileID := c.Params("fileId")

	err := h.queries.DeleteFile(c.Context(), db.DeleteFileParams{
		ID:     fileID,
		RoomID: roomID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) ListFiles(c *fiber.Ctx) error {
	roomID := c.Params("id")

	files, err := h.queries.ListFiles(c.Context(), roomID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list files",
		})
	}

	return c.JSON(files)
}
