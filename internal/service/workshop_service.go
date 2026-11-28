package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/adi-ch1234/campusslots/internal/domain"
)

type WorkshopRepository interface {
	List(ctx context.Context) ([]domain.Workshop, error)
	GetByID(ctx context.Context, id int64) (domain.Workshop, error)
	Create(ctx context.Context, input domain.CreateWorkshopInput) (domain.Workshop, error)
}

// WorkshopCache defines optional read-cache operations for workshops.
type WorkshopCache interface {
	GetList(ctx context.Context) ([]domain.Workshop, bool, error)
	SetList(ctx context.Context, workshops []domain.Workshop) error
	GetByID(ctx context.Context, id int64) (domain.Workshop, bool, error)
	SetByID(ctx context.Context, w domain.Workshop) error
	Invalidate(ctx context.Context, workshopID int64) error
}

type WorkshopService struct {
	repo   WorkshopRepository
	cache  WorkshopCache
	logger *slog.Logger
}

func NewWorkshopService(repo WorkshopRepository, cache WorkshopCache, logger *slog.Logger) *WorkshopService {
	return &WorkshopService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *WorkshopService) List(ctx context.Context) ([]domain.Workshop, error) {
	if s.cache != nil {
		cached, hit, err := s.cache.GetList(ctx)
		if err != nil {
			s.logger.Warn("workshop cache read failed; falling back to database", "error", err)
		} else if hit {
			return cached, nil
		}
	}
	workshops, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	if s.cache != nil {
		if err := s.cache.SetList(ctx, workshops); err != nil {
			s.logger.Warn("Workshop cache right failed", "error", err)
		}
	}
	return workshops, nil
}
func (s *WorkshopService) GetByID(ctx context.Context, id int64) (domain.Workshop, error) {
	if id <= 0 {
		return domain.Workshop{}, domain.ValidationError{Message: "workshop id must be positive"}
	}

	if s.cache != nil {
		cached, hit, err := s.cache.GetByID(ctx, id)
		if err != nil {
			s.logger.Warn("workshop detail cache read failed; falling back to database",
				"workshop_id", id,
				"error", err,
			)
		} else if hit {
			return cached, nil
		}
	}

	workshop, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Workshop{}, err
	}

	if s.cache != nil {
		if err := s.cache.SetByID(ctx, workshop); err != nil {
			s.logger.Warn("workshop detail cache write failed",
				"workshop_id", id,
				"error", err,
			)
		}
	}

	return workshop, nil
}

func (s *WorkshopService) Create(ctx context.Context, input domain.CreateWorkshopInput) (domain.Workshop, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Location = strings.TrimSpace(input.Location)
	input.Status = strings.ToUpper(strings.TrimSpace(input.Status))
	if input.Title == "" {
		return domain.Workshop{}, domain.ValidationError{Message: "title is required"}
	}
	if input.Location == "" {
		return domain.Workshop{}, domain.ValidationError{Message: "location is required"}
	}
	if input.StartsAt.IsZero() {
		return domain.Workshop{}, domain.ValidationError{Message: "starts_at is required"}
	}
	if input.Capacity <= 0 {
		return domain.Workshop{}, domain.ValidationError{Message: "capacity must be greater than 0"}
	}
	if input.Status == "" {
		input.Status = domain.WorkshopStatusOpen
	}
	if input.Status != domain.WorkshopStatusOpen && input.Status != domain.WorkshopStatusClosed {
		return domain.Workshop{}, domain.ValidationError{Message: "status must be OPEN or CLOSED"}
	}

	created, err := s.repo.Create(ctx, input)
	if err != nil {
		return domain.Workshop{}, err
	}
	s.InvalidateWorkshop(ctx, created.ID)
	return created, nil
}
func (s *WorkshopService) InvalidateWorkshop(ctx context.Context, workshopID int64) {
	if s.cache == nil {
		return
	}
	if err := s.cache.Invalidate(ctx, workshopID); err != nil {
		s.logger.Warn("workshop cache invalidation failed",
			"workshop_id", workshopID,
			"error", err,
		)
	}
}
