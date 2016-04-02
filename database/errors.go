package database

const DatabaseConfigurationError = 701

type DatabaseError struct {
	err       string
	ErrorCode int
}

func (dbe DatabaseError) Error() string {
	if dbe.err != "" {
		return dbe.err
	}
	return "The database encountered an error."
}
