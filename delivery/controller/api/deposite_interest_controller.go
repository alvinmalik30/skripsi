package api

import (
	"net/http"
	"polen/delivery/middleware"
	"polen/model/dto"
	"polen/usecase"
	"polen/utils/common"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DepositeInterestController struct {
	depositeUC usecase.DepositeInterestUseCase
	rg         *gin.RouterGroup
}

// @Summary new
// @Description create new data deposito interest
// @Tags deposito interest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.DepositeInterestReq true "Data deposite interest"
// @Success 200 {object} dto.ResponseData
// @Router /depositeinterest [POST]
func (d *DepositeInterestController) createHandler(c *gin.Context) {
	var payload dto.DepositeInterestReq
	var deposite dto.DepositeInterestRequest
	if err := c.ShouldBindJSON(&deposite); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

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
	deposite.Id = uuid.NewString()
	deposite.InterestRate = payload.InterestRate
	deposite.DurationMounth = payload.DurationMounth
	deposite.TaxRate = payload.TaxRate
	code, err := d.depositeUC.CreateNew(deposite)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}

	resutl := dto.ResponseData{
		Message: "success creating data",
		Data:    deposite,
	}
	c.JSON(http.StatusCreated, resutl)
}

// @Summary Get all
// @Description get all data deposito interest
// @Tags deposito interest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param page path int true "page of pagination"
// @Param size path int true "size of pagination"
// @Success 200 {object} dto.ResponsePaging
// @Router /depositeinterest/list/{page}/{size} [GET]
func (d *DepositeInterestController) paggingHandler(c *gin.Context) {
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

	model, pagereturn, err := d.depositeUC.Pagging(payload)
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

func (d *DepositeInterestController) Route() {
	d.rg.POST("/depositeinterest", middleware.AuthMiddleware(), d.createHandler)
	d.rg.GET("/depositeinterest/list/:page/:size", middleware.AuthMiddleware(), d.paggingHandler)
	d.rg.PUT("/depositeinterest/", middleware.AuthMiddleware(), d.updateHandler)
	d.rg.DELETE("/depositeinterest/:id", middleware.AuthMiddleware(), d.deleteHandler)

}

// @Summary delete
// @Description delete data deposito interest
// @Tags deposito interest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param id path string true "id"
// @Success 200 {object} dto.ResponseMessage
// @Router /depositeinterest/{id} [DELETE]
func (d *DepositeInterestController) deleteHandler(c *gin.Context) {
	id := c.Param("id")

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

	err = d.depositeUC.DeleteById(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := dto.ResponseMessage{
		Message: "successfully delete deposite",
	}
	c.JSON(http.StatusCreated, result)
}

// @Summary update
// @Description update data deposito interest
// @Tags deposito interest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.DepositeInterestRequest true "Data deposite interest"
// @Success 200 {object} dto.ResponseData
// @Router /depositeinterest [PUT]
func (d *DepositeInterestController) updateHandler(c *gin.Context) {
	var deposite dto.DepositeInterestRequest
	if err := c.ShouldBindJSON(&deposite); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

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

	err = d.depositeUC.Update(deposite)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	result := dto.ResponseMessage{
		Message: "successfully update",
	}
	c.JSON(http.StatusCreated, result)
}
func NewDepositeInterestController(depositeinterestUC usecase.DepositeInterestUseCase, rg *gin.RouterGroup) *DepositeInterestController {
	return &DepositeInterestController{
		depositeUC: depositeinterestUC,
		rg:         rg,
	}
}
