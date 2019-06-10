package output

import (
	"encoding/json"
	"net/http"

	auth "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/robbiemcmichael/auth-mux/internal/types"
)

type KubernetesTokenReview struct {
	Audience string `yaml:"audience"`
	MaxTTL   int64  `yaml:"maxTTL"`
}

func (o *KubernetesTokenReview) Handler(w http.ResponseWriter, validation types.Validation) error {
	tokenReview := auth.TokenReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: auth.SchemeGroupVersion.String(),
			Kind:       "TokenReview",
		},
		Status: auth.TokenReviewStatus{
			Authenticated: validation.Valid,
			User: auth.UserInfo{
				UID:      validation.Claims.ID,
				Username: validation.Claims.Subject,
				Groups:   validation.Claims.Groups,
			},
			Error: validation.Error,
		},
	}

	return json.NewEncoder(w).Encode(tokenReview)
}
