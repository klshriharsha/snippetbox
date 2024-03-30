package models_test

import (
	"testing"

	"github.com/klshriharsha/snippetbox/internal/assert"
	"github.com/klshriharsha/snippetbox/internal/models"
	"github.com/klshriharsha/snippetbox/internal/testutils"
)

func TestUserExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := testutils.NewTestDB(t)
			model := models.UserModel{conn}
			exists, err := model.Exists(tt.userID)
			assert.NilError(t, err)
			assert.Equal(t, exists, tt.want)
		})
	}
}
