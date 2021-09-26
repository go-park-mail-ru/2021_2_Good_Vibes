package user

type UserInput struct {
	Name     string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
