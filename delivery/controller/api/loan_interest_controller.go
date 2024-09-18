package api

import (
	"fmt"
	"net/http"
	"polen/delivery/middleware"
	"polen/model"
	"polen/model/dto"
	"polen/usecase"
	"polen/utils/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoanInterestController struct {
	loanInterestUC usecase.LoanInterestUseCase
	rg             *gin.RouterGroup
}

// @Summary new
// @Description create new data loan interest
// @Tags loan interest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.LoanInterest true "Data loan interest"
// @Success 200 {object} dto.ResponseData
// @Router /loaninterest [POST]
func (l *LoanInterestController) createHandler(c *gin.Context) {
	var data dto.LoanInterest
	var loanInterest model.LoanInterest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	loanInterest.LoanInterestRate = data.LoanInterestRate
	loanInterest.DurationMonths = data.DurationMonths

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
	loanInterest.Id = common.GenerateID()
	fmt.Println(loanInterest)
	code, err := l.loanInterestUC.CreateNew(loanInterest)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}

	respone := dto.ResponseData{
		Message: "success creating data",
		Data:    loanInterest,
	}
	c.JSON(http.StatusCreated, respone)
}

// @Summary get all
// @Description get all data loan interest
// @Tags loan interest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param page path int true "page of pagination"
// @Param size path int true "size of pagination"
// @Success 200 {object} dto.ResponsePaging
// @Router /loaninterest/list/{page}/{size} [GET]
func (l *LoanInterestController) paggingHandler(c *gin.Context) {
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

	model, pagereturn, err := l.loanInterestUC.Pagging(payload)
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

func (l *LoanInterestController) Route() {
	l.rg.POST("/loaninterest", middleware.AuthMiddleware(), l.createHandler)
	l.rg.GET("/loaninterest/list/:page/:size", middleware.AuthMiddleware(), l.paggingHandler)
	l.rg.PUT("/loaninterest", middleware.AuthMiddleware(), l.updateHandler)
	l.rg.DELETE("/loaninterest/:id", middleware.AuthMiddleware(), l.deleteHandler)

}

// @Summary get all
// @Description get all data loan interest
// @Tags loan interest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param id path string true "id"
// @Success 200 {object} dto.ResponsePaging
// @Router /loaninterest/{id} [DELETE]
func (l *LoanInterestController) deleteHandler(c *gin.Context) {
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

	err = l.loanInterestUC.DeleteById(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	respone := dto.ResponseMessage{
		Message: "successfully delete loan interest",
	}
	c.JSON(http.StatusCreated, respone)
}

// @Summary update
// @Description create update data loan interest
// @Tags loan interest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body model.LoanInterest true "Data loan interest"
// @Success 200 {object} dto.ResponseMessage
// @Router /loaninterest [PUT]
func (l *LoanInterestController) updateHandler(c *gin.Context) {
	var loan model.LoanInterest
	if err := c.ShouldBindJSON(&loan); err != nil {
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

	err = l.loanInterestUC.Update(loan)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	respone := dto.ResponseMessage{
		Message: "success update data",
	}
	c.JSON(http.StatusCreated, respone)
}

func NewLoanInterestController(loanInterestUC usecase.LoanInterestUseCase, rg *gin.RouterGroup) *LoanInterestController {
	return &LoanInterestController{
		loanInterestUC: loanInterestUC,
		rg:             rg,
	}
}
