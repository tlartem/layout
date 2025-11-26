//go:build http_v1

package test

import "gitlab.noway/pkg/httpclient"

type ProfileClient = httpclient.Client

func BuildProfile(s *Suite) {
	s.profile = httpclient.New(httpclient.Config{Host: "localhost", Port: "8080"})
}
