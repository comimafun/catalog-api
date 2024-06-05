package domain

type Error struct {
	Code int     `json:"code"`
	Err  error   `json:"error"`
	Info *string `json:"info"`
}

func NewError(code int, err error, info *string) *Error {
	return &Error{
		Code: code,
		Err:  err,
		Info: info,
	}
}
