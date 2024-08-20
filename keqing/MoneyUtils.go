package keqing

import "strconv"

/*
将人民币元转成分
*/
func Yuan2Fen(yuan string) int64 {
	// Convert cleaned string to float
	f, err := strconv.ParseFloat(yuan, 64)
	if err != nil {
		panic("Failed to parse float")
	}
	// Multiply by 100 and convert to integer
	return int64(f * 100)
}

/*
将人民币分转成元
*/
func Fen2Yuan(cents int64) string {
	// Divide by 100 to get the amount in dollars as a float64
	dollars := float64(cents) / 100

	// Format the float64 to have two decimal places
	formattedDollars := strconv.FormatFloat(dollars, 'f', 2, 64)

	return formattedDollars
}

/*
将多个元相加
*/
func YuanAdd(arr ...string) string {
	if len(arr) == 0 {
		return "0"
	}
	var sum int64 = 0
	for _, v := range arr {
		sum += Yuan2Fen(v)
	}
	return Fen2Yuan(sum)
}

/*
将2个元相减
*/
func YuanSub(a string, b string) string {
	return Fen2Yuan(Yuan2Fen(a) - Yuan2Fen(b))
}

/*
将元乘以一个浮点数
*/
func YuanMul(a string, b float64) string {

	var sum = float64(Yuan2Fen(a)) * b
	return Fen2Yuan(int64(sum))
}

/*
将元除以一个浮点数
*/
func YuanDiv(a string, b float64) string {

	var sum = float64(Yuan2Fen(a)) / b
	return Fen2Yuan(int64(sum))
}
