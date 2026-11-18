package logic

import (
	"errors"
	"strings"
)

//1.Find tailling zero
func TrailingZeroes(number int) (int, error) {
	if number < 0 {
		return 0, errors.New("number can not be negative")
	}
	count := 0
	for number > 0 {
		number /= 5
		count += number
	}
	return count, nil
}

//2.Index of max
func MaxIndex(numbers []int) (int, error) {
	if len(numbers) == 0 {
		return 0, errors.New("list can not blank")
	}
	index := 0
	for i, n := range numbers {
		if n > numbers[index] {
			index = i
		}
	}
	return index, nil
}

//3.Number To Roman
func NumberToRoman(number int) (string, error) {
	if number < 0 {
		return "", errors.New("number can not less than 0")
	}
	if number < 1 || number > 3999 {
		return "", errors.New("roman number must be between 1 and 3999")
	}
	symbols := []struct {
		value int
		text  string
	}{{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"}, {100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"}, {10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"}}
	var result strings.Builder
	for _, symbol := range symbols {
		for number >= symbol.value {
			result.WriteString(symbol.text)
			number -= symbol.value
		}
	}
	return result.String(), nil
}

//3.Number To Thai
func NumberToThai(number int) (string, error) {
	if number < 0 {
		return "", errors.New("number can not less than 0")
	}
	if number > 10000000 {
		return "", errors.New("number must not exceed 10000000")
	}
	if number == 0 {
		return "ศูนย์", nil
	}
	return thai(number), nil
}

func thai(number int) string {
	if number >= 1000000 {
		remainder := number % 1000000
		prefix := thai(number/1000000) + "ล้าน"
		if remainder == 0 {
			return prefix
		}
		if remainder == 1 {
			return prefix + "เอ็ด"
		}
		return prefix + thai(remainder)
	}
	digits := []string{"", "หนึ่ง", "สอง", "สาม", "สี่", "ห้า", "หก", "เจ็ด", "แปด", "เก้า"}
	units := []string{"", "สิบ", "ร้อย", "พัน", "หมื่น", "แสน"}
	var result strings.Builder
	divisor := 100000
	for position := 5; position >= 0; position-- {
		digit := number / divisor % 10
		if digit != 0 {
			switch {
			case position == 1 && digit == 1:
				result.WriteString("สิบ")
			case position == 1 && digit == 2:
				result.WriteString("ยี่สิบ")
			case position == 0 && digit == 1 && number > 1:
				result.WriteString("เอ็ด")
			default:
				result.WriteString(digits[digit])
				result.WriteString(units[position])
			}
		}
		divisor /= 10
	}
	return result.String()
}
