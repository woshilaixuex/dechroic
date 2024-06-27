package common

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 常用常量
 * @Date: 2024-06-25 20:49
 */

const (
	SPLIT     = ","
	COLON     = ":"
	SPACE     = " "
	UNDERLINE = "_"
)

// RedisKeys 结构体包含应用中使用的Redis键前缀。
var RedisKeys = struct {
	ActivityKey                string
	ActivitySKUKey             string
	ActivityCountKey           string
	StrategyKey                string
	StrategyAwardKey           string
	StrategyAwardListKey       string
	StrategyRateTableKey       string
	StrategyRateRangeKey       string
	RuleTreeVOKey              string
	StrategyAwardCountKey      string
	StrategyAwardCountQueryKey string
	StrategyRuleWeightKey      string
	ActivitySKUCountQueryKey   string
	ActivitySKUStockCountKey   string
	ActivitySKUCountClearKey   string
	ActivityAccountLock        string
	ActivityAccountUpdateLock  string
	UserCreditAccountLock      string
}{
	ActivityKey:                "dechoric_activity_key_",
	ActivitySKUKey:             "dechoric_activity_sku_key_",
	ActivityCountKey:           "dechoric_activity_count_key_",
	StrategyKey:                "dechoric_strategy_key_",
	StrategyAwardKey:           "dechoric_strategy_award_key_",
	StrategyAwardListKey:       "dechoric_strategy_award_list_key_",
	StrategyRateTableKey:       "dechoric_strategy_rate_table_key_",
	StrategyRateRangeKey:       "dechoric_strategy_rate_range_key_",
	RuleTreeVOKey:              "rule_tree_vo_key_",
	StrategyAwardCountKey:      "strategy_award_count_key_",
	StrategyAwardCountQueryKey: "strategy_award_count_query_key",
	StrategyRuleWeightKey:      "strategy_rule_weight_key_",
	ActivitySKUCountQueryKey:   "activity_sku_count_query_key",
	ActivitySKUStockCountKey:   "activity_sku_stock_count_key_",
	ActivitySKUCountClearKey:   "activity_sku_count_clear_key_",
	ActivityAccountLock:        "activity_account_lock_",
	ActivityAccountUpdateLock:  "activity_account_update_lock_",
	UserCreditAccountLock:      "user_credit_account_lock_",
}
