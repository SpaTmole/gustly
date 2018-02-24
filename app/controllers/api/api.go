package api

import (
	"github.com/SpaTmole/gustly/app/controllers"
	"github.com/revel/revel"
	// "fmt"
	// "github.com/SpaTmole/gustly/app"
	// "github.com/SpaTmole/gustly/app/models"
	// "net/http"
	// "strings"
)

type APIv1 struct {
	*revel.Controller
}

func init() {
	revel.InterceptFunc(controllers.AuthenticationMiddleware, revel.BEFORE, &APIv1{})
}
