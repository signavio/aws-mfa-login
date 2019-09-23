package action

type ArnParseException struct {
	msg string
}

func (err *ArnParseException) Error() string {
	return err.msg
}
