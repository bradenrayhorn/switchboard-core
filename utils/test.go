package utils

import (
	"github.com/Kamva/mgm/v3"
	"github.com/bradenrayhorn/switchboard-core/config"
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

const privateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgH8URhhPVVlVSgs/JImasa9gX4sCtyZJYVxrfPyNKpNQHb+14GNm
SnQxXRh/tYhQy/KjflqXZk97BcnmXwUlmkITRs+fJhBL3oy9C0BpKzhAPdneVRpJ
naCzVNoH7Latbjvc1NyHVKQrQq8HSsnw1RADnuZWlR278PDO0oVOW1AtAgMBAAEC
gYB72MhPXNGzFDnrKAh1yrssTeIPWgAgYhduuJrAjttVYhj8A0bB2KjrAEjYXW4P
gZ6hw9CafT6KazzC/a7RF0pzdoWRlLyF8pGmxhW4CWhTiI5YPLd3y4MxSq8jCSEP
6RraQDlKUriPfzuwgoR1Jh0Qs6lOT462Gz9VwSBvrMP3yQJBAPVDEjXJ7cjAz3Uw
E9iileqAGIF5u+GH247SjanCLfyid52EpuBQOz/+aWH1GMMQyGhmKshPZw32+9jj
ydaF3kMCQQCEpJgeW/CO7wwS4/3TPFr8iVmGB5vMoyh3ALlnfRsC2wg8s6xlxvTF
nATyEhuwV9UI4mOthAR6UE3kNz6Pw4jPAkEAxHv6H7to0racvN5KV/hQr+/1Mezr
do+XdoD52rDklIs31qqJ4hSEkwznMgHf144fb3vB9H2gKtDeDHDYpxLR7wJAO1HU
yfb6DSIw9x3JPTfHxRqz07IBZjItfZLwV6zmcI9+Do+X8OhaPSm6OHwKsAGHv3Jn
e4kH65+QRhjCvM6IlwJBAPDLfM+h6NH2iYu44bsY8DpcS5BcMmumsK1Hkzabyr8f
Vfqvv7RyNfUj8P8N2IEp2VM7MPPZXyHnmqJW/9BLLpo=
-----END RSA PRIVATE KEY-----
`

const publicKey = `
-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgH8URhhPVVlVSgs/JImasa9gX4sC
tyZJYVxrfPyNKpNQHb+14GNmSnQxXRh/tYhQy/KjflqXZk97BcnmXwUlmkITRs+f
JhBL3oy9C0BpKzhAPdneVRpJnaCzVNoH7Latbjvc1NyHVKQrQq8HSsnw1RADnuZW
lR278PDO0oVOW1AtAgMBAAE=
-----END PUBLIC KEY-----
`

func SetupTestRsaKeys() {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		log.Fatal("failed to generate test rsa key")
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		log.Fatal("failed to generate test rsa key")
	}
	config.RsaPrivate = privateKey
	config.RsaPublic = publicKey
}

func MakeTestUser(username string, passwordHash string) *models.User {
	return &models.User{
		DefaultModel: mgm.DefaultModel{
			IDField: mgm.IDField{ID: primitive.NewObjectID()},
		},
		Username: username,
		Password: passwordHash,
	}
}

func MakeTestGroup(name *string, users []primitive.ObjectID) models.Group {
	return models.Group{
		DefaultModel: mgm.DefaultModel{
			IDField: mgm.IDField{ID: primitive.NewObjectID()},
		},
		Name:    name,
		UserIds: users,
	}
}
