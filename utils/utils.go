package utils

func RemovePercentSign(raw string) string {
	res := make([]rune, 0)
	for _, v := range raw {
		if v != '%' {
			res = append(res, v)
		}
	}
	return string(res)
}
