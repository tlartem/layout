//go:build integration

package test

func (s *Suite) Test_GetProfile_Ok() {
	id, err := s.profile.Create(ctx, "John_Get", 25, "john@gmail.com", "+73003002020")
	s.NoError(err)

	p, err := s.profile.GetProfile(ctx, id.String())
	s.NoError(err)

	s.Equal("John_Get", p.Name)
	s.Equal(25, p.Age)
	s.Equal("john@gmail.com", p.Contacts.Email)
	s.Equal("+73003002020", p.Contacts.Phone)
}

func (s *Suite) Test_GetProfile_NotFound() {
	_, err := s.profile.GetProfile(ctx, "c6799c89-c560-45a2-afda-b3f1eb9bee2b")
	s.NotNil(err)

	s.Contains(err.Error(), "not found")
}
