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
	NO_REVIEW_ERROR      = -60
	REVIEW_EXISTS_ERROR  = -62
	BAD_QUERY_PARAM      = -5
)

const (
	BIND_DESCR            = "can not bind data"
	VALIDATION_DESCR      = "Неправильный формат данных"
	TOKEN_ERROR_DESCR     = "can not get token"
	NO_USER_DESCR         = "Пользователя не существует"
	USER_EXISTS_DESCR     = "Пользователь уже существует"
	WRONG_PASSWORD_DESCR  = "Неверный пароль"
	BD_ERROR_DESCR        = "bd error"
	BAD_INIT_SECRET_KEY   = "bad init secret key"
	HASHER_ERROR_DESCR    = "error hash password"
	NO_REVIEW_DESCR       = "отзыв не существует"
	REVIEW_EXISTS_DESCR   = "отзыв уже существует"
	BAD_QUERY_PARAM_DESCR = "bad query param"
)

var ErrorsMap = map[string]int{
	BIND_DESCR:            BIND_ERROR,
	VALIDATION_DESCR:      VALIDATION_ERROR,
	TOKEN_ERROR_DESCR:     TOKEN_ERROR,
	NO_USER_DESCR:         NO_USER_ERROR,
	USER_EXISTS_DESCR:     USER_EXISTS_ERROR,
	WRONG_PASSWORD_DESCR:  WRONG_PASSWORD_ERROR,
	BD_ERROR_DESCR:        SERVER_ERROR,
	NO_REVIEW_DESCR:       NO_REVIEW_ERROR,
	REVIEW_EXISTS_DESCR:   REVIEW_EXISTS_ERROR,
	BAD_QUERY_PARAM_DESCR: BAD_QUERY_PARAM,
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

func ErrorStringToCode(str string) int {
	code, ok := ErrorsMap[str]
	if !ok {
		return SERVER_ERROR
	}
	return code
}
