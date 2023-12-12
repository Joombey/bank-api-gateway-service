package http

import (
	"net/http"
	"strconv"

	"farukh.go/api-gateway/handlers"
	kk "farukh.go/api-gateway/keycloak"
	"farukh.go/api-gateway/models"
	"github.com/gin-gonic/gin"
	cts "farukh.go/api-gateway/constants"
)



func Init() {
	router := gin.Default()
	router.GET("/credentials/:id", func (ctx *gin.Context)  {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, "Invalid id, not a number")
		}
		handlers.GetProfileHandler(id)
	})

	router.POST("/transfer", func(ctx *gin.Context) {
		data := models.RegisterRequest {}
		err := ctx.BindJSON(&data)
		if err != nil { 
			ctx.IndentedJSON(http.StatusBadRequest, "Invalid body json format")
		}
		err, isAuth := handlers.Transfer(token models.Token, from to value)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		}
		if isAuth {
			ctx.IndentedJSON(http.StatusForbidden, "invalid token, need to reauth")
		}
	})

	router.POST("/auth", func(ctx *gin.Context) {


	})

	router.Run()
}