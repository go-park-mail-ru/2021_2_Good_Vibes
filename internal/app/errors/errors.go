package errors

const (
	BIND_ERROR           = -20
	VALIDATION_ERROR     = -21
	TOKEN_ERROR          = 0
	NO_USER_ERROR        = -30
	USER_EXISTS_ERROR    = -32
	WRONG_PASSWORD_ERROR = -33
	DB_ERROR             = -40
	SERVER_ERROR         = -50
)

const (
	BIND_DESCR           = "can not bind data"
	VALIDATION_DESCR     = "Неправильный формат данных"
	TOKEN_ERROR_DESCR    = "can not get token"
	NO_USER_DESCR        = "Пользователя не существует"
	USER_EXISTS_DESCR    = "Пользователь уже существует"
	WRONG_PASSWORD_DESCR = "Неверный пароль"
	BD_ERROR_DESCR       = "bd error"
	BAD_INIT_SECRET_KEY  = "bad init secret key"
	HASHER_ERROR_DESCR   = "error hash password"
)

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
