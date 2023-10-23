package usecase

import (
	"context"
	"github.com/linqcod/user-service/internal/domain"
	"time"
)

type userUsecase struct {
	userRepository        domain.UserRepository
	ageApiUsecase         domain.AgeApiUsecase
	genderApiUsecase      domain.GenderApiUsecase
	nationalityApiUsecase domain.NationalityApiUsecase
	contextTimeout        time.Duration
}

func NewUserUsecase(
	userRepository domain.UserRepository,
	ageApiUC domain.AgeApiUsecase,
	genderApiUC domain.GenderApiUsecase,
	nationalityApiUC domain.NationalityApiUsecase,
	timeout time.Duration,
) domain.UserUsecase {
	return &userUsecase{
		userRepository:        userRepository,
		ageApiUsecase:         ageApiUC,
		genderApiUsecase:      genderApiUC,
		nationalityApiUsecase: nationalityApiUC,
		contextTimeout:        timeout,
	}
}

func (u userUsecase) Create(ctx context.Context, user *domain.User) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	age, err := u.ageApiUsecase.Get(*user.Name)
	if err != nil {
		return err
	}
	user.Age = &age

	gender, err := u.genderApiUsecase.Get(*user.Name)
	if err != nil {
		return err
	}
	user.Gender = &gender

	nationality, err := u.nationalityApiUsecase.Get(*user.Name)
	if err != nil {
		return err
	}
	user.Nationality = &nationality

	return u.userRepository.Create(c, user)
}

func (u userUsecase) Delete(ctx context.Context, id int64) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepository.Delete(c, id)
}

func (u userUsecase) Update(ctx context.Context, user *domain.User) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepository.Update(c, user)
}

func (u userUsecase) GetById(ctx context.Context, id int64) (*domain.User, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepository.GetById(c, id)
}

func (u userUsecase) GetFilteredUsers(ctx context.Context, count string, nationality string, minAge string, maxAge string, gender string) ([]*domain.User, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepository.GetFilteredUsers(c, count, nationality, minAge, maxAge, gender)
}
