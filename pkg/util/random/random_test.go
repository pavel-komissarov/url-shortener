package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomString_stringLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int
	}{
		{
			name: "size is 1",
			size: 1,
		},
		{
			name: "size is 3",
			size: 3,
		},
		{
			name: "size is 11",
			size: 11,
		},
		{
			name: "size is 16",
			size: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			str, _ := NewRandomString(tt.size)

			assert.Len(t, str, tt.size)
		})
	}
}

func TestNewRandomString_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		errorString string
		size        int
	}{
		{
			name:        "len < 0",
			errorString: "stringLength must be > 0",
			size:        -2,
		},
		{
			name:        "len = 0",
			errorString: "stringLength must be > 0",
			size:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			str, err := NewRandomString(tt.size)

			assert.Equal(t, tt.errorString, err.Error())
			assert.Equal(t, "", str)
		})
	}
}
