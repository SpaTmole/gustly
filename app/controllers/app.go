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
	err = app.DB.SelectOne(&_unusedUser, "select * from \"User\" where \"Email\"=$1", user.Email)
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
	fmt.Println(activation_key)
	var user models.RegistrationProfile
	var err error
	var errorMessage string
	err = app.DB.SelectOne(&user, "select * from \"RegistrationProfile\" where \"ActivationKey\"=$1", activation_key)
	if err != nil {
		revel.ERROR.Println(err)
		errorMessage = "Activation key is invalid."
		c.Validation.Error(errorMessage).Key("activation_key")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.RenderActivation, activation_key)
	}
	if user.IsExpired() || user.Activated {
		errorMessage = "Activation has expired."
		c.Validation.Error(errorMessage).Key("activation_key")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.RenderActivation, activation_key)
	}
	var account = models.User{
		Active:   true,
		Staff:    false,
		Username: user.Username,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: c.Params.Form.Get("password"),
		Verify:   c.Params.Form.Get("verify"),
	}

	fmt.Println(account)
	account.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.RenderActivation, activation_key)
	}
	err = app.DB.Insert(&account)
	if err != nil {
		panic(err)
	}
	user.Activate()
	_, err = app.DB.Update(&user)
	if err != nil {
		panic(err)
	}
	// Create User from password and Registration profile.
	// Deactivate the Registration Profile.
	return c.Redirect("/?activation=success")
}

func (c App) RenderActivation(activation_key string) revel.Result {
	return c.RenderTemplate("App/Activation.html")
}
