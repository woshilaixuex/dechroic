package tr_http

import (
	"context"
	"net/http"

	"github.com/delyr1c/dechoric/src/api/dto"
	infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"
	http_model "github.com/delyr1c/dechoric/src/types/model"
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

// 注册用户
func (uc *UserService) Register(c *gin.Context) {
	var req dto.UserRegisterRequstDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, http_model.Response[any]{
			Code: "400",
			Msg:  "请求体有误",
		})
		return
	}
	ctx := context.Background()
	userVO, err := uc.UserRepo.UserRegister(ctx, req.UserName, req.PassWord, req.Email)
	if err != nil {
		c.JSON(http.StatusOK, http_model.Response[any]{
			Code: "500",
			Msg:  "重复注册或请求有误，请检查用户名与邮箱",
		})
		return
	}

	// 构建成功响应
	c.JSON(http.StatusOK, http_model.Response[dto.UserRegisterResponceDTO]{
		Code: "200",
		Msg:  "注册成功，新用户赠送1000积分！",
		Data: dto.UserRegisterResponceDTO{
			UserId:   userVO.UserId,
			Username: userVO.Username,
		},
	})
}

// 用户登录
func (uc *UserService) Login(c *gin.Context) {
	var req dto.UserLoginRequstDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, http_model.Response[any]{
			Code: "400",
			Msg:  err.Error(),
		})
		return
	}

	ctx := context.Background()
	userVO, err := uc.UserRepo.UserLogin(ctx, req.Email, req.PassWord)
	if err != nil {
		c.JSON(http.StatusOK, http_model.Response[any]{
			Code: "401",
			Msg:  "登录失败，用户名或密码错误",
		})
		return
	}

	// 构建成功响应
	c.JSON(http.StatusOK, http_model.Response[dto.UserLoginResponceDTO]{
		Code: "200",
		Msg:  "登录成功",
		Data: dto.UserLoginResponceDTO{
			UserId:   userVO.UserId,
			Username: userVO.Username,
		},
	})
}

// 获取用户信息
func (uc *UserService) GetUserInfo(c *gin.Context) {
	// 从路径参数获取 userId
	userId := c.Query("user_id")

	// 检查 userId 是否为空
	if userId == "" {
		c.JSON(http.StatusOK, http_model.Response[any]{
			Code: "400",
			Msg:  "userId 不能为空",
		})
		return
	}

	ctx := context.Background()
	userVO, err := uc.UserRepo.UserGetInfo(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, http_model.Response[any]{
			Code: "500",
			Msg:  "获取用户信息失败",
		})
		return
	}

	// 构建成功响应
	c.JSON(http.StatusOK, http_model.Response[dto.UserGetInfoResponceDTO]{
		Code: "200",
		Msg:  "获取用户信息成功",
		Data: dto.UserGetInfoResponceDTO{
			UserId:   userVO.UserId,
			Username: userVO.Username,
			Email:    userVO.Email,
			Credit:   userVO.Credit,
		},
	})
}
