package logic

import (
	"testing"
)

func TestLogic(t *testing.T) {
	for _, tc := range []struct{ n, w int }{{0, 0}, {7, 1}, {25, 6}, {100, 24}} {
		got, err := TrailingZeroes(tc.n)
		if err != nil || got != tc.w {
			t.Fatalf("zeroes %d: %d %v", tc.n, got, err)
		}
	}
	if _, err := TrailingZeroes(-10); err == nil {
		t.Fatal("negative accepted")
	}
	for _, tc := range []struct {
		numbers []int
		want    int
	}{{[]int{1, 2, 1, 3, 5, 6, 4}, 5}, {[]int{-3, -1, -1}, 1}} {
		got, err := MaxIndex(tc.numbers)
		if err != nil || got != tc.want {
			t.Fatal(got, err)
		}
	}
	if _, err := MaxIndex(nil); err == nil {
		t.Fatal("empty accepted")
	}
	for n, want := range map[int]string{0: "ศูนย์", 1: "หนึ่ง", 11: "สิบเอ็ด", 20: "ยี่สิบ", 21: "ยี่สิบเอ็ด", 101: "หนึ่งร้อยเอ็ด", 1000001: "หนึ่งล้านเอ็ด", 10000000: "สิบล้าน"} {
		got, err := NumberToThai(n)
		if err != nil || got != want {
			t.Fatalf("thai %d: %q %v", n, got, err)
		}
	}
	for _, n := range []int{-1, 10000001} {
		if _, err := NumberToThai(n); err == nil {
			t.Fatal(n)
		}
	}
	for n, want := range map[int]string{1: "I", 4: "IV", 9: "IX", 101: "CI", 1994: "MCMXCIV", 3999: "MMMCMXCIX"} {
		got, err := NumberToRoman(n)
		if err != nil || got != want {
			t.Fatal(got, err)
		}
	}
	for _, n := range []int{-1, 0, 4000} {
		if _, err := NumberToRoman(n); err == nil {
			t.Fatal(n)
		}
	}
}
