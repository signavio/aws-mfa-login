package action

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAwsVersion(t *testing.T) {

	var tests = map[string]struct {
		input       string
		expected    string
		errExpected error
	}{
		"detect-windows":    {"aws-cli/1.17.0 Python/3.6.0 Windows/10 botocore/1.14.0", "1.17.0", nil},
		"detect-linux":      {"aws-cli/1.17.1 Python/3.8.1 Linux/4.9.184-linuxkit botocore/1.14.1", "1.17.1", nil},
		"command-not-found": {"bash: aws: command not found", "", &AwsVersionParseException{}},
		"wrong-version":     {"aws-cli/1.17 Python/3.8.1 Linux/4.9.184-linuxkit botocore/1.14.1", "", &AwsVersionParseException{}},
	}
	for _, test := range tests {
		version, err := parseAwsVersion(test.input)
		assert.IsType(t, test.errExpected, err)
		assert.Equal(t, test.expected, version)
	}
}

func TestCheckRequiredAwsVersion(t *testing.T) {
	var tests = map[string]struct {
		input    string
		expected bool
		hasError bool
	}{
		"success":           {"1.17.0", true, false},
		"smaller-version":   {"1.16.0", false, true},
		"no-semver-version": {"1.2", false, true},
	}
	for _, test := range tests {
		fmt.Printf("compare %s against %s\n", test.input, RequiredMinAwsVersion)
		isValid, err := CheckRequiredAwsVersion(test.input)
		if err != nil {
			fmt.Println(err)
		}
		assert.Equal(t, test.expected, isValid)
		assert.Equal(t, test.hasError, err != nil)
	}
}
