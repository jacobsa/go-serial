package serial

import (
	"fmt"
	"testing"
)

func TestIsStandardBaudRate(t *testing.T) {
	testCases := []struct {
		BaudRate   uint
		IsStandard bool
	}{
		{50, true},
		{75, true},
		{110, true},
		{134, true},
		{150, true},
		{200, true},
		{300, true},
		{600, true},
		{1200, true},
		{1800, true},
		{2400, true},
		{4800, true},
		{7200, true},
		{9600, true},
		{14400, true},
		{19200, true},
		{28800, true},
		{38400, true},
		{57600, true},
		{76800, true},
		{115200, true},
		{230400, true},
		{0, false},
		{123, false},
		{14401, false},
	}

	for _, testCase := range testCases {
		testName := fmt.Sprintf("%d", testCase.BaudRate)
		t.Run(testName, func(t *testing.T) {
			result := IsStandardBaudRate(testCase.BaudRate)

			if result != testCase.IsStandard {
				t.Errorf("expected result to be %t, but got %t", testCase.IsStandard, result)
			}
		})
	}
}
