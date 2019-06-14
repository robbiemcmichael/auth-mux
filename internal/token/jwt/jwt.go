package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	joseJWT "gopkg.in/square/go-jose.v2/jwt"

	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type Config struct {
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
	// Assert that the ID includes a prefix
	IDPrefix string `yaml:"idPrefix"`
	// Assert that the subject includes a prefix
	SubjectPrefix string `yaml:"subjectPrefix"`
	// Assert that the groups include a prefix
	GroupPrefix string `yaml:"groupPrefix"`
}

type JWTClaims struct {
	// JWT claim to map to the ID
	ID string `yaml:"id"`
	// JWT claim to map to the subject
	Subject string `yaml:"subject"`
	// JWT claim to map to the groups
	Groups string `yaml:"groups"`
	// JWT claim to map to extra
	Extra string `yaml:"extra"`
}

func (c *Config) Validate(tokenString string) (types.Validation, error) {
	token, err := joseJWT.ParseSigned(tokenString)
	if err != nil {
		invalid := types.Validation{
			Valid: false,
			Error: fmt.Sprintf("Failed to parse token: %v\n", err),
		}
		return invalid, nil
	}

	issuerString, err := getIssuer(*token)
	if err != nil {
		invalid := types.Validation{
			Valid: false,
			Error: fmt.Sprintf("Failed to extract token claims: %v\n", err),
		}
		return invalid, nil
	}

	issuer := c.Issuers[issuerString]
	if issuer == nil {
		invalid := types.Validation{
			Valid: false,
			Error: fmt.Sprintf("Invalid token: unknown issuer %q\n", issuerString),
		}
		return invalid, nil
	}

	if issuer.PublicKey == nil {
		if err := issuer.parsePublicKey(); err != nil {
			return types.Validation{}, fmt.Errorf("parse public key for issuer %q: %v", issuerString, err)
		}
	}

	var publicClaims joseJWT.Claims
	if err := token.Claims(issuer.PublicKey, &publicClaims); err != nil {
		invalid := types.Validation{
			Valid: false,
			Error: fmt.Sprintf("Failed to extract token claims: %v\n", err),
		}
		return invalid, nil
	}

	expected := joseJWT.Expected{
		Audience: c.Audience,
		Time:     time.Now(),
	}

	if err := publicClaims.Validate(expected); err != nil {
		invalid := types.Validation{
			Valid: false,
			Error: fmt.Sprintf("Invalid token: %v\n", err),
		}
		return invalid, nil
	}

	identity, err := getIdentity(*token, c.Claims)
	if err != nil {
		invalid := types.Validation{
			Valid: false,
			Error: fmt.Sprintf("Invalid token: %v\n", err),
		}
		return invalid, nil
	}

	validation := types.Validation{
		Valid:  true,
		Claims: identity,
	}

	assertion := types.Assertion{
		IDPrefix:      issuer.IDPrefix,
		SubjectPrefix: issuer.SubjectPrefix,
		GroupPrefix:   issuer.GroupPrefix,
	}

	validation.Assert(assertion)

	return validation, nil
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

func getIssuer(token joseJWT.JSONWebToken) (string, error) {
	var publicClaims joseJWT.Claims
	if err := token.UnsafeClaimsWithoutVerification(&publicClaims); err != nil {
		return "", err
	}

	return publicClaims.Issuer, nil
}

func getIdentity(token joseJWT.JSONWebToken, fields JWTClaims) (types.Claims, error) {
	var claims interface{}
	if err := token.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return types.Claims{}, err
	}

	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		return types.Claims{}, fmt.Errorf("failed to cast JWT claims to map[string]interface{}")
	}

	id, ok := claimsMap[fields.ID].(string)
	if !ok {
		return types.Claims{}, fmt.Errorf("failed to cast %q claim to string", fields.ID)
	}

	subject, ok := claimsMap[fields.Subject].(string)
	if !ok {
		return types.Claims{}, fmt.Errorf("failed to cast %q claim to string", fields.Subject)
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
		ID:      id,
		Subject: subject,
		Groups:  groups,
		Extra:   claimsMap[fields.Extra],
	}

	return c, nil
}
