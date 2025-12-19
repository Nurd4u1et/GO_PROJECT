package model

type Doctor struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
}
