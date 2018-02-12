package controllers

import (
	"fmt"
	"github.com/SpaTmole/gustly/app"
	"github.com/SpaTmole/gustly/app/models"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) SignUp() revel.Result {
	var user models.User
	c.Params.BindJSON(&user)
	fmt.Println(user)
	c.Validation.Required(user.Verify == user.Password).MessageKey("Passwords don't match").Key("user.verify")
	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		return c.RenderJSON(map[string]interface{}{"result": "fail", "errors": c.Validation.ErrorMap()})
	}
	user.SavePassword()
	err := app.DB.Insert(&user)
	if err != nil {
		panic(err)
	}
	return c.RenderJSON(map[string]string{"result": "success"})
}
