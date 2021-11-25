package postgre

import (
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/labstack/echo/v4"
	"strconv"
)

const (
	maxRatingOnStore       = 5
	maxPriceProductOnStore = 1000000.0
)

const (
	NamePriceMin  = "price_min"
	NamePriceMax  = "price_max"
	NameRatingMin = "rating_min"
	NameRatingMax = "number"
	NameOrder     = "order"
	NameOrderType = "order_type"
	NameCategory  = "category"
)

const (
	TypeOrderPrice      = "price"
	TypeOrderRating     = "rating"
	TypeOrderMin        = "desc"
	TypeOrderMax        = "asc"
	TypeCategoryDefault = "ALL_THINGS"
)

type Filter struct {
	NameCategory string
	MinPrice     float64
	MaxPrice     float64
	MinRating    float64
	MaxRating    float64
	OrderBy      string
	TypeOrder    string
}

func ParseQueryFilter(ctx echo.Context) (*Filter, error) {
	var result Filter
	queryParam := ctx.QueryParams()

	minPriceString := queryParam.Get(NamePriceMin)
	if minPriceString != "" {
		minPrice, err := strconv.ParseFloat(minPriceString, 64)
		if err != nil {
			return nil, err
		}
		result.MinPrice = minPrice
	}

	maxPriceString := queryParam.Get(NamePriceMax)
	if maxPriceString != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceString, 64)
		if err != nil {
			return nil, err
		}
		result.MaxPrice = maxPrice
	} else {
		result.MaxPrice = maxPriceProductOnStore
	}

	minRatingString := queryParam.Get(NameRatingMin)
	if minRatingString != "" {
		minRating, err := strconv.ParseFloat(minRatingString, 64)
		if err != nil {
			return nil, err
		}
		result.MinRating = minRating
	}

	maxRatingString := queryParam.Get(NameRatingMax)
	if maxRatingString != "" {
		maxRating, err := strconv.ParseFloat(maxRatingString, 64)
		if err != nil {
			return nil, err
		}
		result.MaxRating = maxRating
	} else {
		result.MaxRating = maxRatingOnStore
	}

	orderByString := queryParam.Get(NameOrder)
	if orderByString != "" && orderByString != TypeOrderPrice && orderByString != TypeOrderRating {
		return nil, errors.New(customErrors.BAD_QUERY_PARAM_DESCR)
	}
	if orderByString == "" {
		result.OrderBy = TypeOrderRating
	} else {
		result.OrderBy = orderByString
	}

	orderByTypeString := queryParam.Get(NameOrderType)
	if orderByTypeString != "" && orderByTypeString != TypeOrderMin && orderByTypeString != TypeOrderMax {
		return nil, errors.New(customErrors.BAD_QUERY_PARAM_DESCR)
	}
	if orderByTypeString == "" {
		result.TypeOrder = TypeOrderMin
	} else {
		result.TypeOrder = orderByTypeString
	}

	category := queryParam.Get(NameCategory)
	if category == "" {
		result.NameCategory = TypeCategoryDefault
	} else {
		result.NameCategory = category
	}
	return &result, nil
}
