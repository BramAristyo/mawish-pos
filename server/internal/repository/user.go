package repository

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) GetAll() ([]domain.User, error) {
	var u []domain.User

	if err := r.DB.Find(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Store(u *domain.User) (domain.User, error) {
	if err := r.DB.Create(u).Error; err != nil {
		return domain.User{}, err
	}

	return *u, nil
}

func (r *UserRepository) FindById(id uuid.UUID) (domain.User, error) {
	var u domain.User

	if err := r.DB.Where("id = ?", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, usecase_errors.NotFound
		}
		return domain.User{}, err
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (domain.User, error) {
	var u domain.User

	if err := r.DB.Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, usecase_errors.EmailNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepository) Update(id uuid.UUID, u *domain.User) (domain.User, error) {
	var existing domain.User
	if err := r.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, usecase_errors.NotFound
		}
		return domain.User{}, err
	}

	updateData := map[string]any{
		"name":  u.Name,
		"email": u.Email,
		"role":  u.Role,
	}

	if err := r.DB.Model(&existing).Updates(updateData).Error; err != nil {
		return domain.User{}, err
	}

	return existing, nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	result := r.DB.Delete(&domain.User{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *UserRepository) Restore(id uuid.UUID) error {
	result := r.DB.
		Model(&domain.User{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *UserRepository) IsEmailExist(email string) (bool, error) {
	var count int64

	if err := r.DB.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepository) IsEmailTaken(id uuid.UUID, email string) (bool, error) {
	var count int64

	if err := r.DB.Model(&domain.User{}).Where("email = ? AND id != ?", email, id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
