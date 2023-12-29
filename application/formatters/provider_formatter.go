package formatters

import "backend_apps/models"

type (
	RequestByUid struct {
		ID string `param:"id" validate:"required"`
	}

	HeadResponseProvider struct {
		Message string                 `json:"message"`
		Meta    map[string]interface{} `json:"meta"`
		Data    interface{}            `json:"data"`
	}

	RequestBodyProvider struct {
		PrefixCode   string `json:"prefix_code" validate:"required,min=2"`
		NameProvider string `json:"name_provider"  validate:"required,min=2"`
		Status       int    `json:"status"  validate:"required"`
	}

	ResponseProvider struct {
		ProviderId   string `json:"provider_id"`
		PrefixCode   string `json:"prefix_code"`
		NameProvider string `json:"name_provider"`
		Status       int    `json:"status"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	}
)

func FormatProvider(provider models.Provider) ResponseProvider {
	response := ResponseProvider{
		ProviderId:   provider.ProviderId.String(),
		PrefixCode:   provider.PrefixCode,
		NameProvider: provider.NameProvider,
		Status:       provider.Status,
		CreatedAt:    provider.CreatedAt.String(),
		UpdatedAt:    provider.UpdatedAt.String(),
	}
	return response
}
