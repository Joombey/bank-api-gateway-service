package http

import (
	"net/http"
	"strconv"

	"farukh.go/api-gateway/handlers"
	"farukh.go/api-gateway/models"
	tr "farukh.go/api-gateway/tracing"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	ms "github.com/mitchellh/mapstructure"
	otlgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Init() {
	router := gin.Default()

	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	router.Use(p.Instrument())
	
	_, err := tr.InitTracer("jaeger:4318", "api gateway service")
	if err != nil {
		panic(err.Error())
	}
	
	router.Use(otlgin.Middleware("my-server"))
	router.GET("/credentials/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, "Invalid id, not a number")
		}
		response := handlers.GetProfileHandler(id, ctx)
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
		if err != nil {
			panic(err.Error())
		}

		response := handlers.Transfer(
			transferRequest.Username,
			data.Token,
			transferRequest.From,
			transferRequest.To,
			transferRequest.Value,
			ctx,
		)
		ctx.IndentedJSON(response.Err.ErrorCode, response)
	})

	router.POST("/auth", func(ctx *gin.Context) {
		var request models.RegisterRequest

		ctx.BindJSON(&request)
		token, err := handlers.Login(request.Username, request.Password, ctx)

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
		response := handlers.LoadMoney(request, ctx)
		ctx.IndentedJSON(response.Err.ErrorCode, response)
	})

	router.POST("/create", func(ctx *gin.Context) {
		var request models.RegisterRequest
		ctx.BindJSON(&request)
		response := handlers.CreateUser(request.Username, request.Password, ctx)
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
		if err != nil {
			panic(err.Error())
		}

		println(data.Token.AccessToken)
		response := handlers.BlockUser(id, blockRequest.Caller, blockRequest.Username, data.Token, ctx)
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
		if err != nil {
			panic(err.Error())
		}

		handlers.UpdateUser(blockRequest.Caller, blockRequest.Username, role, data.Token, ctx)
	})

	router.Run("0.0.0.0:8080")
}
