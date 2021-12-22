package notifications

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"strings"
)

func FromModelAddressToString(address models.Address) string {
	addressString := strings.Builder{}

	if address.Country != "" {
		addressString.WriteString("Страна: " + address.Country + ", ")
	}
	if address.Region != "" {
		addressString.WriteString("Регион: " + address.Region + ", ")
	}
	if address.City != "" {
		addressString.WriteString("город: " + address.City + ", ")
	}
	if address.Street != "" {
		addressString.WriteString("улица: " + address.Street + ", ")
	}
	if address.House != "" {
		addressString.WriteString("дом: " + address.House + ", ")
	}
	if address.Flat != "" {
		addressString.WriteString("квартира: " + address.Flat + ", ")
	}
	if address.Index != "" {
		addressString.WriteString("индекс: " + address.Index)
	}

	return addressString.String()
}
