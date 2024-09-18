package api

import (
	"polen/delivery/middleware"
	"polen/model/dto"
	"polen/usecase"
	"polen/utils/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SaldoController struct {
	saldoUC usecase.SaldoUsecase
	rg      *gin.RouterGroup
}

// @Summary get all
// @Description get all data saldo
// @Tags saldo
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param page path int true "page of pagination"
// @Param size path int true "size of pagination"
// @Success 200 {object} dto.ResponsePaging
// @Router /saldo/list/{page}/{size} [GET]
func (s *SaldoController) paggingHandler(c *gin.Context) {
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

	model, pagereturn, err := s.saldoUC.Pagging(payload)
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

// @Summary get by id installenment
// @Description get by id data loan
// @Tags saldo
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Success 200 {object} dto.ResponseData
// @Router /saldo [GET]
func (s *SaldoController) showSaldoUser(c *gin.Context) {
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

	model, err := s.saldoUC.FindByIdUser(ucid)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := dto.ResponseData{
		Message: "Success getting data",
		Data:    model,
	}

	c.JSON(200, response)
}

// @Summary get by id installenment
// @Description get by id data loan
// @Tags saldo
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param id path int true "id"
// @Success 200 {object} dto.ResponseData
// @Router /saldo/{id} [GET]
func (s *SaldoController) showSaldoById(c *gin.Context) {
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

	id := c.Param("id")

	model, err := s.saldoUC.FindByIdUser(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := dto.ResponseData{
		Message: "Success getting data",
		Data:    model,
	}

	c.JSON(200, response)
}

func (s *SaldoController) Route() {
	s.rg.GET("/saldo", middleware.AuthMiddleware(), s.showSaldoUser)
	s.rg.GET("/saldo/:id", middleware.AuthMiddleware(), s.showSaldoById)
	s.rg.GET("/saldo/list/:page/:size", middleware.AuthMiddleware(), s.paggingHandler)
}

func NewSaldoController(saldoUC usecase.SaldoUsecase, rg *gin.RouterGroup) *SaldoController {
	return &SaldoController{
		saldoUC: saldoUC,
		rg:      rg,
	}
}
