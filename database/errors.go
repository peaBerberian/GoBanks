package database

const (
	UnknownDatabaseErrorCode = 700 + iota
	DatabaseConfigurationErrorCode
	UnsupportedDatabaseErrorCode
	MissingInformationsErrorCode
	DatabaseQueryErrorCode
	DatabaseConnectionErrorCode
)

type databaseError interface {
	error
	ErrorCode() uint32
}

type genericDatabaseError struct {
	err  string
	code uint32
}
type databaseConfigurationError struct{}
type unsupportedDatabaseError struct{ database string }
type missingInformationsError struct{ field string }
type databaseQueryError struct{ err string }

func (dbe genericDatabaseError) Error() string {
	if dbe.err != "" {
		return dbe.err
	}
	return "The database encountered an error."
}

func (dbe genericDatabaseError) ErrorCode() uint32 {
	if dbe.code == 0 {
		return UnknownDatabaseErrorCode
	}
	return dbe.code
}

func (d databaseConfigurationError) Error() string {
	return "The configuration for this database is not valid."
}

func (d databaseConfigurationError) ErrorCode() uint32 {
	return DatabaseConfigurationErrorCode
}

func (d unsupportedDatabaseError) Error() string {
	if d.database != "" {
		return "The database \"" + d.database + "\" is not supported."
	}
	return "The configured database is not supported."
}

func (e unsupportedDatabaseError) ErrorCode() uint32 {
	return MissingInformationsErrorCode
}

func (e missingInformationsError) Error() string {
	if e.field != "" {
		return "The field \"" + e.field + "\" needs to be filled."
	}
	return "The given request is not complete."
}

func (d missingInformationsError) ErrorCode() uint32 {
	return UnsupportedDatabaseErrorCode
}

func (e databaseQueryError) Error() string {
	if e.err != "" {
		return e.err
	}
	return "The database query was malformed."
}

func (e databaseQueryError) ErrorCode() uint32 {
	return DatabaseQueryErrorCode
}
