package service

import (
	"testing"
	"time"
)

func TestCodeGenerator_Generate(t *testing.T) {
	epoch, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	generator, err := NewCodeGenerator(epoch, 1, 6)
	if err != nil {
		t.Fatalf("Failed to create code generator: %v", err)
	}

	// 测试生成短码
	codes := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		code := generator.Generate()

		// 验证长度
		if len(code) != 6 {
			t.Errorf("Expected code length 6, got %d", len(code))
		}

		// 验证唯一性
		if codes[code] {
			t.Errorf("Duplicate code generated: %s", code)
		}
		codes[code] = true

		// 验汪字符集
		for _, c := range code {
			if !isValidBase62Char(c) {
				t.Errorf("Invalid character in code: %c", c)
			}
		}
	}
}

func TestToBase62(t *testing.T) {
	tests := []struct {
		num    int64
		length int
		want   string
	}{
		{0, 6, "000000"},
		{1, 6, "000001"},
		{61, 6, "00000z"},
		{62, 6, "000010"},
		{1000000, 6, "00w7eO"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := ToBase62(tt.num, tt.length)
			if got != tt.want {
				t.Errorf("ToBase62(%d, %d) = %v, want %v", tt.num, tt.length, got, tt.want)
			}
		})
	}
}

func TestFromBase62(t *testing.T) {
	tests := []struct {
		code string
		want int64
	}{
		{"000000", 0},
		{"000001", 1},
		{"00000z", 61},
		{"000010", 62},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := FromBase62(tt.code)
			if got != tt.want {
				t.Errorf("FromBase62(%v) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestBase62RoundTrip(t *testing.T) {
	// 测试编码和解码的互逆性
	originalValues := []int64{0, 1, 61, 62, 100, 1000, 10000, 100000, 1000000}

	for _, val := range originalValues {
		encoded := ToBase62(val, 6)
		decoded := FromBase62(encoded)

		if decoded != val {
			t.Errorf("Round trip failed for %d: encoded to %s, decoded to %d", val, encoded, decoded)
		}
	}
}

func isValidBase62Char(c rune) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z')
}
