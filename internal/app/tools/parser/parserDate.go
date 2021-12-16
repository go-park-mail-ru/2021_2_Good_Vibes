package parser

//2021-12-12T22:21:28+03:00
func ParseDateFromSql(date string) string {
	var answer string
	counterDoubleTwoPoints := 0
	for _, value := range date {
		if value == '+' {
			break
		}
		if value == '-' {
			answer += "."
			continue
		}
		if value == ':' {
			counterDoubleTwoPoints++
			if counterDoubleTwoPoints == 2 {
				break
			}
		}
		if value != 'T' {
			answer += string(value)
		} else {
			answer += " "
		}
	}
	return answer
}
