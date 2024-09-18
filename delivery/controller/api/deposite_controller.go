package api

import (
	"fmt"
	"polen/delivery/middleware"
	"polen/model/dto"
	"polen/usecase"
	"polen/utils/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepositeController struct {
	depoUc usecase.DepositeUseCase
	rg     *gin.RouterGroup
}

// @Summary New deposito
// @Description Create New deposito Data
// @Tags Deposito
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.DepositeRequest true "Data deposito"
// @Success 200 {object} dto.ResponseData
// @Router /deposite [POST]
func (d *DepositeController) createHandler(c *gin.Context) {
	// get credential
	role, err := common.GetRole(c)
	if err != nil {
		if err.Error() == "unautorized" {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"message": err.Error(),
			// "message": "internal server error",
		})
		return
	}
	if role != "pemodal" {
		c.JSON(403, gin.H{
			"message": "you are not allowed",
		})
		return
	}
	ucid, err := common.GetId(c)
	if err != nil {
		if err.Error() == "unautorized" {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"message": err.Error(),
			// "message": "internal server error",
		})
		return
	}
	// getting input
	var deposite dto.DepositeRequest
	if err := c.ShouldBindJSON(&deposite); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	// creating payload
	var payload dto.DepositeDto
	payload.UserCredential.Id = ucid
	payload.InterestRate.Id = deposite.InterestRateId
	payload.DepositeAmount = deposite.Amount
	// go to usecase
	code, err := d.depoUc.CreateDeposite(payload)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}
	response := dto.ResponseMessage{
		Message: "success creating data",
	}
	c.JSON(200, response)
}

// @Summary get detail deposito user
// @Description Create New deposito Data
// @Tags Deposito
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Success 200 {object} dto.ResponseData
// @Router /deposite/user [GET]
func (d *DepositeController) getDepositeByUserHandler(c *gin.Context) {
	id, err := common.GetId(c)
	if err != nil {
		if err.Error() == "unautorized" {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"message": err.Error(),
			// "message": "internal server error",
		})
		return
	}
	code, result, err := d.depoUc.FindByUcId(id)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
			// "message": "internal server error",
		})
		return
	}
	response := dto.ResponseData{
		Message: "success getting data",
		Data:    result,
	}
	c.JSON(200, response)
}

// @Summary get deposito by id user
// @Description get deposito by id user
// @Tags Deposito
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param id path string true "Id ucredential user"
// @Success 200 {object} dto.ResponseData
// @Router /deposite/user/{id} [GET]
func (d *DepositeController) getDepositeByIdHandler(c *gin.Context) {
	// getting input
	id := c.Param("id")

	code, result, err := d.depoUc.FindById(id)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}
	response := dto.ResponseData{
		Message: "success getting data",
		Data:    result,
	}
	c.JSON(200, response)
}

// @Summary get deposito by id
// @Description get deposito by id deposito
// @Tags Deposito
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param id path string true "Id deposito"
// @Success 200 {object} dto.ResponseData
// @Router /deposite/{id} [GET]
func (d *DepositeController) getDepositeByUserIdHandler(c *gin.Context) {
	// getting input
	id := c.Param("id")

	// get credential
	role, err := common.GetRole(c)
	if err != nil {
		if err.Error() == "unautorized" {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"message": err.Error(),
			// "message": "internal server error",
		})
		return
	}
	if role != "admin" {
		c.JSON(403, gin.H{
			"message": "you are not allowed",
		})
		return
	}

	code, result, err := d.depoUc.FindByUcId(id)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}
	response := dto.ResponseData{
		Message: "success getting data",
		Data:    result,
	}
	c.JSON(200, response)
}

// @Summary Get all
// @Description get all data deposito
// @Tags Deposito
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param page path int true "page of pagination"
// @Param size path int true "size of pagination"
// @Success 200 {object} dto.ResponsePaging
// @Router /deposite/list/{page}/{size} [GET]
func (d *DepositeController) paggingHandler(c *gin.Context) {
	// Mengambil parameter dari URL
	page, _ := strconv.Atoi(c.Param("page"))
	size, _ := strconv.Atoi(c.Param("size"))

	// Memberikan nilai default jika parameter kosong
	if page == 0 {
		page = 1 // Nilai default untuk page
	}

	if size == 0 {
		size = 10 // Nilai default untuk size
	}
	payload := dto.PageRequest{
		Page: page,
		Size: size,
	}

	model, pagereturn, err := d.depoUc.Pagging(payload)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := dto.ResponsePaging{
		Message: "Success getting data",
		Data:    model,
		Paging:  pagereturn,
	}

	c.JSON(200, response)
}

// @Summary update
// @Description update data deposito
// @Tags Deposito
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Success 200 {object} dto.ResponseData
// @Router /deposite/update [PUT]
func (d *DepositeController) updateHandler(c *gin.Context) {
	// get credential
	role, err := common.GetRole(c)
	if err != nil {
		if err.Error() == "unautorized" {
			c.JSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"message": err.Error(),
			// "message": "internal server error",
		})
		return
	}
	if role != "admin" {
		c.JSON(403, gin.H{
			"message": "you are not allowed",
		})
		return
	}
	// update
	fmt.Println("Melakukan pembaruan database...")
	defer fmt.Println("Selesai melakukan pembaruan data...")
	err = d.depoUc.Update()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(200, gin.H{
				"message": "Selesai melakukan pembaruan data",
			})
		}
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
}

func (d *DepositeController) Route() {
	d.rg.POST("/deposite", middleware.AuthMiddleware(), d.createHandler)
	d.rg.GET("/deposite/user", middleware.AuthMiddleware(), d.getDepositeByUserHandler)
	d.rg.GET("/deposite/user/:id", middleware.AuthMiddleware(), d.getDepositeByUserIdHandler)
	d.rg.GET("/deposite/:id", middleware.AuthMiddleware(), d.getDepositeByIdHandler)
	d.rg.GET("/deposite/list/:page/:size", middleware.AuthMiddleware(), d.paggingHandler)
	d.rg.PUT("/deposite/update", middleware.AuthMiddleware(), d.updateHandler)
}

func NewDepositeController(depoUc usecase.DepositeUseCase, rg *gin.RouterGroup) *DepositeController {
	return &DepositeController{
		depoUc: depoUc,
		rg:     rg,
	}
}
