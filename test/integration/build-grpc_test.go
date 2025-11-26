//go:build grpc

package test

import (
	"gitlab.noway/pkg/grpcclient"
)

type ProfileClient = grpcclient.Client

func BuildProfile(s *Suite) {
	var err error
	s.profile, err = grpcclient.New(grpcclient.Config{Host: "localhost", Port: "50051"})
	s.NoError(err)
}
