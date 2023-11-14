package users

type ID string

type User struct {
	ID       ID     `json:"id"`
	Password string `json:"password"`
}
