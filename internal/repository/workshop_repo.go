package repository 
import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adi-ch1234/campusslots/internal/domain"
)

type WorkshopRepository struct{
	pool *pgxpool.Pool
}
func NewWorkshopRepository (pool *pgxpool.Pool) *WorkshopRepository {
	return &WorkshopRepository{pool: pool}
}

func (r *WorkshopRepository) List(ctx context.Context) ([]domain.Workshop, error){
	const query = `
		SELECT id, title, location, starts_at, status, capacity, available_seats, created_at 
		FROM workshops
		ORDER BY starts_at ASC , id ASC
	`
	rows,err:= r.pool.Query(ctx,query)
	if err!=nil{
		return nil, fmt.Errorf("query workshops: %w", err)
	}
	defer rows.Close()
	workshops:=make([]domain.Workshop,0)
	for rows.Next(){
		var w domain.Workshop
		err:=rows.Scan(
			&w.ID,
			&w.Title,
			&w.Location,
			&w.StartsAt,
			&w.Status,
			&w.Capacity,
			&w.AvailableSeats,
			&w.CreatedAt,
		)
		if err!=nil{
			return nil, fmt.Errorf("scan workshop row: %w", err)
		}
		workshops = append(workshops,w)
	}
	if err:=rows.Err();err!=nil{
		return nil, fmt.Errorf("iterate workshop rows: %w", err)
	}
	return workshops,nil
}
func (r *WorkshopRepository) GetByID(ctx context.Context, id int64) (domain.Workshop, error) {
	const query = `
		SELECT id, title, location, starts_at, status, capacity, available_seats, created_at
		FROM workshops
		WHERE id = $1
	`

	var w domain.Workshop
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&w.ID,
		&w.Title,
		&w.Location,
		&w.StartsAt,
		&w.Status,
		&w.Capacity,
		&w.AvailableSeats,
		&w.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Workshop{}, domain.ErrWorkshopNotFound
		}
		return domain.Workshop{}, fmt.Errorf("query workshop %d: %w", id, err)
	}

	return w, nil
}
func (r *WorkshopRepository) Create(ctx context.Context , input domain.CreateWorkshopInput) (domain.Workshop,error){
	const query = `
		INSERT INTO workshops (title, location, starts_at, status, capacity, available_seats)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id, title, location, starts_at, status, capacity, available_seats, created_at
	`
	var w domain.Workshop
	err:=r.pool.QueryRow(ctx,query,input.Title,
		input.Location,
		input.StartsAt,
		input.Status,
		input.Capacity,
	).Scan(
		&w.ID,
		&w.Title,
		&w.Location,
		&w.StartsAt,
		&w.Status,
		&w.Capacity,
		&w.AvailableSeats,
		&w.CreatedAt,
	)
	if err != nil {
		return domain.Workshop{}, fmt.Errorf("insert workshop: %w", err)
	}

	return w, nil
}