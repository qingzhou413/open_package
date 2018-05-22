package exception

type PackageError struct {
	Msg string
}

func (e *PackageError) Error() string{
	return e.Msg
}