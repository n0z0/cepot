package main

import (
	"testing"
)

func TestParseAnswerOption(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"A", "A"},
		{"a", "A"},
		{"B", "B"},
		{"b", "B"},
		{"C", "C"},
		{"c", "C"},
		{"D", "D"},
		{"d", "D"},
		{"E", "E"},
		{"e", "E"},
		{"1", "A"},
		{"2", "B"},
		{"3", "C"},
		{"4", "D"},
		{"5", "E"},
		{"**A**", "A"},
		{"Jawaban: B", "B"},
		{"(C)", "C"},
		{"d.", "D"},
		{"E. Pilihan tepat", "E"},
		{"", ""},
		{"xyz", ""},
	}

	for _, tt := range tests {
		got := ParseAnswerOption(tt.input)
		if got != tt.expected {
			t.Errorf("ParseAnswerOption(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestDirectionMap(t *testing.T) {
	options := []string{"A", "B", "C", "D", "E"}
	for _, opt := range options {
		dir, ok := DirectionMap[opt]
		if !ok {
			t.Fatalf("DirectionMap missing option %s", opt)
		}
		if dir.Option != opt {
			t.Errorf("DirectionMap[%s].Option = %s; want %s", opt, dir.Option, opt)
		}
		if dir.Dx == 0 && dir.Dy == 0 {
			t.Errorf("DirectionMap[%s] has zero movement", opt)
		}
	}
}
