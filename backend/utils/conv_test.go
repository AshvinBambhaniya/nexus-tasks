package utils_test

import (
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/utils"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		text any
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.GetString(tt.text)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}
