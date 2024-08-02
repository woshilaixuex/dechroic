package test

import (
	"fmt"
	"testing"

	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/color"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: logx测试（文档一坨）
 * @Date: 2024-07-28 20:51
 */
func TestLogxWithColor(t *testing.T) {
	coloredText := color.WithColor("警告", color.FgYellow)
	fmt.Println(coloredText)
}
func TestGetSmallestUnitIncrement(t *testing.T) {
	a := common.NewBigFloat().SetPrec(40).SetFloat64(0.121)
	t.Log(common.GetSmallestUnitIncrementByStr(a))
}
