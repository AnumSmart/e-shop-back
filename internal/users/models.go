package users

type User struct {
	Id           int
	Email        string
	HashPass     string
	RefreshToken string
}
