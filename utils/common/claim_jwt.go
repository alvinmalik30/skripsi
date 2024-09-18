package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetRole(c *gin.Context) (string, error) {
	claims, ok := c.Get("claims")

	if ok {
		claimsMap, ok := claims.(jwt.MapClaims) // Assert ke tipe yang sesuai
		if !ok {
			return "", fmt.Errorf("claim is invalid")
		}

		role, ok := claimsMap["role"].(string)
		if !ok {
			return "", fmt.Errorf("data 'role' is invalid")
		}

		return role, nil
	} else {
		return "", fmt.Errorf("unautorized")
	}
}

func GetId(c *gin.Context) (string, error) {
	claims, ok := c.Get("claims")

	if ok {
		claimsMap, ok := claims.(jwt.MapClaims) // Assert ke tipe yang sesuai
		if !ok {
			return "", fmt.Errorf("id is invalid")
		}

		role, ok := claimsMap["id"].(string)
		if !ok {
			return "", fmt.Errorf("data 'id' is invalid")
		}

		return role, nil
	} else {
		return "", fmt.Errorf("unautorized")
	}
}

func GetName(c *gin.Context) (string, error) {
	claims, ok := c.Get("claims")

	if ok {
		claimsMap, ok := claims.(jwt.MapClaims) // Assert ke tipe yang sesuai
		if !ok {
			return "", fmt.Errorf("id is invalid")
		}

		role, ok := claimsMap["username"].(string)
		if !ok {
			return "", fmt.Errorf("data 'name' is invalid")
		}

		return role, nil
	} else {
		return "", fmt.Errorf("unautorized")
	}
}
