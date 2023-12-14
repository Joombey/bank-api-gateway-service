package http

import (
	"net/http"
	"strconv"

	"farukh.go/api-gateway/handlers"
	"farukh.go/api-gateway/models"
	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	
	router.GET("/credentials/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, "Invalid id, not a number")
		}
		response := handlers.GetProfileHandler(id)
		ctx.IndentedJSON(response.Err.ErrorCode, response)
	})

	router.POST("/transfer", func(ctx *gin.Context) {
		data := models.BaseHTTPModel{}
		err := ctx.BindJSON(&data)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, "Invalid body json format")
		}

		transferRequest := data.Body.(models.TransferDTO)
		response := handlers.Transfer(
			transferRequest.Username,
			data.Token,
			transferRequest.From,
			transferRequest.To,
			transferRequest.Value,
		)
		ctx.IndentedJSON(response.Err.ErrorCode, response)
	})

	router.POST("/auth", func(ctx *gin.Context) {
		var request models.RegisterRequest
		
		ctx.BindJSON(&request)
		token, err := handlers.Login(request.Username, request.Password)
		
		var status int
		if err != nil {
			status = http.StatusForbidden
		} else {
			status = http.StatusOK
		}

		ctx.IndentedJSON(status, token)
	})

	router.POST("/load-money", func(ctx *gin.Context) {
		var request models.InsertRequest
		ctx.BindJSON(&request)
		response := handlers.LoadMoney(request)
		ctx.IndentedJSON(response.Err.ErrorCode, response)
	})

	router.GET("/create/:name", func(ctx *gin.Context) {
		
	})
	router.Run()
}