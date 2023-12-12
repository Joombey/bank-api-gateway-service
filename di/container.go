package di

import "farukh.go/api-gateway/services"

type BaseContainer struct {
	Profile  services.ProfileService
	Keycloak services.KeycloakService
}

func GetContiner() BaseContainer {
	return BaseContainer{}
}