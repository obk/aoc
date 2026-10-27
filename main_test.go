package main

import (
	"slices"
	"testing"
)

func TestSplitRunArgs(t *testing.T) {
	tests := []struct {
		args    []string
		nums    []string
		input   string
		example bool
		fails   bool
	}{
		{args: nil},
		{args: []string{"2024", "3", "-e"}, nums: []string{"2024", "3"}, example: true},
		{args: []string{"-e", "3"}, nums: []string{"3"}, example: true},
		{args: []string{"3", "-i", "big.txt"}, nums: []string{"3"}, input: "big.txt"},
		{args: []string{"3", "-i"}, fails: true},
		{args: []string{"-e", "-i", "big.txt"}, fails: true},
		{args: []string{"-x"}, fails: true},
	}
	for _, tt := range tests {
		nums, input, example, err := splitRunArgs(tt.args)
		if tt.fails {
			if err == nil {
				t.Errorf("splitRunArgs(%q) should fail", tt.args)
			}
			continue
		}
		if err != nil || !slices.Equal(nums, tt.nums) || input != tt.input || example != tt.example {
			t.Errorf("splitRunArgs(%q) = %q, %q, %v, %v", tt.args, nums, input, example, err)
		}
	}
}
