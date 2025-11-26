package grpcclient

import (
	"context"
	"fmt"
)

func Example() {
	profile, err := New(Config{Host: "localhost", Port: "50051"})
	if err != nil {
		panic(err)
	}

	defer profile.Close()

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

	err = profile.Delete(ctx, id.String())
	if err != nil {
		panic(err)
	}

	_, err = profile.GetProfile(ctx, id.String())

	fmt.Println("Get request:", err)
}
