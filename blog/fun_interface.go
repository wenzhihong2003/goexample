package main

import (
	"math"
	"strings"
	"strconv"
	"fmt"
)

// https://my.oschina.net/henrylee2cn/blog/741332
// 如何用函数实现接口以及如何检验接口实现
// Golang中下划线“_”表示忽略接收到的值；
//
// const、var、type关键字均支持分组形式，以圆括号“()”包裹，
// 建议将相关声明写在同一分组，如上面代码中Handler和HandlerFunc的声明。
type (
	Handler interface {
		Do(int) error
	}

	HandlerFunc func(int) error
)

func (hf HandlerFunc) Do(i int) error {
	return hf(i)
}

var _ Handler = HandlerFunc(nil) // 这里是将 nil 转为 HandlerFunc 类型

// 参考: https://my.oschina.net/henrylee2cn/blog/742764
// 浮点数比较大小. 展示了Go语言函数式编程的独特魅力

type Accuracy func() float64

func (this Accuracy) Equal(a, b float64) bool {
	return math.Abs(a-b) < this()
}

func (this Accuracy) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > this()
}

func (this Accuracy) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) < this()
}

func (this Accuracy) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < this()
}

func (this Accuracy) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < this()
}

// 一个较完整的处理浮点数的结构体——Floater

type Floater struct {
	numOfDecimalPlaces int
	accuracy           float64
	format             string
}

func NewFloater(numOfDecimalPlaces int) *Floater {
	if numOfDecimalPlaces < 0 || numOfDecimalPlaces > 14 {
		panic("the range of Floater.numOfDecimalPlaces must be between 0 and 14.")
	}
	var accuracy float64 = 1
	if numOfDecimalPlaces > 0 {
		accuracyString := "0." + strings.Repeat("0", numOfDecimalPlaces) + "1"
		accuracy, _ = strconv.ParseFloat(accuracyString, 64)
	}
	return &Floater{
		numOfDecimalPlaces: numOfDecimalPlaces,
		accuracy:           accuracy,
		format:             "%0." + strconv.Itoa(numOfDecimalPlaces) + "f",
	}
}

func (this *Floater) NumofDecimalPlaces() int {
	return this.numOfDecimalPlaces
}

func (this *Floater) Accuracy() float64 {
	return this.accuracy
}

func (this *Floater) Format() string {
	return fmt.Sprintf(this.format, f)
}

func (this *Floater) Atof(s string, bitSize int) (float64, error) {
	f, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return f, err
	}

	return strconv.ParseFloat(fmt.Sprintf(this.format, f), bitSize)
}

func (this *Floater) Ftof(f float64) float64 {
	f, _ = strconv.ParseFloat(fmt.Sprintf(this.format, f), 64)
	return f
}

func (this *Floater) Atoa(s string, bitSize int) (string, error) {
	f, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return s, err
	}
	return fmt.Sprintf(this.format, f), nil
}

func (this *Floater) Equal(a, b float64) bool {
	return math.Abs(a-b) < this.accuracy
}

func (this *Floater) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > this.accuracy
}

func (this *Floater) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > this.accuracy
}

func (this *Floater) GreaterOrEquan(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < this.accuracy
}

func (this *Floater) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < this.accuracy
}

// 使用go build 进行条件编译
// 参见:https://my.oschina.net/henrylee2cn/blog/821146

// Docker的安装配置及使用详解
// https://my.oschina.net/henrylee2cn/blog/821532

// LOG日志级别
// https://my.oschina.net/henrylee2cn/blog/823942

// Golang 平滑关闭／重启与热编译技术
// https://my.oschina.net/henrylee2cn/blog/869059

// 常用 Git 命令清单
// https://my.oschina.net/henrylee2cn/blog/871452

// 使用git-flow来帮助管理git代码
// https://my.oschina.net/henrylee2cn/blog/871403

// gitflow 开发流程
// https://my.oschina.net/henrylee2cn/blog/871407

// 基于SourceTree 下的 Git Flow 模型
// http://blog.haohtml.com/archives/16039

// Git 工作流程
// http://www.ruanyifeng.com/blog/2015/12/git-workflow.html

// Git 使用规范流程
// http://www.ruanyifeng.com/blog/2015/08/git-use-process.html

// 常用 Git 命令清单
// http://www.ruanyifeng.com/blog/2015/12/git-cheat-sheet.html

// Git远程操作详解
// http://www.ruanyifeng.com/blog/2014/06/git_remote.html

// Git分支管理策略
// http://www.ruanyifeng.com/blog/2012/07/git.html

func main() {

}
