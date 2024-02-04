package main

import (
	"github.com/kerem-kaynak/recurb/internal/auth"
	"github.com/kerem-kaynak/recurb/internal/routes"
)

func main() {
	auth.NewAuth()
	router := routes.NewRouter()
	router.Run()
}
