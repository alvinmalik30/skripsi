package api

import (
	"polen/delivery/middleware"
	"polen/model/dto"
	"polen/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userUC usecase.UserUseCase
	authUC usecase.AuthUseCase
	rg     *gin.RouterGroup
}

// @Summary Login
// @Description Login for User
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.AuthLoginRequest true "Data login user"
// @Success 200 {object} dto.ResponseData
// @Router /auth/login [post]
func (a *AuthController) loginHandler(c *gin.Context) {
	var data dto.AuthLoginRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	authResponse, err := a.authUC.Login(data)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": "Failed",
		})
		return
	}

	response := dto.ResponseData{
		Message: "successfully login",
		Data:    authResponse,
	}

	c.JSON(200, response)
	// c.JSON(200, response)
}

// @Summary Register
// @Description Register for User
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.AuthRequest true "Data registration user"
// @Success 200 {object} dto.ResponseMessage
// @Router /auth/register [post]
func (a *AuthController) registerHandler(c *gin.Context) {
	var auth dto.AuthRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := a.userUC.Register(auth)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := dto.ResponseMessage{
		Message: "successfully register",
	}

	c.JSON(200, gin.H{
		"message": response.Message,
	})
}

func (a *AuthController) showUserHandler(c *gin.Context) {
	name := c.Param("name")

	model, err := a.userUC.FindByUsername(name, c)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := gin.H{
		"message": "successfully getting data",
		"data": gin.H{
			"id":       model.Id,
			"username": model.Username,
			"role":     model.Role,
		},
	}

	c.JSON(200, response)
}

// @Summary Get All User
// @Description Getting All User
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param page path int true "page of pagination"
// @Param size path int true "size of pagination"
// @Success 200 {object} dto.ResponsePaging
// @Router /user/{page}/{size} [GET]
func (a *AuthController) paggingUserHandler(c *gin.Context) {
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

	model, pagereturn, err := a.userUC.Paging(payload, c)
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

func (a *AuthController) Route() {
	a.rg.POST("/auth/login", a.loginHandler)
	a.rg.POST("/auth/register", a.registerHandler)
	a.rg.GET("/usercred/:name", middleware.AuthMiddleware(), a.showUserHandler)
	a.rg.GET("/user/:page/:size", middleware.AuthMiddleware(), a.paggingUserHandler)
}

func NewAuthController(userUC usecase.UserUseCase, authUC usecase.AuthUseCase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{
		userUC: userUC,
		authUC: authUC,
		rg:     rg,
	}
}
