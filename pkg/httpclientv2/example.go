package httpclientv2

import (
	"context"
	"fmt"
)

func Example() { //nolint:funlen
	profile, err := New(Config{Address: "http://localhost:8080/noway/layout/api/v2"})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	id, err := profile.Create(ctx, "John", 25, "john@gmail.com", "+73003002020")
	if err != nil {
		panic(err)
	}

	p, err := profile.GetProfile(ctx, id.String())
	if err != nil {
		panic(err)
	}

	fmt.Println(p.ID)
	fmt.Println(p.Age)
	fmt.Println(p.Name)
	fmt.Println(p.Contacts.Email)
	fmt.Println(p.Contacts.Phone)

	var (
		name  = "John Doe"
		age   = 26
		email = "new-john@gmail.com"
		phone = "+73003004000"
	)

	err = profile.Update(ctx, id.String(), &name, &age, &email, &phone)
	if err != nil {
		panic(err)
	}

	p, err = profile.GetProfile(ctx, id.String())
	if err != nil {
		panic(err)
	}

	fmt.Println(p.ID)
	fmt.Println(p.Age)
	fmt.Println(p.Name)
	fmt.Println(p.Contacts.Email)
	fmt.Println(p.Contacts.Phone)

	profiles, err := profile.GetProfiles(ctx, "id", "asc", 0, 10)
	if err != nil {
		panic(err)
	}

	fmt.Println(profiles)

	err = profile.Delete(ctx, id.String())
	if err != nil {
		panic(err)
	}

	_, err = profile.GetProfile(ctx, id.String())

	fmt.Println("Get request:", err)
}
