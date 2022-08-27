package main

import "fmt"

func main() {
	app, err := initApp("Hello! ", "Goodbye!", "Hi! ", "Bye bye!")
	if err != nil {
		return
	}
	app.foo.Print()
	fmt.Println()
	app.bar.Print()
	fmt.Println()
	app.aMao.Print()
}
