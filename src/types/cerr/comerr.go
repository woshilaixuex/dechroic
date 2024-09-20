package cerr

import (
	"fmt"
	"runtime"

	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 统一错误处理（执行过程）
 * @Date: 2024-07-27 12:44
 */

// debug log
func LogError(err error) error {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		logx.Error(fmt.Sprintf("%s: %s:%d: %s", funcName, file, line, err.Error()))
	} else {
		logx.Error(err.Error())
	}
	return err
}
