package input

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	auth "k8s.io/api/authentication/v1"

	"github.com/robbiemcmichael/auth-mux/internal/types"
	"gopkg.in/square/go-jose.v2/jwt"
)

type KubernetesTokenReview struct {
	// Token will be rejected if the audience does not match
	Audience []string `yaml:"audience"`
	// A map containing issuers and their validation configuration
	Issuers map[string]*Issuer `yaml:"issuers"`
	// Fields used to extract claims from the JWT
	Claims JWTClaims `yaml:"claims"`
}

type Issuer struct {
	// Path to the public key used to verify the JWT
	PublicKeyFile string `yaml:"publicKey"`
	// Parsed public key (*rsa.PublicKey, *dsa.PublicKey or *ecdsa.PublicKey)
	PublicKey interface{}
	// If provided, assert that the prefix is included in the user and groups claims
	Prefix string `yaml:"prefix"`
}

type JWTClaims struct {
	UID    string `yaml:"uid"`
	User   string `yaml:"user"`
	Groups string `yaml:"groups"`
	Extra  string `yaml:"extra"`
}

func (i *KubernetesTokenReview) Handler(r *http.Request) (types.Result, error) {
	decoder := json.NewDecoder(r.Body)

	var tokenReview auth.TokenReview
	if err := decoder.Decode(&tokenReview); err != nil {
		return types.Result{}, fmt.Errorf("decode JSON: %v", err)
	}

	result, err := i.validateToken(tokenReview.Spec.Token)
	if err != nil {
		return types.Result{}, fmt.Errorf("validate token: %v", err)
	}

	return result, nil
}

func (i *KubernetesTokenReview) validateToken(tokenString string) (types.Result, error) {
	token, err := jwt.ParseSigned(tokenString)
	if err != nil {
		return types.Result{}, fmt.Errorf("parse token: %v", err)
	}

	issuerString, err := getIssuer(*token)
	if err != nil {
		return types.Result{}, fmt.Errorf("get issuer: %v", err)
	}

	issuer := i.Issuers[issuerString]
	if issuer == nil {
		invalid := types.Result{
			Valid: false,
			Error: fmt.Sprintf("Invalid token: unknown issuer %q", issuerString),
		}
		return invalid, nil
	}

	if issuer.PublicKey == nil {
		if err := issuer.parsePublicKey(); err != nil {
			return types.Result{}, fmt.Errorf("parse public key for issuer %q: %v", issuerString, err)
		}
	}

	var publicClaims jwt.Claims
	if err := token.Claims(issuer.PublicKey, &publicClaims); err != nil {
		return types.Result{}, fmt.Errorf("parse public claims: %v", err)
	}

	expected := jwt.Expected{
		Audience: i.Audience,
		Time:     time.Now(),
	}
	if err := publicClaims.Validate(expected); err != nil {
		invalid := types.Result{
			Valid: false,
			Error: fmt.Sprintf("Invalid token: %v", err),
		}
		return invalid, nil
	}

	claims, err := getClaims(*token, i.Claims)
	if err != nil {
		return types.Result{}, err
	}

	result := types.Result{
		Valid:  true,
		Claims: claims,
	}

	return result, nil
}

func (issuer *Issuer) parsePublicKey() error {
	contents, err := ioutil.ReadFile(issuer.PublicKeyFile)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(contents)
	if block == nil {
		return fmt.Errorf("failed to read PEM block")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	issuer.PublicKey = key
	return nil
}

func getIssuer(token jwt.JSONWebToken) (string, error) {
	var publicClaims jwt.Claims
	if err := token.UnsafeClaimsWithoutVerification(&publicClaims); err != nil {
		return "", err
	}

	return publicClaims.Issuer, nil
}

func getClaims(token jwt.JSONWebToken, fields JWTClaims) (types.Claims, error) {
	var claims interface{}
	if err := token.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return types.Claims{}, err
	}

	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		return types.Claims{}, fmt.Errorf("failed to cast JWT claims to map[string]interface{}")
	}

	uid, ok := claimsMap[fields.UID].(string)
	if !ok {
		return types.Claims{}, fmt.Errorf("failed to cast %q claim to string", fields.UID)
	}

	user, ok := claimsMap[fields.User].(string)
	if !ok {
		return types.Claims{}, fmt.Errorf("failed to cast %q claim to string", fields.User)
	}

	interfaceArray, ok := claimsMap[fields.Groups].([]interface{})
	if !ok {
		return types.Claims{}, fmt.Errorf("failed to cast %q claim to []interface{}", fields.Groups)
	}

	groups := make([]string, len(interfaceArray))
	for i, v := range interfaceArray {
		group, ok := v.(string)
		if !ok {
			return types.Claims{}, fmt.Errorf("failed to cast group to string: %v", v)
		}

		groups[i] = group
	}

	c := types.Claims{
		UID:    uid,
		User:   user,
		Groups: groups,
		Extra:  claimsMap[fields.Extra],
	}

	return c, nil
}
