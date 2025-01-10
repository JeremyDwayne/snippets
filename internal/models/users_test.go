package models

import (
	"testing"

	"github.com/jeremydwayne/snippets/internal/assert"
	"github.com/jeremydwayne/snippets/internal/sqlc"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID string
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: "e30fd85a-efd2-44d0-86ed-88e71a8dfeda",
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: "",
			want:   false,
		},
		{
			name:   "Non-Existent ID",
			userID: "b30fd85a-efd2-44d0-86ed-88e71a8dfeda",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := sqlc.New(newTestDB(t))

			m := UserModel{db}

			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
