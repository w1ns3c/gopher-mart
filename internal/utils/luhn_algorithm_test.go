package utils

import "testing"

func TestLuhnValidator(t *testing.T) {
	tests := []struct {
		name     string
		sequence string
		want     bool
	}{
		{
			name:     "Bank card 1 (even sequence)",
			sequence: "2200700147804412",
			want:     true,
		},
		{
			name:     "Bank card 2 (odd sequence)",
			sequence: "12200700147804411",
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LuhnValidator(tt.sequence); got != tt.want {
				t.Errorf("LuhnValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}
