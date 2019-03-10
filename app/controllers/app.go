package controllers

import (
	"fmt"
	"goost/app/models/randstring"
	"strconv"
	"strings"
	"time"

	"github.com/revel/revel"
)

// secrets contains the stored secrets
var secrets []Secret

// App contains Revel App Controller
type App struct {
	*revel.Controller
}

// Secret represents the actual secret to store, including meta data
type Secret struct {
	id            string
	secret        string
	views         int
	encryptionKey string
	ttl           int
}

// StoreResponse represents a Store JSON response
type StoreResponse struct {
	Success   bool   `json:"success"`
	SecretKey string `json:"secretKey"`
}

// RetrieveResponse represents a Retrive JSON response
type RetrieveResponse struct {
	Success bool   `json:"success"`
	Secret  string `json:"secret"`
}

// ErrorResponse represents an Error JSON response
type ErrorResponse struct {
	ErrorMsg string `json:"errorMsg"`
}

// UnsetSecret removes a Secret from the application
func UnsetSecret(i int) {
	secrets = append(secrets[:i], secrets[i+1:]...)
}

// Store handles POST to /store
func (c App) Store() revel.Result {
	views, _ := strconv.Atoi(c.Params.Form.Get("views"))
	ttl, _ := strconv.Atoi(c.Params.Form.Get("ttl"))

	if ttl <= 0 {
		return c.RenderJSON(ErrorResponse{
			ErrorMsg: "You need to enter a valid period longer than 0 days and 0 hours",
		})
	}

	if views <= 0 {
		return c.RenderJSON(ErrorResponse{
			ErrorMsg: "You need to enter number of views greater than 0",
		})
	}

	secret := Secret{
		id:            randstring.String(13),
		secret:        c.Params.Form.Get("password"),
		views:         views,
		encryptionKey: randstring.String(32),
		ttl:           int(time.Now().Unix()) + 10,
	}

	secrets = append(secrets, secret)

	StoreResponse := StoreResponse{
		Success:   true,
		SecretKey: secret.id + ";" + secret.encryptionKey,
	}
	return c.RenderJSON(StoreResponse)
}

// Retrieve handles POST to /retrieve
func (c App) Retrieve() revel.Result {
	secretKeys := strings.Split(c.Params.Form.Get("secretKey"), ";")
	if len(secretKeys) != 2 {
		return c.RenderJSON(ErrorResponse{
			ErrorMsg: "Invalid Password URL",
		})
	}

	secretID := secretKeys[0]
	encryptionKey := secretKeys[1]

	currentTimestamp := int(time.Now().Unix())

	for i := 0; i < len(secrets); i++ {
		secret := secrets[i]

		ttlLeft := secret.ttl - currentTimestamp
		if ttlLeft < 0 {
			UnsetSecret(i)
			continue
		}

		if secret.id == secretID {
			if secret.encryptionKey == encryptionKey {
				secret.views--

				if secret.views < 1 {
					UnsetSecret(i)
				} else {
					secrets[i] = secret
				}

				return c.RenderJSON(RetrieveResponse{
					Success: true,
					Secret:  secret.secret,
				})
			}

			fmt.Println("Invalid encryption key entered")
		}
	}

	return c.RenderJSON(ErrorResponse{
		ErrorMsg: "No password found for entered keys",
	})
}
