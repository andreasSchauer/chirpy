package main

import (
	"testing"
)

func TestCleanMessage(t *testing.T) {
	tests := []struct {
		name         	string
		inputMessage    string
		expected      	string
	}{
		{
			name:     		"bad word plus dot",
			inputMessage: 	"keRfuffle asidoajd bla bla shArbert fOrnax.",
			expected: 		"**** asidoajd bla bla **** fOrnax.",
		},
		{
			name:     		"all bad words",
			inputMessage: 	"keRfufFle asidoajd bla bla shARbert FOrnax .",
			expected: 		"**** asidoajd bla bla **** **** .",
		},
		{
			name:     		"no bad words",
			inputMessage: 	"I really enjoy turtles",
			expected: 		"I really enjoy turtles",
		},
		{
			name:     		"one word",
			inputMessage: 	"keRfuffle",
			expected: 		"****",
		},
		{
			name:     		"empty message",
			inputMessage: 	"",
			expected: 		"",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := cleanMessage(tc.inputMessage)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected Message: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}