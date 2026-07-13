package utils

import (
	"strings"
	"testing"
)


func TestGenerateShortCode(t *testing.T) {
	
	code,err := GenerateShortCode()
	if err != nil {
		t.Fatalf("generateShortCode() = %v; expected %v", err, nil)
	}
	if len(code) != shortCodeLength {
		t.Fatalf("generated code length = %d, want %d",
	len(code),
	shortCodeLength,
)
	}
	for _,char := range code {
		if !strings.Contains(charset,string(char)) {
			t.Errorf("generated invalid character %q in %q", char, code)
		}
	}
}