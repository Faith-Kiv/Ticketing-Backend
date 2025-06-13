package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Faith-Kiv/Ticketing-Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

const (
	VALID_TOKEN_PARTS = 3
)

type CertResponse struct {
	KID     string   `json:"kid" binding:"required"`
	Kty     string   `json:"kty" binding:"required"`
	Alg     string   `json:"alg" binding:"required"`
	Use     string   `json:"use" binding:"required"`
	N       string   `json:"n" binding:"required"`
	E       string   `json:"e" binding:"required"`
	X5c     []string `json:"x5c" binding:"required"`
	X5t     string   `json:"x5t" binding:"required"`
	Xt_S256 string   `json:"x5t#S256" binding:"required"`
}

var CertStorage = make(map[string]string)

func ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if !strings.HasPrefix(ctx.FullPath(), "/api/v2") {
			ctx.Next()
			return
		}
		headerToken := ctx.GetHeader("Authorization")

		signedToken := strings.Split(headerToken, " ")

		if len(signedToken) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := VerifyAndParse(signedToken[1])
		if err != nil {
			logrus.Info(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		value, ok := claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		roles := value["realm_access"].(map[string]interface{})["roles"].([]interface{})
		s := make([]string, len(roles))
		for i, v := range roles {
			s[i] = fmt.Sprint(v)
		}
		groups, ok := value["Groups"].([]interface{})
		if ok {
			g := make([]string, len(groups))
			for k, m := range groups {
				g[k] = fmt.Sprint(m)
			}
			ctx.Set("groups", g)
		}

		ctx.Set("user_id", value["sub"])
		ctx.Set("email", value["email"])
		ctx.Set("roles", s)

		ctx.Next()

	}
}

func getVerificationCerts() {
	url := os.Getenv("AUTH_CERT_URL")
	keys := struct {
		Keys []CertResponse `json:"keys" binding:"required"`
	}{}

	response, err := utils.Request("", map[string][]string{}, url, "GET")
	if err != nil {
		panic("failed to get auth certs")
	}

	err = json.Unmarshal([]byte(response), &keys)
	if err != nil {
		logrus.Fatal(err)
		panic("failed to process cert response")
	}
	MemoryMapStoreCerts(keys.Keys)
}

func MemoryMapStoreCerts(certs []CertResponse) {
	for _, cert := range certs {
		CertStorage[cert.KID] = cert.X5c[0]
	}
}

func init() {
	getVerificationCerts()
}

func VerifyAndParse(token string) (interface{}, error) {

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		cert, ok := CertStorage[jwtToken.Header["kid"].(string)]
		if !ok {
			return nil, fmt.Errorf("key not found: %s", jwtToken.Header["kid"])
		}
		cert = fmt.Sprintf("-----BEGIN RSA PUBLIC KEY-----\n%s\n-----END RSA PUBLIC KEY-----", cert)
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, fmt.Errorf("key not found: %s", err.Error())
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
