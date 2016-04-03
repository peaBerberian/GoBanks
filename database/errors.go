package database

const DatabaseConfigurationError = 701

type databaseError struct {
	err  string
	code uint32
}

func (dbe databaseError) Error() string {
	if dbe.err != "" {
		return dbe.err
	}
	return "The database encountered an error."
}

func (dbe databaseError) ErrorCode() uint32 {
	return dbe.code
}
