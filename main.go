package main

import "github.com/klaxn/klaxn-api/pkg/app"

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}

	err = a.Run()
	if err != nil {
		a.Logger.Fatal(err)
	}
}
