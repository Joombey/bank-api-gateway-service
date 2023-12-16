package keycloak

import (
	"context"
	"errors"
	"os"

	md "farukh.go/api-gateway/models"
	gocloak "github.com/Nerzal/gocloak/v13"
)

func Init() {
	client = gocloak.NewClient("http://keycloak:8080")
	tryCreateClient()
	tryCreateRoles()
	tryCreateFirstAdmin()
	obtainRoles()
	getClients()
}

func UpdateUser(target, role string) error {
	token := LoginAdmin().AccessToken

	users, err := client.GetUsers(
		context.Background(),
		token,
		realm,
		gocloak.GetUsersParams{ Username: gocloak.StringP(target) },
	)
	if err != nil {
		return err
	}
	
	var userID string = ""
	for _, user := range users {
		if *user.Username == target {
			userID = *user.ID
			break
		}
	}
	if userID == "" {
		return errors.New("no such user")
	}

	return setRoleForNewUser(userID, role, token)
}

func CheckRole(username, role string) (bool, error) {
	token := LoginAdmin().AccessToken
	users, err := client.GetUsersByClientRoleName(
		context.Background(),
			token,
			realm,
			idOfClient,
			role,
			gocloak.GetUsersByRoleParams{},
	)
	if err != nil {
		return false, err
	}

	hasRole := false
	println(len(users))
	for _, v := range users {
		println(*v.Username)
		if *v.Username == username{
			hasRole = true
			break
		}
	}
	return hasRole, nil
}

func DeleteUser(username string) error {
	token := LoginAdmin().AccessToken
	userRepr, err := client.GetUsers(context.Background(), token, realm, gocloak.GetUsersParams{Username: &username})
	if err != nil || len(userRepr) == 0 {
		return nil
	}

	return client.DeleteUser(
		context.Background(),
		token,
		realm,
		*userRepr[0].ID,
	)
}

func CheckToken(token md.Token) (newToken *md.Token, err error) {
	spectResult, err := client.RetrospectToken(
		context.Background(),
		token.AccessToken,
		ClientID,
		secret,
		realm,
	)

	if err != nil {
		return nil, err
	} else if !*spectResult.Active {
		return refreshToken(token)
	} else {
		return &token, nil
	}
}

func LoginAdmin() *gocloak.JWT {
	ctx := context.Background()
	jwt, err := client.LoginAdmin(ctx, "admin", "admin", realm)
	if err != nil {
		panic(err.Error())
	}
	return jwt
}

func RegisterUser(username, password, role string) (userID string, err error) {
	token := LoginAdmin().AccessToken

	userID, err = createUserWithPassword(username, password, token)
	if err != nil {
		return "", err
	}

	err = setRoleForNewUser(userID, role, token)
	return userID, err
}

func LoginUser(username, password string) (md.Token, error) {
	jwt, err := client.Login(
		context.Background(),
		ClientID,
		secret,
		realm,
		username,
		password,
	)
	if err != nil {
		panic(err.Error())
	}
	return md.Token{AccessToken: jwt.AccessToken, RefreshToken: jwt.RefreshToken}, err
}

func tryCreateFirstAdmin() {
	RegisterUser("manager", "manager", RoleAdmin)
}

func refreshToken(token md.Token) (*md.Token, error) {
	jwt, err := client.RefreshToken(context.Background(), token.RefreshToken, ClientID, secret, realm)
	if err != nil {
		return nil, err
	}
	token = md.Token{AccessToken: jwt.AccessToken, RefreshToken: jwt.RefreshToken}
	return &token, err
}

func obtainRoles() {
	jwt := LoginAdmin()
	roles, err := client.GetClientRoles(
		context.Background(),
		jwt.AccessToken,
		realm,
		idOfClient,
		gocloak.GetRoleParams{},
	)

	if err != nil {
		panic(err)
	}
	for _, v := range roles {
		baseRoles[*v.Name] = *v
	}
}

func tryCreateClient() {
	if idOfClient != "" {
		return
	}

	jwt := LoginAdmin()
	id, err := client.CreateClient(
		context.Background(),
		jwt.AccessToken,
		realm,
		gocloak.Client{
			ClientID:                     gocloak.StringP(ClientID),
			Enabled:                      gocloak.BoolP(true),
			Name:                         gocloak.StringP(ClientID),
			PublicClient:                 gocloak.BoolP(false),
			// AuthorizationServicesEnabled: gocloak.BoolP(true),
			DirectAccessGrantsEnabled:    gocloak.BoolP(true),
		},
	)
	if err != nil {
		idOfClient = getClient()
	} else {
		idOfClient = id
	}

	clientRepr, err := client.GetClientSecret(context.Background(), jwt.AccessToken, realm, idOfClient)
	if err != nil {
		panic(err.Error())
	}
	secret = gocloak.PString(clientRepr.Value)
	os.Setenv("ID_OF_CLIENT", idOfClient)
	os.Setenv("SECRET", secret)
}

func getClient() string {
	token := LoginAdmin().AccessToken
	clients, err := client.GetClients(
		context.Background(),
		token,
		realm,
		gocloak.GetClientsParams{ClientID: gocloak.StringP(ClientID)},
	)

	if err != nil {
		panic(err.Error())
	}
	return gocloak.PString(clients[0].ID)
}

func tryCreateRoles() {
	if os.Getenv("ROLES") != "" {
		return
	}

	jwt := LoginAdmin()
	for _, roleName := range Roles {
		role := gocloak.Role{Name: &roleName, ClientRole: gocloak.BoolP(true)}
		client.CreateClientRole(context.Background(), jwt.AccessToken, realm, idOfClient, role)
	}
	os.Setenv("ROLES", "SET")
}

func setRoleForNewUser(userID, roleName, token string) (err error) {
	var role gocloak.Role

	foundRole, err := client.GetClientRole(context.Background(), token, realm, idOfClient, roleName)
	if err != nil {
		role = baseRoles[RoleUser]
	} else {
		role = *foundRole
	}

	return client.AddClientRolesToUser(
		context.Background(),
		token,
		realm,
		idOfClient,
		userID,
		[]gocloak.Role{role},
	)
}

func createUserWithPassword(username, password, token string) (userID string, err error) {
	jwt := LoginAdmin()
	ctx := context.Background()

	user := gocloak.User{Username: &username, Enabled: gocloak.BoolP(true)}

	userID, err = client.CreateUser(ctx, jwt.AccessToken, realm, user)
	if err != nil {
		return "", err
	}

	client.SetPassword(ctx, jwt.AccessToken, userID, realm, password, false)
	return userID, nil
}

func getClients() {
	token := LoginAdmin().AccessToken
	clients, _ := client.GetClients(context.Background(), token, realm, gocloak.GetClientsParams{ClientID: gocloak.StringP("profile-2")})
	for _, v := range clients {
		println(
			gocloak.PString(v.ClientAuthenticatorType),
			gocloak.PBool(v.AuthorizationServicesEnabled),
		)
	}
}

const (
	RoleCardOwner string = "card-owner"
	RoleUser      string = "user"
	RoleAdmin     string = "admin"
	ClientID      string = "profile-app"
)

var Roles = []string{RoleAdmin, RoleUser, RoleCardOwner}

var (
	client     *gocloak.GoCloak
	idOfClient string                  = os.Getenv("ID_OF_CLIENT")
	secret     string                  = os.Getenv("SECRET")
	baseRoles  map[string]gocloak.Role = make(map[string]gocloak.Role, 0)
)

const realm string = "master"
