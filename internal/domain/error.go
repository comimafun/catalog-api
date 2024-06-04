package domain

type Error struct {
	Code int     `json:"code"`
	Err  error   `json:"error"`
	Info *string `json:"info"`
}
