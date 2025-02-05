package main

import (
	"github.com/omaradriano/httpserver/internal/server"
)

const portNum string = ":3333"

func main(){

	srv := server.NewHttpServer(":3333")
	srv.ListenAndServe()

}