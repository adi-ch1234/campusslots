package handler
import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/adi-ch1234/campusslots/internal/domain"
)

// WorkshopService defines the business operations needed by WorkshopHandler.
type WorkshopService interface {
	List(ctx context.Context) ([]domain.Workshop, error)
	GetByID(ctx context.Context, id int64) (domain.Workshop, error)
	Create(ctx context.Context, input domain.CreateWorkshopInput) (domain.Workshop, error)
}
type WorkshopHandler struct {
	svc    WorkshopService
	logger *slog.Logger
}

func NewWorkshopHandler(svc WorkshopService, logger *slog.Logger) *WorkshopHandler {
	return &WorkshopHandler{
		svc:    svc,
		logger: logger,
	}
}
type createWorkshopRequest struct {
	Title    string    `json:"title"`
	Location string    `json:"location"`
	StartsAt time.Time `json:"starts_at"`
	Capacity int       `json:"capacity"`
	Status   string    `json:"status"`
}
func(h *WorkshopHandler) List(w http.ResponseWriter, r * http.Request){
	workshops,err:= h.svc.List(r.Context())
	if err!=nil{
		h.handleError(w,r,err)
		return
	}
	WriteJSON(w,http.StatusOK,map[string]any{
		"workshops":workshops,
	})
}
func (h *WorkshopHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "workshop id must be a positive integer")
		return
	}

	workshop, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, workshop)
}
func (h* WorkshopHandler) Create(w http.ResponseWriter, r * http.Request){
	r.Body = http.MaxBytesReader(w,r.Body,1<<20)
	dec:=json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req createWorkshopRequest
	if err:= dec.Decode(&req);err!=nil{
		WriteError(w,http.StatusBadRequest,"INVALID_REQUEST","Invalid JSON request body")
		return
	}
	if err:=dec.Decode(&struct{}{});!errors.Is(err,io.EOF){
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "request body must contain a single JSON object")
		return
	}
	workshop,err:=h.svc.Create(r.Context(),domain.CreateWorkshopInput{
		Title:req.Title,
		Location:req.Location,
		StartsAt: req.StartsAt,
		Capacity: req.Capacity,
		Status:   req.Status,
	})
	if err!=nil{
		h.handleError(w,r,err)
		return
	}
	WriteJSON(w,http.StatusCreated,workshop)
}
func (h *WorkshopHandler) handleError(w http.ResponseWriter, r *http.Request, err error){
	var valErr domain.ValidationError
	if errors.As(err,&valErr){
		WriteError(w,http.StatusBadRequest,"INVALID_REQUEST",valErr.Message)
		return
	}
	if errors.Is(err,domain.ErrWorkshopNotFound){
		WriteError(w, http.StatusNotFound, "WORKSHOP_NOT_FOUND", "workshop not found")
		return
	}
	h.logger.Error("workshop handler failed",
		"method",r.Method,
		"path",r.URL.Path,
		"error",err,
	)
	WriteError(w, http.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "service temporarily unavailable")
}
