package tr_http

import infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: http方法使用
 * @Date: 2024-10-04 15:55
 */
type RaffeService struct {
	RaffleRepo *infra_repository.StrategyRepository
}

func NewRaffeService(raffleRepo *infra_repository.StrategyRepository) *RaffeService {
	return &RaffeService{
		RaffleRepo: raffleRepo,
	}
}
