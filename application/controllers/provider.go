package controllers

import (
	"backend_apps/database"
	"backend_apps/formatters"
	"backend_apps/models"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateProvider(c *fiber.Ctx) error {
	payload := new(formatters.RequestBodyProvider)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Validation Error"})
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Field is required"})
			} else if err.Tag() == "min" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Field min character"})
			}
		}
	}

	newProvider := models.Provider{
		PrefixCode:   payload.PrefixCode,
		NameProvider: payload.NameProvider,
		Status:       payload.Status,
	}

	if err := database.GetConnection().Save(&newProvider).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	response := formatters.FormatProvider(newProvider)
	var responseCreated formatters.HeadResponseProvider
	responseCreated.Message = "Success Create Provider"
	responseCreated.Meta = nil
	responseCreated.Data = response

	return c.Status(fiber.StatusCreated).JSON(responseCreated)
}

func GetProvider(c *fiber.Ctx) error {
	// Parsing page and pageSize from query parameters
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	sortBy := c.Query("sortBy")
	order := c.Query("d") //DESC or ASC

	if sortBy == "" {
		sortBy = "id"
	}

	if c.Query("d") == "" {
		order = "desc"
	}

	/* no pagination */
	if c.Query("page") == "" && c.Query("pageSize") == "" {
		resp, err := GetProviderMassive(c)
		if err != nil {
			return err
		}
		// Response Created
		var responseCreated formatters.HeadResponseProvider
		responseCreated.Message = "Success Get Provider"
		responseCreated.Meta = nil
		responseCreated.Data = resp

		return c.Status(fiber.StatusOK).JSON(responseCreated)
	}

	if pageSize == 0 {
		pageSize = 10
	}

	maxpage := float64(CalculateProvider()) / float64(pageSize)
	ceiled := math.Ceil(maxpage)
	if page > int(ceiled) {
		page = int(ceiled)
	}

	if page < 1 {
		page = 1
	}

	// Fetch paginated groups
	resp, err := GetPaginatedProvider(c, page, pageSize, sortBy, order)

	if err != nil {
		return err
	}

	extrameta := map[string]interface{}{
		"currentPage": page,
		"pageSize":    pageSize,
		"pageCount":   ceiled,
		"sortBy":      sortBy,
		"order":       order,
	}

	var responseCreated formatters.HeadResponseProvider
	responseCreated.Message = "Success Get Provider"
	responseCreated.Meta = extrameta
	responseCreated.Data = resp

	return c.Status(fiber.StatusOK).JSON(responseCreated)
}

func GetDetailProvider(c *fiber.Ctx) error {
	// Check Payload
	payload := new(formatters.RequestByUid)
	if err := c.ParamsParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Get Detail Provider
	var provider models.Provider
	query := database.GetConnection()
	query = query.Model(&models.Provider{}).Where("provider_id = ?", payload.ID).First(&provider)

	if query.Error == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Not Found Provider ID"})
	}

	// Response
	response := formatters.FormatProvider(provider)
	var responseCreated formatters.HeadResponseProvider
	responseCreated.Message = "Success Get Detail Provider"
	responseCreated.Meta = nil
	responseCreated.Data = response
	return c.Status(fiber.StatusOK).JSON(responseCreated)
}

func GetDetailProviders(uid string, c *fiber.Ctx) (models.Provider, error) {
	// Get Detail Provider
	var provider models.Provider
	query := database.GetConnection()
	query = query.Model(&models.Provider{}).Where("provider_id = ?", uid).First(&provider)

	if query.Error == gorm.ErrRecordNotFound {
		return provider, query.Error
	}
	return provider, nil
}

func UpdateProvider(c *fiber.Ctx) error {
	// Check Payload
	payload := new(formatters.RequestByUid)
	if err := c.ParamsParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	var providers models.Provider

	// Parse request body
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Update user data
	if err := database.GetConnection().Model(&models.Provider{}).Where("provider_id = ?", payload.ID).Updates(requestBody).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Provider not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update provider"})
	}

	// Check changes
	if err := database.GetConnection().Where("provider_id = ?", payload.ID).First(&providers).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Provider not found"})
	}

	// Serialize the response data
	var responseUpdate formatters.HeadResponseProvider
	responseUpdate.Message = "Success Get Provider"
	responseUpdate.Meta = nil
	responseUpdate.Data = providers

	return c.Status(fiber.StatusOK).JSON(responseUpdate)
}

func DeleteProvider(c *fiber.Ctx) error {
	// Check Payload
	payload := new(formatters.RequestByUid)
	if err := c.ParamsParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	var providers models.Provider
	resProvider := database.GetConnection().Where("provider_id = ?", payload.ID).First(&providers)
	if resProvider.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Provider Not Found",
		})
	}

	var contact models.Contact
	resDelInvItems := database.GetConnection().Where("provider_id = ?", providers.ID).Delete(&contact)
	if resDelInvItems.RowsAffected == 0 {
		log.Println("Contact item has been deleted")
	}

	resDelParInv := database.GetConnection().Where("id = ?", providers.ID).Delete(&providers)
	if resDelParInv.RowsAffected == 0 {
		log.Print("Provider has been deleted")
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success has been deleted provider",
	})
}

/*------------------------- Controllers Item ------------------------*/

/* Get Provider without Pagination */
func GetProviderMassive(c *fiber.Ctx) ([]formatters.ResponseProvider, error) {
	var providers []models.Provider
	if err := database.GetConnection().Find(&providers).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Empty Provider"})
		}
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	providerFormatteds := []formatters.ResponseProvider{}
	var providerFormatted formatters.ResponseProvider

	for _, provider := range providers {
		providerFormatted = formatters.FormatProvider(provider)
		providerFormatteds = append(providerFormatteds, providerFormatted)
	}
	return providerFormatteds, nil
}

/* Pagination Provider */
func GetPaginatedProvider(c *fiber.Ctx, offset, limit int, sortBy, order string) ([]formatters.ResponseProvider, error) {
	var providers []models.Provider
	db := database.GetConnection().Model(&providers)
	page := (offset - 1) * limit
	order_by := fmt.Sprintf("%s %s", sortBy, order)

	if err := db.Offset(page).Limit(limit).Order(order_by).Find(&providers).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Empty Provider"})
		}
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	var dataProvider formatters.ResponseProvider

	datacs := []formatters.ResponseProvider{}
	for _, provider := range providers {
		dataProvider = formatters.FormatProvider(provider)
		datacs = append(datacs, dataProvider)
	}
	return datacs, nil
}

/* Calculate Count Provider */
func CalculateProvider() int64 {
	var count int64
	var providers []models.Provider
	result := database.GetConnection().Model(&providers).Count(&count)
	if result.Error != nil {
		log.Print("Error occurred : ", result.Error)
		return 0
	}
	return count
}
