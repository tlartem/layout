//go:build http_v2

package test

func (s *Suite) Test_A_GetProfiles_Ok() {
	_, err := s.profile.Create(ctx, "John1_Get", 25, "john1@gmail.com", "+73003002021")
	s.NoError(err)
	_, err = s.profile.Create(ctx, "John2_Get", 26, "john2@gmail.com", "+73003002022")
	s.NoError(err)

	profiles, err := s.profile.GetProfiles(ctx, "name", "asc", 0, 10)
	s.NoError(err)

	s.Equal(2, len(profiles))

	p := profiles[0]

	s.Equal("John1_Get", p.Name)
	s.Equal(25, p.Age)
	s.Equal("john1@gmail.com", p.Contacts.Email)
	s.Equal("+73003002021", p.Contacts.Phone)

	p = profiles[1]

	s.Equal("John2_Get", p.Name)
	s.Equal(26, p.Age)
	s.Equal("john2@gmail.com", p.Contacts.Email)
	s.Equal("+73003002022", p.Contacts.Phone)
}
