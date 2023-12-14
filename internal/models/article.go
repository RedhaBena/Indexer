package models

type Article struct {
	Id    string `json:"_id"`
	Title string `json:"title"`

	Authors    []Author `json:"authors"`
	References []string `json:"references"`
}

type Reference struct {
	From string
	To   string
}
