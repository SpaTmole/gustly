package controllers

import (
	"encoding/json"
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
	var count = 0
	var user models.RegistrationProfile
	c.Params.BindJSON(&user)
	fmt.Println(user)
	// c.Validation.Required(user.Verify == user.Password).MessageKey("Passwords don't match").Key("user.verify")
	errors := user.Validate(c.Validation)
	if errors != nil {
		return c.RenderJSON(map[string]interface{}{"result": "fail", "errors": errors})
	}
	app.DB.Model(&models.User{}).Where("email = ?", user.Email).Or("username = ?", user.Username).Count(&count)
	if count != 0 {
		return c.RenderJSON(map[string]interface{}{"result": "fail", "errors": [1]string{"User with this email or username already exist."}})
	}
	activation_key := user.GenerateKey()
	err = app.DB.Create(&user).Error
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
	var user models.RegistrationProfile
	var err error
	var errorMessage string
	var form models.PasswordSubmition
	err = app.DB.Where("activation_key = ?", activation_key).First(&user).Error
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
	c.Params.Bind(&form.Password, "password")
	c.Params.Bind(&form.Verify, "verify")
	form.Validate(c.Validation)
	if c.Validation.HasErrors() {
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
	}
	account.SavePassword(form.Password)
	err = app.DB.Create(&account).Error
	if err != nil {
		panic(err)
	}
	user.Activate()
	err = app.DB.Save(&user).Error
	if err != nil {
		panic(err)
	}
	return c.Redirect("/?activation=success")
}

func (c App) RenderActivation(activation_key string) revel.Result {
	return c.RenderTemplate("App/Activation.html")
}

func (c App) Login() revel.Result {
	var err error
	defer func(err error) {
		if err != nil {
			revel.ERROR.Println(err)
		}
	}(err)
	var credentials = models.Credentials{}
	var user = models.User{}
	var jsonResponse map[string]interface{}
	c.Params.BindJSON(&credentials)
	errors := credentials.Validate(c.Validation)
	if errors != nil {
		return c.RenderJSON(map[string]interface{}{"result": "fail", "errors": errors})
	}
	err = app.DB.Where("username = ?", credentials.Username).First(&user).Error
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"result": "fail", "errors": []string{"Username or password is incorrect"}})
	}
	revel.INFO.Println(user)
	token, expires := user.Login(&credentials)
	if token == nil {
		return c.RenderJSON(map[string]interface{}{"result": "fail", "errors": []string{"Username or password is incorrect"}})
	}
	err = app.DB.Create(&token).Error
	if err != nil {
		panic(err)
	}
	marshaled, _ := json.Marshal(token)
	json.Unmarshal(marshaled, &jsonResponse)
	jsonResponse["expires_in"] = expires
	return c.RenderJSON(jsonResponse)
}

func (c App) Logout() revel.Result {
	return c.RenderJSON(map[string]interface{}{"result": "fail"})
}
