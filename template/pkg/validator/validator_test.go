package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"-" validate:"required"`
}

func TestValidate(t *testing.T) {
	type testcase struct {
		title         string
		data          data
		expected      error
		validateError []string
	}

	testcases := []testcase{
		{
			title:    "should return no error when data is valid",
			data:     data{Name: "Valid Name", Age: 30},
			expected: nil,
		},
		{
			title: "should return error when name is empty",
			data:  data{Name: ""},
			expected: errorMap{
				"name": "required",
				"Age":  "required",
			},
			validateError: []string{"name: required", "Age: required"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			validator := NewReqValidator()
			err := validator.Validate(tc.data)
			assert.Equal(t, tc.expected, err)
			if err != nil {
				for _, ve := range tc.validateError {
					assert.Contains(t, err.Error(), ve)
				}
			}
		})
	}
}
