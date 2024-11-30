package postgres

import (
	"context"
	"errors"
	"github.com/minishop/internal/domain"
	"gorm.io/gorm"
	"time"
)

type UserRepository struct {
	db *gorm.DB
}

type user struct {
	ID        uint64    `json:"id" gorm:"primarykey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserRepository) Create(ctx context.Context, dUser domain.UserCreateParameters) (domain.User, error) {
	// Parse domain.User into user
	repoUser := user{
		Username:  dUser.Username,
		Password:  dUser.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := u.db.WithContext(ctx).Create(&repoUser)
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return domain.User{
		ID:        repoUser.ID,
		Username:  repoUser.Username,
		Password:  repoUser.Password,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
	}, nil
}

func (u *UserRepository) GetByUsername(ctx context.Context, s string) (domain.User, error) {
	repoUser := user{}
	result := u.db.WithContext(ctx).First(&repoUser, "username = ?", s)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.NotFoundError
		}

		return domain.User{}, result.Error
	}

	return domain.User{
		ID:        repoUser.ID,
		Username:  repoUser.Username,
		Password:  repoUser.Password,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
	}, nil
}

func (u *UserRepository) GetByID(ctx context.Context, userId uint64) (domain.User, error) {
	repoUser := user{}
	result := u.db.WithContext(ctx).First(&repoUser, userId)
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return domain.User{
		ID:        repoUser.ID,
		Username:  repoUser.Username,
		Password:  repoUser.Password,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
	}, nil
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) AutoMigrate() error {
	return u.db.AutoMigrate(&user{})
}
