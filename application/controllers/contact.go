package controllers

import (
	"backend_apps/database"
	"backend_apps/formatters"
	"backend_apps/models"
	util "backend_apps/package"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateContact(c *fiber.Ctx) error {
	payload := new(formatters.RequestBodyContact)
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

	// Get Provider
	provider, err := GetDetailProviders(payload.ProviderId, c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Not Found Provider"})
	}

	key := util.Getenv("SECRET_ECRYPT", "")
	encryptedData, err := util.Encrypt(payload.Phone, key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	number, _ := strconv.Atoi(payload.Phone)

	var ganjilGenep string
	if number%2 == 0 {
		ganjilGenep = "Genap"
	} else {
		ganjilGenep = "Ganjil"
	}

	newContact := models.Contact{
		Phone:       encryptedData,
		ProviderId:  provider.ID,
		Status:      1,
		GanjilGenep: ganjilGenep,
	}

	if err := database.GetConnection().Save(&newContact).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	response := formatters.FormatContact(newContact)
	var responseCreated formatters.HeadResponseContact
	responseCreated.Message = "Success Create Contact"
	responseCreated.Meta = nil
	responseCreated.Data = response
	decryptedData, _ := util.Decrypt(newContact.Phone, "key")

	Broadcast([]byte("Success Create Number Phone ? " + decryptedData))

	return c.Status(fiber.StatusCreated).JSON(responseCreated)
}

func GetContact(c *fiber.Ctx) error {
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
		resp, err := GetContactMassive(c)
		if err != nil {
			return err
		}
		// Response Created
		var responseCreated formatters.HeadResponseContact
		responseCreated.Message = "Success Get Contact"
		responseCreated.Meta = nil
		responseCreated.Data = resp
		return c.Status(fiber.StatusOK).JSON(responseCreated)
	}

	if pageSize == 0 {
		pageSize = 10
	}

	maxpage := float64(CalculateContact()) / float64(pageSize)
	ceiled := math.Ceil(maxpage)
	if page > int(ceiled) {
		page = int(ceiled)
	}

	if page < 1 {
		page = 1
	}

	// Fetch paginated groups
	resp, err := GetPaginatedContact(c, page, pageSize, sortBy, order)

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

	var responseCreated formatters.HeadResponseContact
	responseCreated.Message = "Success Get Contact"
	responseCreated.Meta = extrameta
	responseCreated.Data = resp

	// Broadcast([]byte("Data baru ditambahkan"))

	return c.Status(fiber.StatusOK).JSON(responseCreated)
}

func GetDetailContact(c *fiber.Ctx) error {
	// Check Payload
	payload := new(formatters.RequestByUid)
	if err := c.ParamsParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Get Detail Invoice
	var contact models.Contact
	query := database.GetConnection()
	query = query.Preload("Provider").Model(&models.Contact{}).Where("contact_id = ?", payload.ID).First(&contact)

	if query.Error == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Not Found Contact ID"})
	}

	// Response
	response := formatters.FormatContact(contact)
	var responseCreated formatters.HeadResponseContact
	responseCreated.Message = "Success Get Detail Contact"
	responseCreated.Meta = nil
	responseCreated.Data = response
	// Broadcast([]byte("Data baru ditambahkan"))

	return c.Status(fiber.StatusOK).JSON(responseCreated)
}

func UpdateContact(c *fiber.Ctx) error {
	// Check Payload
	payload := new(formatters.RequestByUid)
	if err := c.ParamsParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	var contacts models.Contact

	// Parse request body
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Handle Provider
	provider, _ := requestBody["provider_id"].(string)
	providers, err := GetDetailProviders(provider, c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found Provider"})
	}
	requestBody["provider_id"] = providers.ID

	// Handle Encrypt Phone Number
	phones, _ := requestBody["phone"].(string)
	number, _ := strconv.Atoi(phones)

	var ganjilGenep string
	if number%2 == 0 {
		ganjilGenep = "Genap"
	} else {
		ganjilGenep = "Ganjil"
	}

	encryptedData, err := util.Encrypt(phones, "key")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	requestBody["phone"] = encryptedData
	requestBody["ganjil_genep"] = ganjilGenep

	// Update user data
	if err := database.GetConnection().Model(&models.Contact{}).Where("contact_id = ?", payload.ID).Updates(requestBody).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Contact not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update contact"})
	}

	// Check changes
	if err := database.GetConnection().Where("contact_id = ?", payload.ID).First(&contacts).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Contacts not found"})
	}

	// Serialize the response data
	var responseUpdate formatters.HeadResponseProvider
	responseUpdate.Message = "Success Get Contacts"
	responseUpdate.Meta = nil
	responseUpdate.Data = contacts
	decryptedData, _ := util.Decrypt(contacts.Phone, "key")
	Broadcast([]byte("Success Update Number Phone ? " + decryptedData))
	return c.Status(fiber.StatusOK).JSON(responseUpdate)
}

func DeleteContact(c *fiber.Ctx) error {
	// Check Payload
	payload := new(formatters.RequestByUid)
	if err := c.ParamsParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	var contacts models.Contact
	getContact := database.GetConnection().Where("contact_id = ?", payload.ID).First(&contacts)
	if getContact.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Contact Not Found",
		})
	}

	decryptedData, _ := util.Decrypt(contacts.Phone, "key")
	Broadcast([]byte("Success Delete Number Phone ? " + decryptedData))

	contact := database.GetConnection().Where("id = ?", contacts.ID).Delete(&contacts)
	if contact.RowsAffected == 0 {
		log.Println("Contact item has been deleted")
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success has been deleted contact",
	})
}

/*------------------------- Controllers Item ------------------------*/

/* Get Contact without Pagination */
func GetContactMassive(c *fiber.Ctx) ([]formatters.ResponseContact, error) {
	var contacts []models.Contact
	if err := database.GetConnection().Preload("Provider").Model(&models.Contact{}).Find(&contacts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Empty Contacts"})
		}
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	contactFormatteds := []formatters.ResponseContact{}
	var contactFormatted formatters.ResponseContact

	for _, provider := range contacts {
		contactFormatted = formatters.FormatContact(provider)
		contactFormatteds = append(contactFormatteds, contactFormatted)
	}
	return contactFormatteds, nil
}

/* Pagination Contact */
func GetPaginatedContact(c *fiber.Ctx, offset, limit int, sortBy, order string) ([]formatters.ResponseContact, error) {
	var contacts []models.Contact
	db := database.GetConnection().Preload("Provider").Model(&contacts)
	page := (offset - 1) * limit
	order_by := fmt.Sprintf("%s %s", sortBy, order)

	if err := db.Offset(page).Limit(limit).Order(order_by).Find(&contacts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Empty Contact"})
		}
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	var dataContacts formatters.ResponseContact

	datacs := []formatters.ResponseContact{}
	for _, contact := range contacts {
		dataContacts = formatters.FormatContact(contact)
		datacs = append(datacs, dataContacts)
	}
	return datacs, nil
}

/* Calculate Count Contact */
func CalculateContact() int64 {
	var count int64
	var contacts []models.Contact
	result := database.GetConnection().Model(&contacts).Count(&count)
	if result.Error != nil {
		log.Print("Error occurred : ", result.Error)
		return 0
	}
	return count
}
