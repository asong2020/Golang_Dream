package middleware

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/model/request"
	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/util/response"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, "非法访问", c)
			c.Abort()
			return
		}
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, err.Error(), c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.AsongServer.Jwt.Signkey),
	}
}

func (j *JWT) GenerateToken(claims request.UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(t string) (*request.UserClaims, error) {
	token, err := jwt.ParseWithClaims(t, &request.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("That's not even a token")
			} else if v.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.New("Token is expired")
			} else if v.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("Token not active yet")
			} else {
				return nil, errors.New("Couldn't handle this token:")
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.UserClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.New("Couldn't handle this token:")

	} else {
		return nil, errors.New("Couldn't handle this token:")

	}
}

func (j *JWT) RefreshToken(t string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(t, &request.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*request.UserClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.GenerateToken(*claims)
	}
	return "", errors.New("Couldn't handle this token:")
}
