package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Mshivam2409/AnC-Courses/models"
	"github.com/gofiber/fiber/v2"
	kratos "github.com/ory/kratos-client-go/client"
	"github.com/ory/kratos-client-go/client/admin"
	kratosmodels "github.com/ory/kratos-client-go/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Ory struct {
	*kratos.OryKratos
}

var OryClient = &Ory{}

func InitalizeOryServices() {
	OryClient.OryKratos = kratos.NewHTTPClientWithConfig(nil, &kratos.TransportConfig{
		Host:     viper.GetString("kratos.admin_url"),
		BasePath: "/",
		Schemes:  []string{"http"},
	})
}

// CLEARANCE : Clearance levels required for access
var CLEARANCE = map[string]int{
	"REGISTER":         1,
	"READ_COURSE_DATA": 1,
	"ADD_COURSE":       2,
	"ADD_REVIEW":       1,
	"GOOGLE_DRIVE":     2,
	"TOGGLE_REVIEW":    2,
	"ELEVATE_USER":     3,
	"DEMOTE_USER":      3,
}

// IsAuthorized Checks access for a given resource
func (ory *Ory) IsAuthorized(tokenString string, resource string) (bool, error) {
	c, err := ory.GetClearanceLevel(tokenString)
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	if c >= CLEARANCE[resource] {
		return true, err
	}

	return false, err
}

func (ory *Ory) CanRegister(username string) (bool, error) {
	u := &models.MGMUser{}
	name := strings.Replace(username, " ", "", -1)
	filter := bson.M{"username": name}
	err := MongoClient.Users.Collection("ug").FindOne(context.TODO(), filter).Decode(u)
	// TODO Implement banning
	if err != nil {
		log.Printf("Unable to check access for %s: %v", name, err)
	}
	return true, err
}

func (ory *Ory) SetKratosID(kid string, username string) error {
	u := &models.MGMUser{}
	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "kid", Value: kid}}}}
	fmt.Println(kid, 11)
	err := MongoClient.Users.Collection("ug").FindOneAndUpdate(context.TODO(), filter, update).Decode(u)
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	return err
}

// CreateUser sdfsd
func (ory *Ory) CreateUser(username string) error {
	ok, err := ory.CanRegister(username)
	if err != nil {
		log.Printf("Unable to check regsitration : %v", err)
		return err
	}
	if !ok {
		log.Printf("Unable to check regsitration : %v", err)
		return err
	}
	// s := viper.GetString("kratos.schema")
	i, err := ory.OryKratos.Admin.CreateIdentity(&admin.CreateIdentityParams{Body: &kratosmodels.CreateIdentity{Traits: fiber.Map{"email": username + "@iitk.ac.in"}}, Context: context.TODO()})
	if err != nil {
		log.Printf("Unable to check regsitration : %v %v", err, err.Error())
		return err
	}
	id := i.GetPayload().ID
	fmt.Print(id)
	link, err := ory.OryKratos.Admin.CreateRecoveryLink(&admin.CreateRecoveryLinkParams{Body: &kratosmodels.CreateRecoveryLink{IdentityID: id}, Context: context.TODO()})
	if err != nil {
		log.Printf("Unable to check regsitration : %v", err)
		return err
	}
	url := *link.Payload.RecoveryLink
	message := "Dear User,\nPlease use the following link to recover your account:\n" + url
	SendMail(message, username+"@iitk.ac.in")
	ory.SetKratosID(string(id), username)
	return err
}

func (ory *Ory) DeleteUser(id string) error {
	_, err := ory.OryKratos.Admin.DeleteIdentity(&admin.DeleteIdentityParams{ID: id})
	return err
}

func (ory *Ory) GetClearanceLevel(tokenString string) (int, error) {
	var claims map[string]interface{}
	token, err := jwt.ParseSigned(tokenString)
	if err != nil {
		return 0, err
	}
	err = token.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		return 0, err
	}
	u := &models.MGMUser{}
	filter := bson.D{{Key: "username", Value: claims["username"]}}
	err = MongoClient.Users.Collection("ug").FindOne(context.TODO(), filter).Decode(u)
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	return u.Clearance, err
}
