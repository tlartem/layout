//go:build integration

package test

func (s *Suite) Test_DeleteProfile() {
	id, err := s.profile.Create(ctx, "John_Delete", 25, "john@gmail.com", "+73003002020")
	s.NoError(err)

	p, err := s.profile.GetProfile(ctx, id.String())
	s.NoError(err)

	s.Equal("John_Delete", p.Name)
	s.Equal(25, p.Age)
	s.Equal("john@gmail.com", p.Contacts.Email)
	s.Equal("+73003002020", p.Contacts.Phone)

	err = s.profile.Delete(ctx, id.String())
	s.NoError(err)

	_, err = s.profile.GetProfile(ctx, id.String())
	s.Contains(err.Error(), "not found")
}
