package http

type body struct {
	Command string `json:"command"`
	Id      string `json:"id"`
	Message string `json:"message"`
}
