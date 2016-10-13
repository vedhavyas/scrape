package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterSlice(t *testing.T) {
	cases := []struct {
		BaseSlice     []string
		NewSlice      []string
		FilteredSlice []string
	}{
		{
			BaseSlice:     []string{"a", "b", "c", "d"},
			NewSlice:      []string{"c", "d", "e"},
			FilteredSlice: []string{"e"},
		},

		{
			BaseSlice:     []string{"a", "b", "c", "d", "e"},
			NewSlice:      []string{"c", "d", "e", "f", "g"},
			FilteredSlice: []string{"f", "g"},
		},

		{
			BaseSlice:     []string{"a", "b", "c", "d", "a"},
			NewSlice:      []string{"c", "d", "e", "e"},
			FilteredSlice: []string{"e", "e"},
		},

		{
			BaseSlice:     []string{"a", "b", "c", "d", "a"},
			NewSlice:      []string{"c", "d", "e", "e", "d"},
			FilteredSlice: []string{"e", "e"},
		},
	}

	for _, testCase := range cases {
		actualSlice := FilterSlice(testCase.BaseSlice, testCase.NewSlice)
		assert.EqualValues(t, testCase.FilteredSlice, actualSlice)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	cases := []struct {
		Slice  []string
		Result []string
	}{
		{
			Slice:  []string{"a", "b", "c", "d"},
			Result: []string{"a", "b", "c", "d"},
		},

		{
			Slice:  []string{"a", "b", "c", "d", "e", "a", "c", "a"},
			Result: []string{"a", "b", "c", "d", "e"},
		},

		{
			Slice:  []string{"a", "b", "b", "c", "d", "d", "e", "a", "c", "e"},
			Result: []string{"a", "b", "c", "d", "e"},
		},

		{
			Slice:  []string{"a"},
			Result: []string{"a"},
		},
	}

	for _, testCase := range cases {
		actualSlice := RemoveDuplicates(testCase.Slice)
		assert.EqualValues(t, len(testCase.Result), len(actualSlice))
	}
}
