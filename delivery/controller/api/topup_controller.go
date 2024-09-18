package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"polen/delivery/middleware"
	"polen/model/dto"
	"polen/usecase"
	"polen/utils/common"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TopUpController struct {
	topupUC usecase.TopUpUseCase
	bioUc   usecase.BiodataUserUseCase
	rg      *gin.RouterGroup
}

// @Summary get by id user
// @Description get data top up by user id
// @Tags topup
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param id path int true "id"
// @Success 200 {object} dto.ResponseData
// @Router /topup/user/{id} [GET]
func (t *TopUpController) getByIdUserId(c *gin.Context) {
	id := c.Param("id")
	// ambil role
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
	// jika bukan pemodal maka tidak boleh
	if role != "admin" {
		c.JSON(403, gin.H{
			"message": "you are not allowed",
		})
		return
	}
	data, err := t.topupUC.FindByIdUser(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success getting data",
		"data":    data,
	})
}

// @Summary get by id data
// @Description get data top up by id data
// @Tags topup
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param id path int true "id"
// @Success 200 {object} dto.ResponseData
// @Router /topup/{id} [GET]
func (t *TopUpController) getById(c *gin.Context) {
	id := c.Param("id")
	data, err := t.topupUC.FindById(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := dto.ResponseData{
		Message: "success getting data",
		Data:    data,
	}
	c.JSON(200, response)
}

// @Summary confirm
// @Description confirm data topup
// @Tags topup
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.TopUpupdate true "Data top up"
// @Success 200 {object} dto.ResponseMessage
// @Router /topup/confirm [PUT]
func (t *TopUpController) ConfirmUpload(c *gin.Context) {
	// ambil role
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

	// cek apakah user adalah admin
	if role != "admin" {
		c.JSON(403, gin.H{
			"message": "you are not allowed",
		})
		return
	}

	// ambil data request
	var data dto.TopUpupdate
	var topup dto.TopUpUser
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	topup.Accepted = data.Accepted
	topup.Id = data.Id
	topup.Status = data.Status
	code, err := t.topupUC.ConfimUploadFile(topup)
	if err != nil {
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}
	endresult := dto.ResponseMessage{
		Message: "success creating data",
	}
	c.JSON(http.StatusCreated, endresult)
}

// @Summary get data updated
// @Description get updated data topup
// @Tags topup
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Success 200 {object} dto.ResponseData
// @Router /topup/confirm [GET]
func (t *TopUpController) UploadedFile(c *gin.Context) {
	// ambil role
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

	// cek apakah user adalah admin
	if role != "admin" {
		c.JSON(403, gin.H{
			"message": "you are not allowed",
		})
		return
	}

	data, err := t.topupUC.FindUploadedFile()
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
			// "message": "internal server error",
		})
	}

	response := dto.ResponseData{
		Message: "success getting data",
		Data:    data,
	}
	c.JSON(200, response)
}

// @Summary new
// @Description create new data topup
// @Tags topup
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param request body dto.TopUpReq true "Data topup"
// @Success 200 {object} dto.ResponseData
// @Router /topup [POST]
func (t *TopUpController) createHandler(c *gin.Context) {
	// ambil id
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

	// ambil role
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

	// cek apakah user adalah pemodal
	if role != "pemodal" {
		c.JSON(403, gin.H{
			"message": "you are not allowed",
		})
		return
	}

	// cek apakah user sudah terkonfirmasi boleh melakukan top up
	bio, err := t.bioUc.FindByUserCredential(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	fmt.Println(bio.IsAglible, bio.Information)
	if !bio.IsAglible {
		c.JSON(403, gin.H{
			"message": []string{
				"you are not allowed to di this transaction",
				bio.Information,
			},
		})
		return
	}

	var data dto.TopUpReq
	var topup dto.TopUpUser
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	topup.Id = uuid.NewString()
	topup.TopUpAmount = data.TopUpAmount
	topup.UserCredential.Id = ucid
	result, err := t.topupUC.CreateNew(topup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	endresult := dto.ResponseData{
		Message: "success creating data",
		Data:    result,
	}
	c.JSON(http.StatusCreated, endresult)
}

// @Summary get by user login
// @Description get data top up by user login
// @Tags topup
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Success 200 {object} dto.ResponseData
// @Router /topup/user [GET]
func (t *TopUpController) getByIdUserLoginHandler(c *gin.Context) {
	// ambil role
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
	// jika bukan pemodal maka tidak boleh
	if role != "pemodal" {
		c.JSON(403, gin.H{
			"message": "you are not allowed",
		})
		return
	}
	// user id
	ucid, err := common.GetId(c)
	fmt.Println(ucid)
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
	// get data
	data, err := t.topupUC.FindByIdUser(ucid)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
			// "message": "internal server error",
		})
		return
	}
	response := dto.ResponseData{
		Message: "success getting data",
		Data:    data,
	}
	c.JSON(200, response)
}

// @Summary upload
// @Description upload recipe
// @Tags topup
// @Accept mpfd
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param file formData file true "File to upload"
// @Param id formData string true "ID parameter"
// @Success 200 {object} dto.ResponseMessage
// @Router /topup/upload [POST]
func (t *TopUpController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Membuat nama file unik dengan menambahkan timestamp di belakangnya
	fileName := filepath.Base(file.Filename)
	fileExt := filepath.Ext(fileName)
	timestamp := time.Now().Format("20060102150405") // Format timestamp yang diinginkan (YYYYMMDDHHmmss)
	uniqueFileName := fileName[:len(fileName)-len(fileExt)] + "_" + timestamp + fileExt

	// Menggunakan path file yang aman
	filePath := "uploads/" + uniqueFileName

	id := c.PostForm("id")

	// get user credential
	iduc, err := common.GetId(c)
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

	// payload
	payload := dto.TopUpUser{
		Id:             id,
		File:           filePath,
		UserCredential: dto.GetAuthResponse{Id: iduc},
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	code, err := t.topupUC.UploadFile(payload)
	if err != nil {
		_ = os.Remove(payload.File)
		c.JSON(code, gin.H{
			"message": err.Error(),
		})
		return
	}

	endresult := dto.ResponseMessage{
		Message: "success uploaded data",
	}
	c.JSON(http.StatusCreated, endresult)
}

// @Summary get all
// @Description get all data saldo
// @Tags topup
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token" default(Bearer <token>)
// @Param page path int true "page of pagination"
// @Param size path int true "size of pagination"
// @Success 200 {object} dto.ResponsePaging
// @Router /topup/list/{page}/{size} [GET]
func (t *TopUpController) paggingHandler(c *gin.Context) {
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

	model, pagereturn, err := t.topupUC.Pagging(payload)
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

func (t *TopUpController) Route() {
	t.rg.POST("/topup", middleware.AuthMiddleware(), t.createHandler)
	t.rg.POST("/topup/upload", middleware.AuthMiddleware(), t.UploadFile)
	t.rg.PUT("/topup/confirm", middleware.AuthMiddleware(), t.ConfirmUpload)
	t.rg.GET("/topup/uploaded", middleware.AuthMiddleware(), t.UploadedFile)
	t.rg.GET("/topup/user/:id", middleware.AuthMiddleware(), t.getByIdUserId)
	t.rg.GET("/topup/:id", middleware.AuthMiddleware(), t.getById)
	t.rg.GET("/topup/user", middleware.AuthMiddleware(), t.getByIdUserLoginHandler)
	t.rg.GET("/topup/list/:page/:size", middleware.AuthMiddleware(), t.paggingHandler)
	// t.rg.PUT("/topup/update", t.updateHandler)

}

func NewTopUpController(topupUC usecase.TopUpUseCase, bioUc usecase.BiodataUserUseCase, rg *gin.RouterGroup) *TopUpController {
	return &TopUpController{
		topupUC: topupUC,
		rg:      rg,
		bioUc:   bioUc,
	}
}
