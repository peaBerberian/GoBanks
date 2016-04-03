package api

const (
	UnknownOperationErrorCode = 700 + iota
	BodyParsingErrorCode
	QueryOperationErrorCode
	MissingParameterErrorCode
	NotPermittedOperationErrorCode
)

type OperationError interface {
	error
	ErrorCode() uint32
}

type bodyParsingError struct{}
type genericOperationError struct{}
type queryOperationError struct{}
type missingParameterError struct{ parameter string }
type notPermittedOperationError struct{}

func (e genericOperationError) Error() string {
	return "The operation failed."
}

func (e genericOperationError) ErrorCode() uint32 {
	return UnknownOperationErrorCode
}

func (e bodyParsingError) ErrorCode() uint32 {
	return BodyParsingErrorCode
}

func (e bodyParsingError) Error() string {
	return "Could not read your request body."
}

func (e queryOperationError) Error() string {
	return "The query to perform the wanted operation failed."
}

func (e queryOperationError) ErrorCode() uint32 {
	return QueryOperationErrorCode
}

func (e missingParameterError) Error() string {
	if e.parameter != "" {
		return "Invalid request. Missing parameter: " + e.parameter
	}
	return "Invalid request. Missing parameters."
}

func (e missingParameterError) ErrorCode() uint32 {
	return MissingParameterErrorCode
}

func (e notPermittedOperationError) Error() string {
	return "The wanted operation is not permitted."
}

func (e notPermittedOperationError) ErrorCode() uint32 {
	return NotPermittedOperationErrorCode
}
