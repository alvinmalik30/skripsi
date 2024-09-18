package api

import (
	"net/http"
	"polen/delivery/middleware"
	"polen/model"
	"polen/model/dto"
	"polen/usecase"
	"polen/utils/common"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LatePaymentFeeController struct {
	LatePaymentFeeUC usecase.LatePaymentFeeUsecase
	rg               *gin.RouterGroup
}

// @Summary new
// @Description create new data late fee payment
// @Tags late fee payment
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.LatePaymentFeeReq true "Data late fee"
// @Success 200 {object} dto.ResponseData
// @Router /latepaymentfee [POST]
func (p *LatePaymentFeeController) createHandler(c *gin.Context) {
	var app dto.LatePaymentFeeReq
	var payload model.LatePaymentFee
	payload.Name = app.Name
	payload.Nominal = app.Nominal
	payload.Unit = app.Unit
	if err := c.ShouldBindJSON(&app); err != nil {
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
	payload.Id = uuid.NewString()
	code, err := p.LatePaymentFeeUC.CreateNew(payload)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := dto.ResponseData{
		Message: "success creating data",
		Data:    payload,
	}
	c.JSON(http.StatusCreated, result)
}

// @Summary get all
// @Description get all data late fee payment
// @Tags late fee payment
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param page path int true "page of pagination"
// @Param size path int true "size of pagination"
// @Success 200 {object} dto.ResponsePaging
// @Router /latepaymentfee/list/{page}/{size} [GET]
func (p *LatePaymentFeeController) paggingHandler(c *gin.Context) {
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

	model, pagereturn, err := p.LatePaymentFeeUC.Pagging(payload)
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

func (p *LatePaymentFeeController) Route() {
	p.rg.POST("/latepaymentfee", middleware.AuthMiddleware(), p.createHandler)
	p.rg.GET("/latepaymentfee/list/:page/:size", middleware.AuthMiddleware(), p.paggingHandler)
	p.rg.PUT("/latepaymentfee/", middleware.AuthMiddleware(), p.updateHandler)
	p.rg.DELETE("/latepaymentfee/:id", middleware.AuthMiddleware(), p.deleteHandler)

}

// @Summary delete
// @Description delete data late fee payment
// @Tags late fee payment
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param id path string true "id"
// @Success 200 {object} dto.ResponseMessage
// @Router /latepaymentfee/{id} [DELETE]
func (p *LatePaymentFeeController) deleteHandler(c *gin.Context) {
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

	err = p.LatePaymentFeeUC.DeleteById(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := dto.ResponseMessage{
		Message: "success creating data",
	}
	c.JSON(http.StatusCreated, result)
}

// @Summary update
// @Description update data late fee payment
// @Tags late fee payment
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body model.LatePaymentFee true "Data late fee"
// @Success 200 {object} dto.ResponseMessage
// @Router /latepaymentfee [PUT]
func (p *LatePaymentFeeController) updateHandler(c *gin.Context) {
	var app model.LatePaymentFee
	if err := c.ShouldBindJSON(&app); err != nil {
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

	err = p.LatePaymentFeeUC.Update(app)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	result := dto.ResponseMessage{
		Message: "success creating data",
	}
	c.JSON(http.StatusCreated, result)
}

func NewLatePaymentFeeController(LatePaymentFeeUC usecase.LatePaymentFeeUsecase, rg *gin.RouterGroup) *LatePaymentFeeController {
	return &LatePaymentFeeController{
		LatePaymentFeeUC: LatePaymentFeeUC,
		rg:               rg,
	}
}
