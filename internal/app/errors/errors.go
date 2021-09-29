package errors

const (
	BIND_ERROR = 20
	VALIDATION_ERROR = 21
	TOKEN_ERROR = 22
	NO_USER_ERROR = 30
	USER_EXISTS_ERROR = 32
	WRONG_PASSWORD_ERROR = 33
	DB_ERROR = 40
	SERVER_ERROR = 50
)

const (
	BIND_DESCR = "can not bind data"
	VALIDATION_DESCR = "can not validate data"
	TOKEN_ERROR_DESCR = "can not get token"
	NO_USER_DESCR = "user does not exist"
	USER_EXISTS_DESCR = "user already exists"
	WRONG_PASSWORD_DESCR = "wrong password"
)