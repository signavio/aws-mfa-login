package action

type ArnParseException struct {
	msg string
}

func (err *ArnParseException) Error() string {
	return err.msg
}

type AwsVersionParseException struct {
	msg string
}

func (err *AwsVersionParseException) Error() string {
	return err.msg
}
