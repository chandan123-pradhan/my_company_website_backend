package main

import (
	"fmt"
	"log"
	"net/http"

	"blog_project.com/controllers"
	"blog_project.com/routers"
)

func main() {
    fmt.Println("Hello, world!")
    controllers.InitDB()
    defer controllers.DB.Close()
    controllers.Initialize(controllers.DB);

    r:= routers.NewRouter();

    log.Println("starting server on : 8080")
    if err:= http.ListenAndServe(":8080", r);err!=nil{
        log.Fatal(err)
    }
}