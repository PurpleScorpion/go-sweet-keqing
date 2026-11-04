package keqing

import (
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
)

var decimalAmount = regexp.MustCompile(`^[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?$`)

// Yuan2Fen 将十进制元转换为分，超出分的余数向零截断；非法输入和溢出会 panic。
func Yuan2Fen(yuan string) int64 {
	if !decimalAmount.MatchString(yuan) {
		panic("invalid decimal amount: " + yuan)
	}
	// 此检查只用于拒绝超大指数及非有限值，不参与金额计算。
	f, err := strconv.ParseFloat(yuan, 64)
	if err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
		panic("amount out of range: " + yuan)
	}
	if f == 0 {
		return 0
	} // 极小金额向零截断，避免构造巨大的十进制分母。
	amount, ok := new(big.Rat).SetString(yuan)
	if !ok {
		panic("invalid decimal amount: " + yuan)
	}
	return moneyInt64(amount.Mul(amount, big.NewRat(100, 1)))
}

// Fen2Yuan 精确格式化整数分，包括 int64 的完整范围。
func Fen2Yuan(cents int64) string {
	whole, fraction := cents/100, cents%100
	sign := ""
	if cents < 0 {
		sign, whole, fraction = "-", -whole, -fraction
	}
	return fmt.Sprintf("%s%d.%02d", sign, whole, fraction)
}

func moneyInt64(value *big.Rat) int64 {
	n := new(big.Int).Quo(value.Num(), value.Denom())
	if !n.IsInt64() {
		panic("amount exceeds int64 cents")
	}
	return n.Int64()
}

func moneyFactor(value float64) *big.Rat {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		panic("money factor must be finite")
	}
	// 最短十进制表示避免 0.1 等倍率的二进制尾差。
	r, ok := new(big.Rat).SetString(strconv.FormatFloat(value, 'g', -1, 64))
	if !ok {
		panic("invalid money factor")
	}
	return r
}

func YuanAdd(arr ...string) string {
	if len(arr) == 0 {
		return "0"
	}
	sum := new(big.Rat)
	for _, value := range arr {
		sum.Add(sum, new(big.Rat).SetInt64(Yuan2Fen(value)))
	}
	return Fen2Yuan(moneyInt64(sum))
}

func YuanSub(a, b string) string {
	value := new(big.Rat).Sub(new(big.Rat).SetInt64(Yuan2Fen(a)), new(big.Rat).SetInt64(Yuan2Fen(b)))
	return Fen2Yuan(moneyInt64(value))
}

func YuanMul(a string, b float64) string {
	value := new(big.Rat).Mul(new(big.Rat).SetInt64(Yuan2Fen(a)), moneyFactor(b))
	return Fen2Yuan(moneyInt64(value))
}

func YuanDiv(a string, b float64) string {
	if b == 0 {
		panic("The divisor cannot be 0")
	}
	value := new(big.Rat).Quo(new(big.Rat).SetInt64(Yuan2Fen(a)), moneyFactor(b))
	return Fen2Yuan(moneyInt64(value))
}
