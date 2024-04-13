package main

import "ai-detective/backend"

func main() {
	b := backend.NewServer()
	b.RunServer()
}
