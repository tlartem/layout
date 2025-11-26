//go:build http_v2

package test

import (
	"gitlab.noway/pkg/httpclientv2"
)

type ProfileClient = httpclientv2.Client

func BuildProfile(s *Suite) {
	var err error
	s.profile, err = httpclientv2.New(httpclientv2.Config{Address: "http://localhost:8080/noway/layout/api/v2"})
	s.NoError(err)
}
