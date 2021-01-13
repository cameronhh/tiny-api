package main

import "os"

func main() {
	a := App{}
	a.Initialize(os.Getenv("DB_CONNECTION"))
	a.Run()
}
