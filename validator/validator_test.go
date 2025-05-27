package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	Name string `json:"name" validate:"required"`
}

func TestValidate(t *testing.T) {
	type testcase struct {
		title   string
		data    data
		wantErr bool
	}

	testcases := []testcase{
		{
			title:   "should return no error when data is valid",
			data:    data{Name: "Valid Name"},
			wantErr: false,
		},
		{
			title:   "should return error when name is empty",
			data:    data{Name: ""},
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			validator := NewReqValidator()
			err := validator.Validate(tc.data)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
