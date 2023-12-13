package di

import "farukh.go/api-gateway/services"

type BaseContainer struct {
	Profile  services.ProfileService
}

func GetContiner() BaseContainer {
	return BaseContainer{}
}