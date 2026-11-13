package audit

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func From(c *gin.Context, action string) Entry {
	return Entry{
		ActorID:    c.GetString("user_id"),
		ActorRole:  c.GetString("role"),
		FacilityID: c.GetString("facility_id"),
		Action:     action,
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
	}
}

func (s *Service) Log(ctx context.Context, e Entry) error {
	err := s.repo.Insert(ctx, e)
	if err != nil {
		log.Printf("audit log failed for %s: %v", e.Action, err)
	}

	return err
}

func (s *Service) Find(ctx context.Context, f Filter) ([]Record, error) {
	return s.repo.Find(ctx, f)
}

func (s *Service) ForPatientUser(
	ctx context.Context,
	userID string,
) ([]Record, error) {
	return s.repo.ForPatientUser(ctx, userID)
}
