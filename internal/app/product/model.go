package product

type Product struct {
	Id     int     `json:"id"`
	Image  string  `json:"image"`
	Name   string  `json:"name"`
	Price  int     `json:"price"`
	Rating float32 `json:"rating"`
}

func NewProduct(id int, image string, name string, price int, rating float32) Product {
	return Product{
		Id:     id,
		Image:  image,
		Name:   name,
		Price:  price,
		Rating: rating,
	}
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
