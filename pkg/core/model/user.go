package model

type User struct {
	FirstName    string `json:"first_name"    db:"first_name"`
	LastName     string `json:"last_name"     db:"last_name"`
	Email        string `json:"email"         db:"email"`
	PasswordHash string `json:"-"             db:"password_hash"`
	Model        `            db:",inline"`
}
