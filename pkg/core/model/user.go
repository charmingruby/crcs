package model

type User struct {
	FirstName    string `json:"first_name" db:"first_name"`
	LastName     string `json:"last_name"  db:"last_name"`
	Email        string `json:"email"      db:"email"`
	PasswordHash string `json:"-"          db:"password_hash"`
	Model        `                         db:",inline"`
}

type UserInput struct {
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
}

func NewUser(input UserInput) *User {
	return &User{
		Model:        NewModel(),
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
	}
}
