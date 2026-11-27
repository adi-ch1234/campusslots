package service
import (
	"context"
	"strings"

	"github.com/adi-ch1234/campusslots/internal/domain"

)
type WorkshopRepository interface {
	List(ctx context.Context) ([]domain.Workshop,error)
	GetByID(ctx context.Context,id int64) (domain.Workshop , error)
	Create(ctx context.Context, input domain.CreateWorkshopInput) (domain.Workshop,error)
}
type WorkshopService struct {
	repo WorkshopRepository
}
func NewWorkshopService(repo WorkshopRepository) *WorkshopService{
	return &WorkshopService{repo:repo}
}
func (s *WorkshopService) List(ctx context.Context) ([]domain.Workshop,error){
	return s.repo.List(ctx)
}
func (s *WorkshopService) GetByID(ctx context.Context ,id int64) (domain.Workshop,error){
	if id<=0 {
		return domain.Workshop{} , domain.ValidationError{Message:"Workshop id must be positive"}
	}
	return s.repo.GetByID(ctx,id)
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

	return s.repo.Create(ctx, input)
}
