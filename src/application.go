package main

import (
	"fmt"
	"net/http"

	"github.com/copchase/user-status/controller"
	"github.com/julienschmidt/httprouter"
)

func main() {
	fmt.Println("Hello world!")
	router := httprouter.New()
	router.GET("/userinfo", controller.GetUserInfo)

	http.ListenAndServe(":5000", router)
}
