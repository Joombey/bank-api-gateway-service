package http

import (
	"net/http"
	"strconv"

	"farukh.go/api-gateway/handlers"
	"farukh.go/api-gateway/models"
	"github.com/gin-gonic/gin"
	ms "github.com/mitchellh/mapstructure"
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

		var transferRequest models.TransferDTO
		err = ms.Decode(data.Body, &transferRequest)
		if err!= nil{
			panic(err.Error())
		}

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

	router.POST("/create", func(ctx *gin.Context) {
		var request models.RegisterRequest
		ctx.BindJSON(&request)
		
		response := handlers.CreateUser(request.Username, request.Password)
		ctx.IndentedJSON(response.Err.ErrorCode, response)
	})

	router.POST("/block/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, "invalid id not a number")
		}

		data := models.BaseHTTPModel{}
		err = ctx.BindJSON(&data)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, "Invalid body json format")
		}

		var blockRequest models.DeleteUserRequest
		err = ms.Decode(data.Body, &blockRequest)
		if err!= nil{
			panic(err.Error())
		}
		
		println(data.Token.AccessToken)
		response := handlers.BlockUser(id, blockRequest.Caller, blockRequest.Username, data.Token)
		ctx.IndentedJSON(response.Err.ErrorCode, response)
	})

	router.POST("/update/:role", func(ctx *gin.Context) {
		role := ctx.Param("role")
		data := models.BaseHTTPModel{}
		err := ctx.BindJSON(&data)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, "Invalid body json format")
		}

		var blockRequest models.DeleteUserRequest
		err = ms.Decode(data.Body, &blockRequest)
		if err!= nil{
			panic(err.Error())
		}
		
		handlers.UpdateUser(blockRequest.Caller, blockRequest.Username, role, data.Token)
	})
	
	router.Run("0.0.0.0:8080")
}
