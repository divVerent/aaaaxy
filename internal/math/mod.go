package math

func Mod(a, b int) int {
	m := a % b
	if m < 0 {
		m += b
	}
	return m
}
