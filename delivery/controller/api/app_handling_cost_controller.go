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

type AppHandlingCostController struct {
	appHandlingCostUC usecase.AppHandlingCostUsecase
	rg                *gin.RouterGroup
}

// @Summary New Handling Cost
// @Description Create New Apps Handling Cost Data
// @Tags Handling Cost
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.AppHandlingCostReq true "Data application handling cost"
// @Success 200 {object} dto.ResponseData
// @Router /apphandlingcost [POST]
func (p *AppHandlingCostController) createHandler(c *gin.Context) {
	var app dto.AppHandlingCostReq
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	var apps model.AppHandlingCost
	apps.Name = app.Name
	apps.Nominal = app.Nominal
	apps.Unit = app.Unit

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
	apps.Id = uuid.NewString()
	code, err := p.appHandlingCostUC.CreateNew(apps)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := dto.ResponseData{
		Message: "successfully update verification user",
		Data:    apps,
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Get all Handling Cost
// @Description get all Apps Handling Cost Data
// @Tags Handling Cost
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param page path int true "page of pagination"
// @Param size path int true "size of pagination"
// @Success 200 {object} dto.ResponsePaging
// @Router /apphandlingcost/list/{page}/{size} [GET]
func (p *AppHandlingCostController) paggingHandler(c *gin.Context) {
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

	model, pagereturn, err := p.appHandlingCostUC.Pagging(payload)
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

func (p *AppHandlingCostController) Route() {
	p.rg.POST("/apphandlingcost", middleware.AuthMiddleware(), p.createHandler)
	p.rg.GET("/apphandlingcost/list/:page/:size", middleware.AuthMiddleware(), p.paggingHandler)
	p.rg.PUT("/apphandlingcost/", middleware.AuthMiddleware(), p.updateHandler)
	p.rg.DELETE("/apphandlingcost/:id", middleware.AuthMiddleware(), p.deleteHandler)

}

// @Summary delete data Handling Cost
// @Description Create New Apps Handling Cost Data
// @Tags Handling Cost
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param id path string true "id data"
// @Success 200 {object} dto.ResponsePaging
// @Router /apphandlingcost/{id} [DELETE]
func (p *AppHandlingCostController) deleteHandler(c *gin.Context) {
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

	err = p.appHandlingCostUC.DeleteById(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := dto.ResponseMessage{
		Message: "successfully delete app handling cost",
	}
	c.JSON(200, response)
}

// @Summary Update Handling Cost
// @Description Create New Apps Handling Cost Data
// @Tags Handling Cost
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body model.AppHandlingCost true "Data application handling cost"
// @Success 200 {object} dto.ResponseMessage
// @Router /apphandlingcost [PUT]
func (p *AppHandlingCostController) updateHandler(c *gin.Context) {
	var app model.AppHandlingCost
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

	err = p.appHandlingCostUC.Update(app)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	response := dto.ResponseMessage{
		Message: "successfully update",
	}
	c.JSON(200, response)
}

func NewAppHandlingCostController(appHandlingCostUC usecase.AppHandlingCostUsecase, rg *gin.RouterGroup) *AppHandlingCostController {
	return &AppHandlingCostController{
		appHandlingCostUC: appHandlingCostUC,
		rg:                rg,
	}
}
