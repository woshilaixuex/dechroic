package chain

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 过滤链接口
 * @Date: 2024-08-15 23:10
 */
type ILogiChain interface {
	Logic(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error)
	ModelType() string
	// 不对外提供
	ILogiChainBase
}
type ILogiChainBase interface {
	AppendNext(next ILogiChain) ILogiChain
	Next() ILogiChain
}
