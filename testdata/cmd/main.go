package main

func main() {
	app, err := initApp("Hello! ", "Goodbye!")
	if err != nil {
		return
	}
	app.foo.Print()
}
