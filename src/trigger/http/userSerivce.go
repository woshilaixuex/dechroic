package tr_http

import (
	"context"
	"net/http"

	"github.com/delyr1c/dechoric/src/api/dto"
	infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"
	"github.com/gin-gonic/gin"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-05 09:56
 */
type UserService struct {
	UserRepo *infra_repository.UserRepository
}

func NewUserService(repository *infra_repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: repository,
	}
}

func (uc *UserService) Register(c *gin.Context) {
	var req dto.UserRegisterRequstDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	userVO, err := uc.UserRepo.UserRegister(ctx, req.UserName, req.PassWord, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userVO)
}

func (uc *UserService) Login(c *gin.Context) {
	var req dto.UserLoginRequstDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	userVO, err := uc.UserRepo.UserLogin(ctx, req.Email, req.PassWord)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userVO)
}
