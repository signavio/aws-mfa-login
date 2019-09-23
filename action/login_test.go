package action

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseUsername(t *testing.T) {
	var tests = map[string]struct {
		caller           *sts.GetCallerIdentityOutput
		expectedUsername string
		errExpected      error
	}{
		"successful-1": {caller: &sts.GetCallerIdentityOutput{Arn: aws.String("someacc/someuser")}, expectedUsername: "someuser", errExpected: nil},
		"successful-2": {caller: &sts.GetCallerIdentityOutput{Arn: aws.String("arn:aws:iam::123456:user/someuser2")}, expectedUsername: "someuser2", errExpected: nil},
		"failed-1":     {caller: &sts.GetCallerIdentityOutput{Arn: aws.String("someacc/")}, expectedUsername: "", errExpected: &ArnParseException{}},
		"failed-2":     {caller: &sts.GetCallerIdentityOutput{Arn: aws.String("someacc")}, expectedUsername: "", errExpected: &ArnParseException{}},
	}
	for _, test := range tests {
		name, err := ParseUsername(test.caller)
		assert.IsType(t, test.errExpected, err)
		assert.Equal(t, test.expectedUsername, name)
	}
}
