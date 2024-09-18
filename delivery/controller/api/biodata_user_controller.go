package api

import (
	"polen/delivery/middleware"
	"polen/model/dto"
	"polen/usecase"
	"polen/utils/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BiodataUserController struct {
	biodataUC usecase.BiodataUserUseCase
	rg        *gin.RouterGroup
}

// @Summary Detail Biodata User
// @Description Get detail Biodata
// @Tags Biodata
// @Accept json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Produce json
// @Success 200 {object} dto.BiodataResponse
// @Router /biodata [GET]
func (u *BiodataUserController) listHandler(c *gin.Context) {
	biodata, err := u.biodataUC.FindByUserCredential(c)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, biodata)
}

// @Summary Find User Updated
// @Description Find User Which has Been Updated the Biodata
// @Tags Biodata
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Success 200 {object} dto.ResponseData
// @Router /biodata/updated [GET]
func (u *BiodataUserController) listUserUpdated(c *gin.Context) {
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
	biodata, err := u.biodataUC.FindUserUpdated()
	if err != nil {
		c.JSON(500, gin.H{
			// "message": err.Error(),
			"message": "internal server error",
		})
		return
	}
	response := dto.ResponseData{
		Message: "Successfully Getting Data",
		Data:    biodata,
	}

	c.JSON(200, response)
}

// @Summary Find User Updated
// @Description Find User Which has Been Updated the Biodata
// @Tags Biodata
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.UpdateBioRequest true "data vefication user"
// @Success 200 {object} dto.ResponseMessage
// @Router /biodata/verified [PUT]
func (u *BiodataUserController) updateAdmin(c *gin.Context) {
	var req dto.UpdateBioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, gin.H{
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
	code, err := u.biodataUC.AdminUpdate(req, c)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
			// "message": "Internal Server Error",
		})
		return
	}
	response := dto.ResponseMessage{
		Message: "successfully update verification user",
	}
	c.JSON(code, response)
}

// @Summary Update Biodata
// @Description Do Update User Biodata
// @Tags Biodata
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.BiodataRequest true "data biodata user"
// @Success 200 {object} dto.ResponseMessage
// @Router /biodata/update [PUT]
func (u *BiodataUserController) updateUser(c *gin.Context) {
	var biodata dto.BiodataRequest
	if err := c.ShouldBindJSON(&biodata); err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}

	code, err := u.biodataUC.UserUpdate(biodata, c)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}
	response := dto.ResponseMessage{
		Message: "successfully update biodata",
	}
	c.JSON(code, response)
}

// @Summary Get all data
// @Description Do Update User Biodata
// @Tags Biodata
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param page path int true "page of pagination"
// @Param size path int true "size of pagination"
// @Success 200 {object} dto.ResponsePaging
// @Router /biodata/list/{page}/{size} [GET]
func (b *BiodataUserController) paggingBiodataHandler(c *gin.Context) {
	// get role
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

	model, pagereturn, err := b.biodataUC.Paging(payload)
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

func (b *BiodataUserController) Route() {
	// b.rg.PUT("/biodata/update", middleware.AuthMiddleware(), b.createHandler)
	b.rg.GET("/biodata", middleware.AuthMiddleware(), b.listHandler)
	b.rg.GET("/biodata/updated", middleware.AuthMiddleware(), b.listUserUpdated)
	b.rg.PUT("/biodata/update", middleware.AuthMiddleware(), b.updateUser)
	b.rg.PUT("/biodata/verified", middleware.AuthMiddleware(), b.updateAdmin)
	b.rg.GET("/biodata/list/:page/:size", middleware.AuthMiddleware(), b.paggingBiodataHandler)
}

func NewBiodataController(biodataUC usecase.BiodataUserUseCase, rg *gin.RouterGroup) *BiodataUserController {
	return &BiodataUserController{
		biodataUC: biodataUC,
		rg:        rg,
	}
}
