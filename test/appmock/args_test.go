package appmock

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockObj struct {
	mock.Mock
}

func (m *mockObj) DoSomething(args ...any) {
	m.Called(args...)
}

func TestReturnArgs(t *testing.T) {
	type testcase struct {
		title         string
		shouldPanic   bool
		args          []any
		expectedValue []any
		expectedPanic string
	}

	testcases := []testcase{
		{
			title:         "should not panic when args are correct",
			shouldPanic:   false,
			args:          []any{new(int), new(string)},
			expectedValue: []any{999, "Test"},
		},
		{
			title:         "should panic when length of args is not equal to length of expectedValue",
			shouldPanic:   true,
			args:          []any{new(int), new(string)},
			expectedValue: []any{999},
			expectedPanic: "expected 1 arguments, got 2",
		},
		{
			title:         "should panic when args is nil",
			shouldPanic:   true,
			args:          []any{new(int), nil},
			expectedValue: []any{999, "Test"},
			expectedPanic: "argument 1 is nil",
		},
		{
			title:         "should panic when args is not a pointer",
			shouldPanic:   true,
			args:          []any{new(int), "Test"},
			expectedValue: []any{999, 0},
			expectedPanic: "argument 1 is not a pointer",
		},
		{
			title:         "should panic when type mismatch",
			shouldPanic:   true,
			args:          []any{new(int), new(string)},
			expectedValue: []any{999, 0},
			expectedPanic: "type mismatch at argument 1: cannot assign int to string",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			m := &mockObj{}

			m.On("DoSomething", tc.args...).Run(ReturnArgs(tc.expectedValue...))

			// Act
			fn := func() {
				m.DoSomething(tc.args...)
			}

			// Assert
			if tc.shouldPanic {
				assert.PanicsWithValue(t, tc.expectedPanic, fn)
			} else {
				assert.NotPanics(t, fn)
			}

			m.AssertExpectations(t)
			if tc.shouldPanic {
				return
			}

			for i, arg := range tc.args {
				val := reflect.ValueOf(arg).Elem().Interface()
				assert.Equal(t, tc.expectedValue[i], val)
			}
		})
	}
}
