package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/luxarts/jsend-go"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/utils/sanitizer"
	"net/http"
	"time"
)

type CategoriesService interface {
	Create(req *domain.CategoryCreateRequest, userID string) (*domain.CategoryCreateResponse, error)
	GetAll(userID int64) (*[]domain.CategoryCreateResponse, error)
	GetByName(name string, userID int64) (*domain.CategoryCreateResponse, error)
	GetByEmoji(emoji string, userID int64) (*domain.CategoryCreateResponse, error)
	GetByID(ID int64, userID string) (*domain.CategoryCreateResponse, error)
	DeleteByID(id int64, userID int64) error
	UpdateByID(categoryDTO *domain.CategoryCreateRequest) (*domain.CategoryCreateResponse, error)
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
func (s *categoriesService) GetAll(userID int64) (*[]domain.CategoryCreateResponse, error) {
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
func (s *categoriesService) GetByName(name string, userID int64) (*domain.CategoryCreateResponse, error) {
	name = sanitizer.Sanitize(name)

	category, err := s.repo.GetByName(name, userID)
	if err != nil {
		return nil, err
	}

	return category.ToCategoryCreateResponse(), nil
}
func (s *categoriesService) GetByEmoji(emoji string, userID int64) (*domain.CategoryCreateResponse, error) {
	category, err := s.repo.GetByEmoji(emoji, userID)
	if err != nil {
		return nil, err
	}

	return category.ToCategoryCreateResponse(), nil
}
func (s *categoriesService) DeleteByID(id int64, userID int64) error {
	return s.repo.DeleteByID(id, userID)
}
func (s *categoriesService) UpdateByID(categoryDTO *domain.CategoryCreateRequest) (*domain.CategoryCreateResponse, error) {
	category := categoryDTO.ToCategory()

	category.SanitizedName = sanitizer.Sanitize(category.Name)

	category, err := s.repo.UpdateByID(category)
	if err != nil {
		return nil, err
	}

	return category.ToCategoryCreateResponse(), nil
}
func (s *categoriesService) GetByID(ID int64, userID string) (*domain.CategoryCreateResponse, error) {
	category, err := s.repo.GetByID(ID, userID)
	if err != nil {
		return nil, err
	}

	return category.ToCategoryCreateResponse(), nil
}
