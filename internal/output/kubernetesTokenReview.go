package output

import (
	"fmt"
)

type KubernetesTokenReview struct {
	Audience string `yaml:"audience"`
	MaxTTL   int64  `yaml:"maxTTL"`
}

func (o *KubernetesTokenReview) Config() string {
	return fmt.Sprintf("%+v", o)
}
