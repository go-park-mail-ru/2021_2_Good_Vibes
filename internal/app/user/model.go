package user

type UserInput struct {
	Name     string `json:"username" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type UserStorage struct {
	Id       int    `json:"id"`
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,customPassword"`
}

type Error struct {
	ErrorCode        int    `json:"error code"`
	ErrorDescription string `json:"error description"`
}

func NewError(errorCode int, errorDesc string) *Error {
	return &Error{
		ErrorCode:        errorCode,
		ErrorDescription: errorDesc,
	}
}

func NewUserStorage(id int, name string, email string, password string) UserStorage {
	return UserStorage{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
}
