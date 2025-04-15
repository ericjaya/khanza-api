package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ikti-its/khanza-api/internal/app/exception"
	web "github.com/ikti-its/khanza-api/internal/app/model"
	"github.com/ikti-its/khanza-api/internal/modules/tindakan/internal/model"
	"github.com/ikti-its/khanza-api/internal/modules/tindakan/internal/usecase"
)

type TindakanController struct {
	UseCase *usecase.TindakanUseCase
}

func NewTindakanController(useCase *usecase.TindakanUseCase) *TindakanController {
	return &TindakanController{
		UseCase: useCase,
	}
}

func (c *TindakanController) Create(ctx *fiber.Ctx) error {
	var request model.TindakanRequest
	fmt.Println("📥 Received POST /tindakan")

	if err := ctx.BodyParser(&request); err != nil {
		fmt.Println("❌ Error parsing body:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(web.Response{
			Code:   fiber.StatusBadRequest,
			Status: "Bad Request",
			Data:   "Invalid request body",
		})
	}

	response, err := c.UseCase.Create(&request)
	if err != nil {
		fmt.Println("❌ Error in usecase.Create():", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(web.Response{
			Code:   fiber.StatusInternalServerError,
			Status: "Error",
			Data:   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(web.Response{
		Code:   fiber.StatusCreated,
		Status: "Created",
		Data:   response,
	})
}

func (c *TindakanController) GetAll(ctx *fiber.Ctx) error {
	response, err := c.UseCase.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(web.Response{
			Code:   fiber.StatusInternalServerError,
			Status: "Error",
			Data:   err.Error(),
		})
	}
	return ctx.JSON(web.Response{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (c *TindakanController) GetByNomorRawat(ctx *fiber.Ctx) error {
	nomorRawat := ctx.Params("nomor_rawat")

	response, err := c.UseCase.GetByNomorRawat(nomorRawat)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(web.Response{
			Code:   fiber.StatusNotFound,
			Status: "Not Found",
			Data:   err.Error(),
		})
	}

	return ctx.JSON(web.Response{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (c *TindakanController) Update(ctx *fiber.Ctx) error {
	nomorRawat := ctx.Params("nomor_rawat")
	var request model.TindakanRequest

	if err := ctx.BodyParser(&request); err != nil {
		panic(&exception.BadRequestError{Message: "Invalid request body"})
	}

	response, err := c.UseCase.Update(nomorRawat, &request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(web.Response{
			Code:   fiber.StatusInternalServerError,
			Status: "Error",
			Data:   err.Error(),
		})
	}
	return ctx.JSON(web.Response{
		Code:   fiber.StatusOK,
		Status: "OK",
		Data:   response,
	})
}

func (c *TindakanController) Delete(ctx *fiber.Ctx) error {
	nomorRawat := ctx.Params("nomor_rawat")
	jamRawat := ctx.Params("jam_rawat")

	if nomorRawat == "" || jamRawat == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"status": "Bad Request",
			"data":   "nomor_rawat and jam_rawat are required",
		})
	}

	err := c.UseCase.Delete(nomorRawat, jamRawat)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Error",
			"data":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"code":   200,
		"status": "Success",
		"data":   "Tindakan deleted successfully",
	})
}
