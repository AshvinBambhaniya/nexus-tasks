package utils

import (
	"errors"
	"testing"

	"gopkg.in/go-playground/validator.v9"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		want    bool
		wantErr bool
	}{
		{
			name:    "valid simple email",
			email:   "user@example.com",
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid email with dot",
			email:   "first.last@domain.org",
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid email with sub-domain",
			email:   "user@mail.example.co.uk",
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid - no @",
			email:   "userexample.com",
			want:    false,
			wantErr: false,
		},
		{
			name:    "invalid - no domain",
			email:   "user@",
			want:    false,
			wantErr: false,
		},
		{
			name:    "invalid - no TLD",
			email:   "user@domain",
			want:    false,
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateEmail() got = %v, want %v for email %q", got, tt.want, tt.email)
			}
		})
	}
}

func TestValidatorErrorString(t *testing.T) {
	v := validator.New()

	type User struct {
		Email string `validate:"required,email"`
		Name  string `validate:"required"`
	}

	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "nil error",
			err:  nil,
			want: "",
		},
		{
			name: "multiple field errors",
			err:  v.Struct(User{}),
			want: "email,name fields are invalid.",
		},
		{
			name: "single field error",
			err:  v.Struct(User{Email: "valid@mail.com"}), // Name is missing
			want: "name fields are invalid.",
		},
		{
			name: "non-validator error",
			err:  errors.New("some other error"),
			want: "panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want == "panic" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("ValidatorErrorString() did not panic as expected for non-validator error")
					}
				}()
				ValidatorErrorString(tt.err)
				return
			}

			got := ValidatorErrorString(tt.err)
			if got != tt.want {
				t.Errorf("ValidatorErrorString() got = %q, want %q", got, tt.want)
			}
		})
	}
}
