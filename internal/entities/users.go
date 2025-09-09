package entities

type User struct {
	ID       string `validate:"uuid4"`
	Nickname string `validate:"required"`
	Email    string `validate:"email"`
	//
	IsAdmin bool
}

func (user *User) GetID() string {
	return user.ID
}

func (user *User) GetEmail() string {
	return user.Email
}

type UserCreds struct {
	User
	//
	Password string `validate:"required"`
}
