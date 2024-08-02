package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/delyr1c/dechoric/src/types/cerr"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 自定义高精度浮点
 * @Date: 2024-06-25 20:48
 */

// 高精度浮点
type BigFloat struct {
	*big.Float
}

func NewBigFloat() *BigFloat {
	return &BigFloat{new(big.Float)}
}

func (f *BigFloat) Cmp(y *BigFloat) int {
	return f.Float.Cmp(y.Float)
}

// 获取最小浮点单位增量
func GetSmallestUnitIncrementByStr(f *big.Float) (*big.Float, error) {
	increment := new(big.Float).SetPrec(f.Prec())
	nums, err := f.MarshalText()
	if err != nil {
		return nil, cerr.LogError(errors.New("the big.Float`s formatis err"))
	}
	newnums := make([]byte, 0)
	isFound := false
	index := 0
	for _, num := range nums {
		switch num {
		case '.':
			index += 1
			newnums = append(newnums, '0', '.')
			isFound = true
		default:
			if isFound {
				index += 1
				newnums = append(newnums, '0')
			}
		}
	}
	if index < 2 {
		return new(big.Float).SetFloat64(0), cerr.LogError(errors.New("the big.Float`s index less than 2"))
	}
	newnums[index] = '1'
	increment.SetString(string(newnums))
	return increment, nil
}

// sql/database基本类型接口实现
func (f *BigFloat) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if f.Float == nil {
		f.Float = new(big.Float)
	}
	switch v := value.(type) {
	case []byte:
		fVal, _, err := big.ParseFloat(string(v), 10, 512, big.ToNearestEven)
		if err != nil {
			return err
		}
		f.Set(fVal)
		return nil
	case string:
		fVal, _, err := big.ParseFloat(v, 10, 512, big.ToNearestEven)
		if err != nil {
			return err
		}
		f.Set(fVal)
		return nil
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type *big.Float", value)
	}
}
func (f *BigFloat) Value() (driver.Value, error) {
	return f.String(), nil
}

// 序列化接口实现
func (f *BigFloat) MarshalJSON() ([]byte, error) {
	if f.Float == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(f.String())
}
func (f *BigFloat) UnmarshalJSON(data []byte) error {
	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if f.Float == nil {
		f.Float = new(big.Float)
	}
	switch v := value.(type) {
	case nil:
		*f = BigFloat{}
	case []byte:
		floatVal, _, err := big.ParseFloat(string(v), 10, 256, big.ToNearestEven)
		if err != nil {
			return err
		}
		f.Float.Set(floatVal)
	case string:
		floatVal, _, err := big.ParseFloat(v, 10, 256, big.ToNearestEven)
		if err != nil {
			return err
		}
		f.Float.Set(floatVal)
	case float64:
		f.Float.SetFloat64(v)
	default:
		return errors.New("invalid type for BigFloat")
	}
	return nil
}
