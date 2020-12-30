package math

func Mod(a, b int) int {
	if a < 0 {
		return b - 1 - ^a%b
	}
	return a % b
}

func Div(a, b int) int {
	if a < 0 {
		return ^(^a / b)
	}
	return a / b
}
