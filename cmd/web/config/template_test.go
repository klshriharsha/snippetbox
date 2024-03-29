package config

import (
	"testing"
	"time"

	"github.com/klshriharsha/snippetbox/internal/assert"
)

func TestStandardDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 27, 10, 30, 0, 0, time.UTC),
			want: "27 Mar 2024 at 10:30",
		},
		{
			name: "IST",
			tm:   time.Date(2024, 3, 27, 10, 30, 0, 0, time.FixedZone("IST", 60*60*5.5)),
			want: "27 Mar 2024 at 05:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := standardDate(tt.tm)
			assert.Equal(t, sd, tt.want)
		})
	}
}
