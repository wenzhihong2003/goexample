package models_test

// import "github.com/gopherguides/training/fundamentals/testing/src/cover"
import (
	"testing"

	"github.com/gopherguides/training/fundamentals/testing/src/cover"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		user  models.User
		valid bool
	}{
		{user: models.User{}, valid: false},
		{user: models.User{First: "Rob"}, valid: false},
		{user: models.User{Last: "Pike"}, valid: false},
		{user: models.User{First: "Rob", Last: "Pike"}, valid: true},
	}
	for _, test := range tests {
		err := test.user.Validate()
		if err != nil && test.valid {
			t.Errorf("unexpected error: %s", err)
		}
	}
}
