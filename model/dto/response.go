package dto

type ResponseMessage struct {
	Message string `json:"message"`
}

type ResponseData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponsePaging struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Paging  Paging      `json:"paging"`
}
