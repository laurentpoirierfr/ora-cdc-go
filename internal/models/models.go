package models

type Info struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Health struct {
	Status string `json:"status"`
}
