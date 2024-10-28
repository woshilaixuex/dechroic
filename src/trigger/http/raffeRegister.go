package tr_http

import (
	"context"
	"net/http"

	"github.com/delyr1c/dechoric/src/api/dto"
	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/raffle"
	infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"
	http_model "github.com/delyr1c/dechoric/src/types/model"
	"github.com/gin-gonic/gin"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: http方法使用
 * @Date: 2024-10-04 15:55
 */
type RaffeService struct {
	RaffleRepo            *infra_repository.StrategyRepository
	Armory                *armory.StrategyArmory
	DefaultRaffleStrategy *raffle.DefaultRaffleStrategy
}

func NewRaffeService(raffleRepo *infra_repository.StrategyRepository) *RaffeService {
	armory := armory.NewStrategyArmory(*repository.NewStrategyService(raffleRepo))
	return &RaffeService{
		RaffleRepo:            raffleRepo,
		Armory:                armory,
		DefaultRaffleStrategy: raffle.NewDefaultRaffleStrategy(*repository.NewStrategyService(raffleRepo), armory),
	}
}

func (rc *RaffeService) GetPrize(c *gin.Context) {
	// 解析请求参数
	var req dto.PrizeRequsetDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, http_model.Response[any]{
			Code: "400",
			Msg:  "请求体有误",
		})
		return
	}

	// 获取 context 和 strategyId
	ctx := context.Background()
	strategyId := req.StrategyId

	// 调用仓库层获取奖品信息
	prizesVO, err := rc.RaffleRepo.QueryPrize(ctx, strategyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, http_model.Response[any]{
			Code: "500",
			Msg:  "后台获取奖品信息错误",
		})
		return
	}

	// 构建响应数据
	prizeResponse := make([]dto.PrizeResponseDTO, len(prizesVO.PrizeVOs))
	for i, prizeVO := range prizesVO.PrizeVOs {
		prizeResponse[i] = dto.PrizeResponseDTO{
			SortId:     prizeVO.SortId,
			AwardId:    prizeVO.AwardId,
			StrategyId: prizeVO.StrategyId,
			AwardTitle: prizeVO.AwardTitle,
		}
	}

	// 返回成功响应
	c.JSON(http.StatusOK, http_model.Response[dto.PrizesResponseDTO]{
		Code: "200",
		Msg:  "获取奖品列表成功",
		Data: dto.PrizesResponseDTO{
			PrizeList: prizeResponse,
		},
	})
}
func (rc *RaffeService) GetAwardPrize(c *gin.Context) {
	var req dto.AwardPrizeRequsetDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, http_model.Response[any]{
			Code: "400",
			Msg:  "请求体有误",
		})
		return
	}
	ctx := context.Background()
	// 检查用户积分
	isOK, _ := rc.RaffleRepo.DoUserCredit(ctx, req.UserId)
	if !isOK {
		c.JSON(http.StatusOK, http_model.Response[any]{
			Code: "200",
			Msg:  "用户积分不足请进行充值",
		})
		return
	}
	raffleAwardEntity := &StrategyEntity.RaffleFactorEntity{
		UserId:     req.UserId,
		StrategyId: req.StrategyId,
	}
	award, err := rc.DefaultRaffleStrategy.PerformRaffle(ctx, raffleAwardEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, http_model.Response[any]{
			Code: "500",
			Msg:  "抽奖发生错误",
		})
		return
	}
	go func() {
		rc.RaffleRepo.PerformAward(ctx, req.UserId, int32(award.AwardId))
	}()
	desc, err := rc.RaffleRepo.AddStrategyAwardHistory(ctx, req.UserId, award)
	if err != nil {
		c.JSON(http.StatusInternalServerError, http_model.Response[any]{
			Code: "500",
			Msg:  "历史添加发生错误",
		})
		return
	}
	c.JSON(http.StatusOK, http_model.Response[dto.AwardPrizeResponseDTO]{
		Code: "200",
		Msg:  "获取奖品列表成功",
		Data: dto.AwardPrizeResponseDTO{
			AwardId:    int32(award.AwardId),
			AwardTitle: desc,
		},
	})
}
func (rc *RaffeService) GetHistory(c *gin.Context) {
	var req dto.GetHistoryRequsetDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, http_model.Response[any]{
			Code: "400",
			Msg:  "请求体有误",
		})
		return
	}
	ctx := context.Background()
	historys, totalPage, err := rc.RaffleRepo.QueryStrategyAwardHistory(ctx, req.UserId, int(req.Page))
	if err != nil {
		c.JSON(http.StatusBadRequest, http_model.Response[any]{
			Code: "400",
			Msg:  "请求内容不合法",
		})
		return
	}
	if req.Page > int32(totalPage) && int32(totalPage) != 0 {
		c.JSON(http.StatusBadRequest, http_model.Response[any]{
			Code: "400",
			Msg:  "请求内容不合法",
		})
		return
	}
	c.JSON(http.StatusOK, http_model.Response[dto.GetHistoryResponseDTO]{
		Code: "200",
		Msg:  "获取历史记录",
		Data: dto.GetHistoryResponseDTO{
			Page:       req.Page,
			TotalPages: int32(totalPage),
			Historys:   historys,
		},
	})
}
