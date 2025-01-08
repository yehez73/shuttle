package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShuttleHandler struct {
	ShuttleService services.ShuttleServiceInterface
	DB             *sqlx.DB 
}

func NewShuttleHandler(shuttleService services.ShuttleServiceInterface) *ShuttleHandler {
	return &ShuttleHandler{
		ShuttleService: shuttleService,
	}
}

func (h *ShuttleHandler) GetShuttleTrackByParent(c *fiber.Ctx) error {
	// Ambil userUUID dari context
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok || userUUID == "" {
		log.Println("Invalid or missing userUUID")
		return utils.BadRequestResponse(c, "Invalid or missing userUUID", nil)
	}

	// Parse userUUID ke dalam UUID format
	parentUUID, err := uuid.Parse(userUUID)
	if err != nil {
		log.Println("Invalid userUUID format:", err)
		return utils.BadRequestResponse(c, "Invalid userUUID format", nil)
	}

	// Panggil service dengan parameter tambahan
	log.Println("Fetching shuttle track for parentUUID:", parentUUID)
	shuttles, err := h.ShuttleService.GetShuttleTrackByParent(parentUUID)
	if err != nil {
		log.Println("Shuttle data not found:", err)
		return utils.NotFoundResponse(c, "Shuttle data not found", nil)
	}

	// Kirim response
	log.Println("Successfully fetched shuttle data:", shuttles)
	return c.Status(http.StatusOK).JSON(shuttles)
}

func (h *ShuttleHandler) GetAllShuttleByParent(c *fiber.Ctx) error {
    userUUID, ok := c.Locals("userUUID").(string)
    if !ok || userUUID == "" {
        return utils.BadRequestResponse(c, "Invalid or missing userUUID", nil)
    }
    
    parentUUID, err := uuid.Parse(userUUID)
    if err != nil {
        return utils.BadRequestResponse(c, "Invalid userUUID format", nil)
    }
    
    // Debug log
    fmt.Println("ParentUUID:", parentUUID)

    shuttles, err := h.ShuttleService.GetAllShuttleByParent(parentUUID)
    if err != nil {
        return utils.NotFoundResponse(c, "Shuttle data not found", nil)
    }

    return c.Status(http.StatusOK).JSON(shuttles)
}

func (h *ShuttleHandler) GetAllShuttleByDriver(c *fiber.Ctx) error {
    userUUID, ok := c.Locals("userUUID").(string)
    if !ok || userUUID == "" {
        return utils.BadRequestResponse(c, "Invalid or missing userUUID", nil)
    }
    
    driverUUID, err := uuid.Parse(userUUID)
    if err != nil {
        return utils.BadRequestResponse(c, "Invalid userUUID format", nil)
    }
    
    // Debug log
    fmt.Println("ParentUUID:", driverUUID)

    shuttles, err := h.ShuttleService.GetAllShuttleByDriver(driverUUID)
    if err != nil {
        return utils.NotFoundResponse(c, "Shuttle data not found", nil)
    }

    return c.Status(http.StatusOK).JSON(shuttles)
}

func (h *ShuttleHandler) GetSpecShuttle(c *fiber.Ctx) error {
	shuttleUUIDParam := c.Params("id")
	if shuttleUUIDParam == "" {
		log.Println("Missing shuttle UUID in request URL")
		return utils.BadRequestResponse(c, "Missing shuttle UUID in request URL", nil)
	}

	shuttleUUID, err := uuid.Parse(shuttleUUIDParam)
	if err != nil {
		log.Println("Invalid shuttle UUID format:", err)
		return utils.BadRequestResponse(c, "Invalid shuttle UUID format", nil)
	}

	log.Println("Fetching shuttle spec data for shuttleUUID:", shuttleUUID)
	shuttle, err := h.ShuttleService.GetSpecShuttle(shuttleUUID)
	if err != nil {
		log.Println("Error fetching shuttle data:", err)
		return utils.NotFoundResponse(c, "Shuttle data not found", nil)
	}

	if len(shuttle) == 0 {
		log.Println("Shuttle data not found for shuttleUUID:", shuttleUUID)
		return utils.NotFoundResponse(c, "Shuttle data not found", nil)
	}

	log.Println("Successfully fetched shuttle data:", shuttle)
	return c.Status(http.StatusOK).JSON(shuttle)
}

func (h *ShuttleHandler) AddShuttle(c *fiber.Ctx) error {
	// Log: Check if userUUID is valid
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok || userUUID == "" {
		log.Println("AddShuttle: Invalid or missing userUUID")
		return utils.BadRequestResponse(c, "Invalid or missing userUUID", nil)
	}
	log.Printf("AddShuttle: userUUID found - %s", userUUID)

	username := c.Locals("user_name").(string)
	driverUUID, err := uuid.Parse(userUUID)
	if err != nil {
		log.Printf("AddShuttle: Invalid userUUID format - %s", userUUID)
		return utils.BadRequestResponse(c, "Invalid userUUID format", nil)
	}
	log.Printf("AddShuttle: Parsed driverUUID - %s", driverUUID.String())

	shuttleReq := new(dto.ShuttleRequest)
	if err := c.BodyParser(shuttleReq); err != nil {
		log.Println("AddShuttle: Invalid request data")
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}
	log.Println("AddShuttle: Parsed shuttle request")

	if shuttleReq.Status == "" {
		shuttleReq.Status = "waiting_to_be_taken_to_school"
		log.Println("AddShuttle: Set default status to 'waiting_to_be_taken_to_school'")
	}

	// Log: Validate shuttle request
	if err := utils.ValidateStruct(c, shuttleReq); err != nil {
		log.Printf("AddShuttle: Validation error - %s", err.Error())
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}
	log.Println("AddShuttle: Shuttle request validated")

	// Log: Attempt to add shuttle
	if err := h.ShuttleService.AddShuttle(*shuttleReq, driverUUID.String(), username); err != nil {
		log.Println("AddShuttle: Failed to add shuttle")
		return utils.InternalServerErrorResponse(c, "Failed to add shuttle", nil)
	}
	log.Println("AddShuttle: Shuttle added successfully")

	// Log: Return success response
	return utils.SuccessResponse(c, "Shuttle added successfully", nil)
}

func (h *ShuttleHandler) EditShuttle(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequestResponse(c, "Missing shuttleUUID in URL", nil)
	}

	var statusReq struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.BodyParser(&statusReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", nil)
	}

	if err := utils.ValidateStruct(c, statusReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid status: "+err.Error(), nil)
	}

	if err := h.ShuttleService.EditShuttleStatus(id, statusReq.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundResponse(c, "Shuttle not found", nil)
		}
		return utils.InternalServerErrorResponse(c, "Failed to edit shuttle", nil)
	}

	return utils.SuccessResponse(c, "Shuttle status updated successfully", nil)
}
