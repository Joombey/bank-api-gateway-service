package keycloak

import (
	"context"

	md "farukh.go/api-gateway/models"
	gocloak "github.com/Nerzal/gocloak/v13"
)

var client *gocloak.GoCloak
var clientID string
var idOfClient string
var roles []gocloak.Role

func Init() {
	client = gocloak.NewClient("http://localhost:8086")
	tryCreateClient()
	tryCreateRoles()
}

func LoginAdmin() *gocloak.JWT {
	ctx := context.Background()
	jwt, err := client.LoginAdmin(ctx, "admin", "admin", "master")
	if err != nil {
		panic(err.Error())
	}
	return jwt
}

func RegisterUser(username, password, role string) (string, error) {
	jwt := LoginAdmin()
	ctx := context.Background()

	user := gocloak.User{Username: &username, Enabled: gocloak.BoolP(true)}
	
	userID, err := client.CreateUser(ctx, jwt.AccessToken, "master", user)
	
	err = client.AddClientRolesToUser(ctx, jwt.AccessToken, "master", idOfClient, userID, getRolesByName(role))
	return userID, err
}

func getRolesByName(name string) ([]gocloak.Role) {
	for id, role := range roles {
		if *role.Name == name {
			return roles[id:id]
		}
	}
	return []gocloak.Role{}
}

func CreateUser(username string) string {
	jwt := LoginAdmin()
	user := gocloak.User{
		Username: gocloak.StringP(username),
		Enabled:  gocloak.BoolP(true),
	}
	userId, err := client.CreateUser(context.Background(), jwt.AccessToken, "master", user)
	if err != nil {
		println(err.Error())
	} else {
		println(jwt)
	}
	return userId
}

func GetUser(userID string) *gocloak.User {
	jwt := LoginAdmin()
	user, err := client.GetUserByID(context.Background(), jwt.AccessToken, "master", userID)
	if err != nil {
		println(err.Error())
	} else {
		println(jwt)
	}
	return user
}

func Register(req md.RegisterRequest) string {
	jwt := LoginAdmin()
	userID := CreateUser(req.Username)
	if userID == "" {
		return "not OK"
	}
	err := client.SetPassword(context.Background(), jwt.AccessToken, userID, "master", req.Password, false)
	if err != nil {
		panic(err.Error())
	}

	return userID
}

func GetClients() {
	jwt := LoginAdmin()
	clients, err := client.GetClients(context.Background(), jwt.AccessToken, "master", gocloak.GetClientsParams{})
	if err != nil {
		panic(err.Error())
	}
	for _, client := range clients {
		println(*client.Name, *client.ClientID, *client.ID)
	}
}

func Auth(req md.RegisterRequest) string {
	// jwt := LoginAdmin()
	_, err := client.LoginClient(context.Background(), "profile-2", "gcVzad8Ua3qA8O2GqKRo3GpgsvGAf5gg", "master")
	if err != nil {
		panic(err.Error())
	}
	// client.GetClients()

	if err != nil {
		panic(err.Error())
	}
	userJWT, err := client.Login(context.Background(), "profile-2", "gcVzad8Ua3qA8O2GqKRo3GpgsvGAf5gg", "master", req.Username, req.Password)
	if err != nil {
		panic(err.Error())
	}
	println(
		userJWT.AccessToken,
		userJWT.ExpiresIn,
		userJWT.RefreshToken,
	)
	return userJWT.AccessToken
}

func Decode(token string) {
	println(token)
	tk, _, err := client.DecodeAccessToken(context.Background(), token, "master")
	if err != nil {
		panic(err.Error())
	}
	println(tk.Valid)
}

func RoleCheck(req md.RegisterRequest) {
	jwt := LoginAdmin()

	clients, err := client.GetClients(context.Background(), jwt.AccessToken, "master", gocloak.GetClientsParams{
		ClientID: gocloak.StringP("profile-2"),
	})

	if err != nil {
		panic(err.Error())
	}

	role := gocloak.Role{
		Name:       gocloak.StringP("profiler"),
		ClientRole: gocloak.BoolP(true),
	}

	roleId, err := client.CreateClientRole(context.Background(), jwt.AccessToken, "master", *clients[0].ID, role)
	if err != nil {
		panic(err.Error())
	}

	println(clients[0].String())

	rolemap := make(map[string][]string)
	stringMap := make([]string, 1)
	stringMap = append(stringMap, "user")

	rolemap["admin123"] = stringMap

	userId, err := client.CreateUser(context.Background(), jwt.AccessToken, "master", gocloak.User{
		Username:    &req.Username,
		ClientRoles: &rolemap,
		Enabled:     gocloak.BoolP(true),
	})

	if err != nil {
		println(err.Error())
	}

	client.SetPassword(context.Background(), jwt.AccessToken, userId, "master", req.Password, false)

	if err != nil {
		panic(err.Error())
	}

	println(roleId)

	role2, err := client.GetClientRoleByID(context.Background(), jwt.AccessToken, "master", roleId)

	if err != nil {
		panic(err)
	}

	println(role2.String())

	roleSlice := []gocloak.Role{*role2}
	err = client.AddClientRolesToUser(context.Background(), jwt.AccessToken, "master", "profile-2", userId, roleSlice)

	if err != nil {
		panic(err)
	}

	// userActual, err := client.GetUserByID(context.Background(), jwt.AccessToken, "master", userId)

	// if err != nil {
	// 	println(err.Error())
	// }
}

func GetUserByUsername(username string) string {
	jwt := LoginAdmin()
	result, err := client.GetUsers(
		context.Background(),
		jwt.AccessToken,
		"master",
		gocloak.GetUsersParams{
			Username: &username,
		},
	)

	if err != nil {
		panic(err.Error())
	}
	clients, err := client.GetClients(context.Background(), jwt.AccessToken, "master", gocloak.GetClientsParams{ClientID: gocloak.StringP("profile-2")})

	if err != nil {
		panic(err.Error())
	}
	AddRole(*result[0].ID, *clients[0].ID)

	return result[0].String()
}

func AddRole(userID string, clientId string) {
	jwt := LoginAdmin()
	role1, err := client.GetClientRole(
		context.Background(),
		jwt.AccessToken,
		"master",
		clientId,
		"profiler",
	)

	if err != nil {
		panic(err.Error())
	}

	err = client.AddClientRolesToUser(
		context.Background(),
		jwt.AccessToken,
		"master",
		clientId,
		userID,
		[]gocloak.Role{*role1},
	)

	if err != nil {
		panic(err.Error())
	}
}
