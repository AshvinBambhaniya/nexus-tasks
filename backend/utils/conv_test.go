package utils_test

import (
	"testing"

	"github.com/AshvinBambhaniya/nexus-tasks/v2/utils"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name string
		text any
		want string
	}{
		{
			name: "string input",
			text: "hello",
			want: "hello",
		},
		{
			name: "int input",
			text: 123,
			want: "123",
		},
		{
			name: "bool input",
			text: true,
			want: "true",
		},
		{
			name: "nil input",
			text: nil,
			want: "<nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.GetString(tt.text); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}
