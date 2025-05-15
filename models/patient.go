package models

type Patient struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Diagnosis string `json:"diagnosis"`
}
