package user_repository

import (
	"context"

	user_vo "github.com/delyr1c/dechoric/src/domain/usercenter/model/vo"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-04 21:22
 */
type IUserRepository interface {
	UserRegister(ctx context.Context, username, password, email string) (*user_vo.UserTypeVO, error)
	UserLogin(ctx context.Context, email, password string) (*user_vo.UserTypeVO, error)
}
type UserRepository struct {
	repo IUserRepository
}

func NewUserRepository(repo IUserRepository) *UserRepository {
	return &UserRepository{
		repo: repo,
	}
}
func (r *UserRepository) UserRegister(ctx context.Context, username, password, email string) (*user_vo.UserTypeVO, error) {
	return r.repo.UserRegister(ctx, username, password, email)
}
func (r *UserRepository) UserLogin(ctx context.Context, email, password string) (*user_vo.UserTypeVO, error) {
	return r.repo.UserLogin(ctx, email, password)
}
