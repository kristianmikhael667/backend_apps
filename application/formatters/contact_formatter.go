package formatters

import (
	"backend_apps/models"
	util "backend_apps/package"
)

type (
	HeadResponseContact struct {
		Message string                 `json:"message"`
		Meta    map[string]interface{} `json:"meta"`
		Data    interface{}            `json:"data"`
	}

	RequestBodyContact struct {
		Phone      string `json:"phone" validate:"required,min=10"`
		ProviderId string `json:"provider_id" validate:"required"`
	}

	ResponseContact struct {
		Providers   Providers `json:"provider"`
		ContactId   string    `json:"contact_id"`
		Phone       string    `json:"phone"`
		ProviderId  int16     `json:"provider_id"`
		Status      int       `json:"status"`
		GanjilGenep string    `json:"ganjil_genep"`
		CreatedAt   string    `json:"created_at"`
		UpdatedAt   string    `json:"updated_at"`
	}

	Providers struct {
		ProviderId   string `json:"provider_id"`
		PrefixCode   string `json:"prefix_code"`
		NameProvider string `json:"name_provider"`
		Status       int    `json:"status"`
	}
)

func FormatContact(contact models.Contact) ResponseContact {
	decryptedData, _ := util.Decrypt(contact.Phone, "key")
	response := ResponseContact{
		Providers: Providers{
			ProviderId:   contact.Provider.ProviderId.String(),
			PrefixCode:   contact.Provider.PrefixCode,
			NameProvider: contact.Provider.NameProvider,
			Status:       contact.Provider.Status,
		},
		ContactId:   contact.ContactId.String(),
		Phone:       decryptedData,
		ProviderId:  contact.ProviderId,
		Status:      contact.Status,
		GanjilGenep: contact.GanjilGenep,
		CreatedAt:   contact.ContactId.String(),
		UpdatedAt:   contact.UpdatedAt.String(),
	}
	return response
}
