package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/utils/sanitizer"
)

type CategoriesService interface {
	Create(categoryDTO *domain.CategoryDTO) (*domain.CategoryDTO, error)
	GetAll(userID int) (*[]domain.CategoryDTO, error)
	GetByName(name string, userID int) (*domain.CategoryDTO, error)
	GetByEmoji(emoji string, userID int) (*domain.CategoryDTO, error)
	GetByID(ID int, userID int) (*domain.CategoryDTO, error)
	DeleteByID(id int, userID int) error
	UpdateByID(categoryDTO *domain.CategoryDTO) (*domain.CategoryDTO, error)
}
type categoriesService struct {
	repo repository.CategoriesRepository
}

func NewCategoriesService(categoriesRepo repository.CategoriesRepository) CategoriesService {
	return &categoriesService{
		repo: categoriesRepo,
	}
}
func (s *categoriesService) Create(categoryDTO *domain.CategoryDTO) (*domain.CategoryDTO, error) {
	category := categoryDTO.ToCategory()

	category.SanitizedName = sanitizer.Sanitize(category.Name)

	category, err := s.repo.Create(category)
	if err != nil {
		return nil, err
	}

	return category.ToDTO(), nil
}
func (s *categoriesService) GetAll(userID int) (*[]domain.CategoryDTO, error) {
	categories, err := s.repo.GetAll(userID)
	if err != nil {
		return nil, err
	}

	var categoryDTOs = make([]domain.CategoryDTO, 0)
	for _, category := range *categories {
		categoryDTOs = append(categoryDTOs, *category.ToDTO())
	}

	return &categoryDTOs, nil
}
func (s *categoriesService) GetByName(name string, userID int) (*domain.CategoryDTO, error) {
	name = sanitizer.Sanitize(name)

	category, err := s.repo.GetByName(name, userID)
	if err != nil {
		return nil, err
	}

	return category.ToDTO(), nil
}
func (s *categoriesService) GetByEmoji(emoji string, userID int) (*domain.CategoryDTO, error) {
	category, err := s.repo.GetByEmoji(emoji, userID)
	if err != nil {
		return nil, err
	}

	return category.ToDTO(), nil
}
func (s *categoriesService) DeleteByID(id int, userID int) error {
	return s.repo.DeleteByID(id, userID)
}
func (s *categoriesService) UpdateByID(categoryDTO *domain.CategoryDTO) (*domain.CategoryDTO, error) {
	category := categoryDTO.ToCategory()

	category.SanitizedName = sanitizer.Sanitize(category.Name)

	category, err := s.repo.UpdateByID(category)
	if err != nil {
		return nil, err
	}

	return category.ToDTO(), nil
}
func (s *categoriesService) GetByID(ID int, userID int) (*domain.CategoryDTO, error) {
	category, err := s.repo.GetByID(ID, userID)
	if err != nil {
		return nil, err
	}

	return category.ToDTO(), nil
}
