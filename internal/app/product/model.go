package product

type Product struct {
	Id int `json:"id:"`
	Image string `json:"image"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Rating float32 `json:"rating"`
}

type Error struct {
	ErrorCode int `json:"error code" validate:"required"`
	ErrorDescription string `json:"error description" validate:"required"`
}

func NewError(errorCode int, errorDesc string) *Error {
	return &Error{
		ErrorCode: errorCode,
		ErrorDescription: errorDesc,
	}
}
