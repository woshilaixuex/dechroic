package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 自定义基础类型封装
 * @Date: 2024-06-25 20:48
 */

// 高精度浮点
type BigFloat big.Float

func (f *BigFloat) String() string {
	return (*big.Float)(f).String()
}
func (f *BigFloat) Set(value *big.Float) {
	(*big.Float)(f).Set(value)
}

// sql/database基本类型接口实现
func (f *BigFloat) Scan(value interface{}) error {
	if value == nil {
		return nil
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
	if f == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(f.String())
}
func (f *BigFloat) UnmarshalJSON(data []byte) error {
	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	switch v := value.(type) {
	case nil:
		f = nil
	case string:
		floatVal, _, err := big.ParseFloat(v, 10, 256, big.ToNearestEven)
		if err != nil {
			return err
		}
		f.Set(floatVal)
	case float64:
		(*big.Float)(f).SetFloat64(v)
	default:
		return errors.New("invalid type for BigFloat")
	}
	return nil
}
