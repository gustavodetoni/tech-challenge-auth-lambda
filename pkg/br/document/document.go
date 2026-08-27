package document

import "unicode"

func Normalize(value string) string {
	out := make([]rune, 0, len(value))
	for _, r := range value {
		if unicode.IsDigit(r) {
			out = append(out, r)
		}
	}
	return string(out)
}

func IsValidCPF(value string) bool {
	digits := Normalize(value)
	if len(digits) != 11 || allEqual(digits) {
		return false
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(digits[i]-'0') * (10 - i)
	}
	first := checkDigit(sum)
	if first != int(digits[9]-'0') {
		return false
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(digits[i]-'0') * (11 - i)
	}
	second := checkDigit(sum)
	return second == int(digits[10]-'0')
}

func IsValidCNPJ(value string) bool {
	digits := Normalize(value)
	if len(digits) != 14 || allEqual(digits) {
		return false
	}

	firstWeights := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	secondWeights := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	first := weightedCheckDigit(digits[:12], firstWeights)
	if first != int(digits[12]-'0') {
		return false
	}

	second := weightedCheckDigit(digits[:13], secondWeights)
	return second == int(digits[13]-'0')
}

func allEqual(value string) bool {
	if value == "" {
		return true
	}
	first := value[0]
	for i := 1; i < len(value); i++ {
		if value[i] != first {
			return false
		}
	}
	return true
}

func checkDigit(sum int) int {
	rest := (sum * 10) % 11
	if rest == 10 {
		return 0
	}
	return rest
}

func weightedCheckDigit(digits string, weights []int) int {
	sum := 0
	for i := range digits {
		sum += int(digits[i]-'0') * weights[i]
	}
	rest := sum % 11
	if rest < 2 {
		return 0
	}
	return 11 - rest
}
