package models

type Author struct {
	Id   *string `json:"_id,omitempty"`
	Name string  `json:"name"`
}
