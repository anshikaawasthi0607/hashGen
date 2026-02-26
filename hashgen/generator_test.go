package hashgen

import (
	"strings"
	"testing"
)

func TestGenerate_Deterministic(t *testing.T) {
	for _, input := range []string{"hello", "Hello", "abc123", "GoLang2024"} {
		h1, err := Generate(input)
		if err != nil {
			t.Fatalf("Generate(%q) error: %v", input, err)
		}
		h2, _ := Generate(input)
		if h1 != h2 {
			t.Errorf("Generate(%q) is not deterministic: %q != %q", input, h1, h2)
		}
	}
}

func TestGenerate_Length(t *testing.T) {
	inputs := []string{"a", "z", "A", "Z", "0", "9", "HelloWorld123", strings.Repeat("x", 512)}
	for _, input := range inputs {
		h, err := Generate(input)
		if err != nil {
			t.Fatalf("Generate(%q) error: %v", input, err)
		}
		if len(h) != 10 {
			t.Errorf("Generate(%q) length = %d, want 10", input, len(h))
		}
	}
}

func TestGenerate_AvalancheEffect(t *testing.T) {
	h1, _ := Generate("hello")
	h2, _ := Generate("hellp") // one character differs
	if h1 == h2 {
		t.Error("Avalanche effect failure: similar inputs produced identical hashes")
	}
	// Expect most characters to differ
	diff := 0
	for i := range h1 {
		if h1[i] != h2[i] {
			diff++
		}
	}
	if diff < 5 {
		t.Errorf("Weak avalanche effect: only %d/10 characters differ", diff)
	}
}
