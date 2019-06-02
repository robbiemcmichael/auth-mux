package input

import (
	"fmt"
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

func (i *KubernetesTokenReview) Config() string {
	return fmt.Sprintf("%+v", i)
}
