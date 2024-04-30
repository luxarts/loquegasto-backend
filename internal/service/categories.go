package service

import (
	"errors"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/utils/sanitizer"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/luxarts/jsend-go"
)

type CategoriesService interface {
	Create(req *domain.CategoryCreateRequest, userID string) (*domain.CategoryCreateResponse, error)
	GetAll(userID string) (*[]domain.CategoryCreateResponse, error)
	GetByName(name string, userID string) (*domain.CategoryCreateResponse, error)
	GetByEmoji(emoji string, userID string) (*domain.CategoryCreateResponse, error)
	GetByID(ID string, userID string) (*domain.CategoryCreateResponse, error)
	DeleteByID(id string, userID string) error
	UpdateByID(req *domain.CategoryUpdateRequest, id string, userID string) (*domain.CategoryUpdateResponse, error)
}
type categoriesService struct {
	repo repository.CategoriesRepository
}

func NewCategoriesService(categoriesRepo repository.CategoriesRepository) CategoriesService {
	return &categoriesService{
		repo: categoriesRepo,
	}
}
func (s *categoriesService) Create(req *domain.CategoryCreateRequest, userID string) (*domain.CategoryCreateResponse, error) {
	sanitizedName := sanitizer.Sanitize(req.Name)

	// Check if name already exists for the given user
	c, err := s.repo.GetBySanitizedName(sanitizedName, userID)
	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && err != nil && *jsendErr.Code != http.StatusNotFound {
		return nil, err
	}
	if c != nil {
		return nil, jsend.NewError("category name already exists", nil, http.StatusBadRequest)
	}

	category := req.ToCategory()
	category.ID = uuid.NewString()
	category.SanitizedName = sanitizedName
	category.UserID = userID
	category.CreatedAt = time.Now()

	category, err = s.repo.Create(category)
	if err != nil {
		return nil, err
	}

	return category.ToCategoryCreateResponse(), nil
}
func (s *categoriesService) GetAll(userID string) (*[]domain.CategoryCreateResponse, error) {
	categories, err := s.repo.GetAll(userID)
	if err != nil {
		return nil, err
	}

	var categoryDTOs = make([]domain.CategoryCreateResponse, 0)
	for _, category := range *categories {
		categoryDTOs = append(categoryDTOs, *category.ToCategoryCreateResponse())
	}

	return &categoryDTOs, nil
}
func (s *categoriesService) GetByName(name string, userID string) (*domain.CategoryCreateResponse, error) {
	name = sanitizer.Sanitize(name)

	category, err := s.repo.GetByName(name, userID)
	if err != nil {
		return nil, err
	}

	return category.ToCategoryCreateResponse(), nil
}
func (s *categoriesService) GetByEmoji(emoji string, userID string) (*domain.CategoryCreateResponse, error) {
	category, err := s.repo.GetByEmoji(emoji, userID)
	if err != nil {
		return nil, err
	}

	return category.ToCategoryCreateResponse(), nil
}
func (s *categoriesService) DeleteByID(id string, userID string) error {
	return s.repo.DeleteByID(id, userID)
}
func (s *categoriesService) UpdateByID(req *domain.CategoryUpdateRequest, id string, userID string) (*domain.CategoryUpdateResponse, error) {
	sanitizedName := sanitizer.Sanitize(req.Name)

	// Check if the name already exists for the given user
	c, err := s.repo.GetBySanitizedName(sanitizedName, userID)
	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && err != nil && *jsendErr.Code != http.StatusNotFound {
		return nil, err
	}
	if c != nil {
		return nil, jsend.NewError("category name already used", nil, http.StatusBadRequest)
	}

	category := req.ToCategory()

	category, err = s.repo.UpdateByID(category, id, userID)
	if err != nil {
		return nil, err
	}

	return category.ToCategoryUpdateResponse(), nil
}
func (s *categoriesService) GetByID(ID string, userID string) (*domain.CategoryCreateResponse, error) {
	category, err := s.repo.GetByID(ID, userID)
	if err != nil {
		return nil, err
	}

	return category.ToCategoryCreateResponse(), nil
}
