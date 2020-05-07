package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqSlice(t *testing.T) {
	var testData = []struct {
		input  []interface{}
		expect []interface{}
	}{
		{[]interface{}{1, 2, 3, 4, 5, 4, 3, 2, 1}, []interface{}{1, 2, 3, 4, 5}},
		{[]interface{}{"a", "b", "c", "a"}, []interface{}{"a", "b", "c"}},
		{[]interface{}{1.0, 2.01, 3.0, 1.00, 2.00, 1.0}, []interface{}{1.0, 2.01, 3.0, 2.0}},
	}

	for _, ts := range testData {
		result, err := UniqSliceWithError(ts.input)
		if assert.NoError(t, err) {
			assert.Equal(t, result, ts.expect, "UniqSliceWithError")
		}
	}
}

func TestInSlice(t *testing.T) {
	var testData = []struct {
		input  []interface{}
		value  interface{}
		expect bool
	}{
		{[]interface{}{1, 2, 3, 4, 5, 4, 3, 2, 1}, 6, false},
		{[]interface{}{1, 2, 3, 4, 5, 4, 3, 2, 1}, 5, true},
		{[]interface{}{"a", "b", "c", "a"}, "a", true},
		{[]interface{}{"a", "b", "c", "a"}, "d", false},
		{[]interface{}{1.0, 2.01, 3.0, 1.00, 2.00, 1.0}, 2.0, true},
		{[]interface{}{1.0, 2.01, 3.0, 1.00, 2.00, 1.0}, 2.2, false},
		{[]interface{}{1.0, 2.01, 3.0, 1.00, 2.00, 1.0}, 2, false},
	}

	for _, ts := range testData {
		result, err := InSliceWithError(ts.input, ts.value)
		if assert.NoError(t, err) {
			assert.True(t, result == ts.expect, "InSliceWithError")
		}
	}
}
