package main

import (
	"testing"
)

func TestEmailValidationHelper(t *testing.T) {

	t.Parallel() // marks test suite as capable of running in parallel with other tests

	var tests = []struct {
		email  string
		result bool
	}{
		{"aaa", false},
		{"aaa@", false},
		{"aaa@ccc", false},
		{"aaa@ccc.", false},
		{"aaa@ccc.com", true},
		{"@ccc", false},
		{"@ccc.", false},
		{"@ccc.com", false},
		{"ccc.com", false},
	}
	for _, testCase := range tests {
		testCase := testCase // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(testCase.email, func(t *testing.T) {

			t.Parallel() // marks each test case as capable of running in parallel with each other

			actualResult := isValidEmail(testCase.email)
			if actualResult != testCase.result {
				t.Errorf("Email validation for %s should return %t but returned %t instead", testCase.email, testCase.result, actualResult)
			}
		})
	}
}
