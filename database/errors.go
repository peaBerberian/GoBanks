package database

const DatabaseConfigurationError = 701

type DatabaseError struct {
	err  string
	code uint32
}

func (dbe DatabaseError) Error() string {
	if dbe.err != "" {
		return dbe.err
	}
	return "The database encountered an error."
}

func (dbe DatabaseError) ErrorCode() uint32 {
	return dbe.code
}
