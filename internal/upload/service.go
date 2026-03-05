package upload

import (
	"time"

	"github.com/google/uuid"
)

// Session represents the state of a file upload session.
// It mirrors the JSON contract currently used by the HTTP API and CLI.
type Session struct {
	ID             string    `json:"uploadid"`
	Filename       string    `json:"filename"`
	TotalSize      int64     `json:"totalsize"`
	ChunkSize      int64     `json:"chunksize"`
	TotalChunks    int64     `json:"totalchunks"`
	ReceivedChunks int64     `json:"receivedchunks"`
	VerifiedChunks int64     `json:"verifiedchunks"`
	Status         string    `json:"status"`
	OverallSHA256  string    `json:"overallSha256,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Service defines the behavior for managing upload sessions.
type Service interface {
	InitiateSession(s Session) Session
}

type service struct{}

// NewService creates a new upload service instance.
func NewService() Service {
	return &service{}
}

// InitiateSession initializes a new upload session with identifiers and timestamps.
func (s *service) InitiateSession(ses Session) Session {
	ses.ID = uuid.New().String()
	ses.Status = "INITIATED"
	now := time.Now()
	ses.CreatedAt = now
	if ses.UpdatedAt.IsZero() {
		ses.UpdatedAt = now
	}
	return ses
}

