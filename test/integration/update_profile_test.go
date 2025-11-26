//go:build integration

package test

func (s *Suite) Test_UpdateProfile() {
	id, err := s.profile.Create(ctx, "John_Update", 25, "john@gmail.com", "+73003002020")
	s.NoError(err)

	p, err := s.profile.GetProfile(ctx, id.String())
	s.NoError(err)

	s.Equal(25, p.Age)

	var (
		name  = "New John_Update"
		age   = 26
		email = "new-john@gmail.com"
		phone = "+73003004000"
	)

	err = s.profile.Update(ctx, id.String(), &name, &age, &email, &phone)
	s.NoError(err)

	p, err = s.profile.GetProfile(ctx, id.String())
	s.NoError(err)

	s.Equal(name, p.Name)
	s.Equal(age, p.Age)
	s.Equal(email, p.Contacts.Email)
	s.Equal(phone, p.Contacts.Phone)
}
