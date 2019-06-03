package input

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "k8s.io/api/authentication/v1"

	"github.com/robbiemcmichael/auth-mux/internal/types"
        "gopkg.in/square/go-jose.v2/jwt"
)

type KubernetesTokenReview struct {
	PublicKey    string    `yaml:"publicKey"`
	Audience     string    `yaml:"audience"`
	Claims       JWTClaims `yaml:"claims"`
	Prefix       string    `yaml:"prefix"`
}

type JWTClaims struct {
	UIDClaim    string `yaml:"uid"`
	UserClaim   string `yaml:"user"`
	GroupsClaim string `yaml:"groups"`
	ExtraClaim  string `yaml:"extra"`
}

func (i *KubernetesTokenReview) Handler(r *http.Request) (types.Result, error) {
	decoder := json.NewDecoder(r.Body)

	var tokenReview auth.TokenReview
	if err := decoder.Decode(&tokenReview); err != nil{
		return types.Result{}, fmt.Errorf("decode JSON: %+v", err)
	}

	claims, err := i.getClaims(tokenReview.Spec.Token)
	if err != nil{
		return types.Result{}, fmt.Errorf("get token claims: %+v", err)
	}

	result := types.Result{
		Valid: true,
		Claims: claims,
	}

	return result, nil
}

func (i *KubernetesTokenReview) getClaims(tokenString string) (types.Claims, error) {
        token, err := jwt.ParseSigned(tokenString)
        if err != nil {
                return types.Claims{}, err
        }

        var publicClaims jwt.Claims
        if err := token.UnsafeClaimsWithoutVerification(&publicClaims); err != nil {
                return types.Claims{}, fmt.Errorf("failed to get public claims: %s", err)
        }

	claims := types.Claims{
		UID: publicClaims.Subject,
		User: publicClaims.Subject,
	}

        return claims, nil
}
