package infra_repository

import (
	"context"
	"errors"

	user_vo "github.com/delyr1c/dechoric/src/domain/usercenter/model/vo"
	user_repository "github.com/delyr1c/dechoric/src/domain/usercenter/repository"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/aiModel"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/aiUsage"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/user"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/redis"
	"github.com/delyr1c/dechoric/src/types/cerr"
	"github.com/delyr1c/dechoric/src/types/encrypt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-04 21:27
 */

var _ user_repository.IUserRepository = (*UserRepository)(nil)

type UserRepository struct {
	RedisService redis.RedisService
	UserModel    user.UserModel
	AIModel      aiModel.AiModelsModel
	AiUsageModel aiUsage.AiUsageModel
}

func NewUserRepository(sqlConn sqlx.SqlConn, redis redis.RedisService) *UserRepository {
	return &UserRepository{
		RedisService: redis,
		UserModel:    user.NewUserModel(sqlConn),
		AIModel:      aiModel.NewAiModelsModel(sqlConn),
		AiUsageModel: aiUsage.NewAiUsageModel(sqlConn),
	}
}
func (r *UserRepository) UserRegister(ctx context.Context, username, password, email string) (*user_vo.UserTypeVO, error) {
	// 生成用户ID
	userId, err := r.RedisService.GenerateID(ctx)
	if err != nil {
		return nil, err
	}
	// 密码加密
	password, err = encrypt.Encrypt(password)
	if err != nil {
		return nil, err
	}

	// 创建用户数据并插入到数据库
	userData := &user.User{
		UserId:   userId,
		Username: username,
		Password: password,
		Email:    email,
	}
	_, err = r.UserModel.Insert(ctx, userData)
	if err != nil {
		return nil, cerr.LogError(err)
	}

	// 查询所有 AI 模型
	aiModels, err := r.AIModel.FindAll(ctx)
	if err != nil {
		return nil, cerr.LogError(err)
	}

	// 为每个 AI 模型添加使用记录
	for _, model := range aiModels {
		aiUsageData := &aiUsage.AiUsage{
			UserId:     userId,
			ModelId:    model.ModelId,
			ModelName:  model.ModelName,
			QueryCount: 20, // 默认剩余使用次数为20
		}
		_, err := r.AiUsageModel.Insert(ctx, aiUsageData)
		if err != nil {
			return nil, cerr.LogError(err)
		}
	}
	// 返回注册成功的用户信息
	return &user_vo.UserTypeVO{
		UserId:   userId,
		Username: username,
		Email:    email,
	}, nil
}

func (r *UserRepository) UserLogin(ctx context.Context, email, password string) (*user_vo.UserTypeVO, error) {
	userData, err := r.UserModel.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, cerr.LogError(err)
	}
	isValid, err := encrypt.DeCode(password, userData.Password)
	if err != nil {
		return nil, cerr.LogError(err)
	}
	if !isValid {
		return nil, cerr.LogError(errors.New("invalid username or password"))
	}
	return &user_vo.UserTypeVO{
		UserId:   userData.UserId,
		Username: userData.Username,
		Email:    userData.Email,
	}, nil
}
