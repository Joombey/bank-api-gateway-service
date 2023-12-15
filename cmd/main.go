package main

import (
	"farukh.go/api-gateway/http"
	"farukh.go/api-gateway/keycloak"
)

func main() {
	keycloak.Init()
	http.Init()
}
