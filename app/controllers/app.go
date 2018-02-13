package controllers

import (
	"fmt"
	"github.com/SpaTmole/gustly/app"
	"github.com/SpaTmole/gustly/app/mail"
	"github.com/SpaTmole/gustly/app/models"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	activation := c.Params.Get("activation")
	return c.Render(activation)
}

func (c App) SignUp() revel.Result {
	var err error
	var user models.RegistrationProfile
	var _unusedUser models.User
	c.Params.BindJSON(&user)
	fmt.Println(user)
	// c.Validation.Required(user.Verify == user.Password).MessageKey("Passwords don't match").Key("user.verify")
	errors := user.Validate(c.Validation)
	if errors != nil {
		return c.RenderJSON(map[string]interface{}{"result": "fail", "errors": errors})
	}
	err = app.DB.SelectOne(&_unusedUser, "select * from user where email=?", user.Email)
	if err == nil {
		return c.RenderJSON(map[string]interface{}{"result": "fail", "errors": [1]string{"User with this email already exist."}})
	}
	// user.SavePassword()
	activation_key := user.GenerateKey()
	err = app.DB.Insert(&user)
	if err != nil {
		panic(err)
	}
	err = app.EmailService.SendMail(user.Email, []string{}, mail.MakeActivationMessage(activation_key))
	if err != nil {
		panic(err)
	}
	return c.RenderJSON(map[string]string{"result": "success"})
}

func (c App) Activate(activation_key string) revel.Result {
	fmt.Println(c.Params.Form)
	// Create User from password and Registration profile.
	// Deactivate the Registration Profile.
	return c.Redirect("/?activation=success")
}

func (c App) RenderActivation() revel.Result {
	return c.RenderTemplate("App/Activation.html")
}
