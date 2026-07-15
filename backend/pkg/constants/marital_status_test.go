package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaritalStatusValue(t *testing.T) {
	tests := []struct {
		name    string
		status  MaritalStatus
		wantNil bool
		wantVal string
	}{
		{"empty returns nil", "", true, ""},
		{"TK returns TK", MaritalStatusSingle, false, "TK"},
		{"K returns K", MaritalStatusMarried, false, "K"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.status.Value()
			assert.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, val)
			} else {
				assert.Equal(t, tt.wantVal, val)
			}
		})
	}
}

func TestMaritalStatusScan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    MaritalStatus
		wantErr bool
	}{
		{"nil sets empty", nil, "", false},
		{"TK string", "TK", MaritalStatusSingle, false},
		{"K string", "K", MaritalStatusMarried, false},
		{"TK bytes", []byte("TK"), MaritalStatusSingle, false},
		{"int fails", 123, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s MaritalStatus
			err := s.Scan(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, s)
			}
		})
	}
}
