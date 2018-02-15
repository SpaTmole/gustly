package controllers

import (
	"fmt"
	"github.com/SpaTmole/gustly/app"
	"github.com/SpaTmole/gustly/app/models"
	"github.com/revel/revel"
	"net/http"
	"strings"
	"time"
)

type Private struct {
	*revel.Controller
}

func (c Private) Logout() revel.Result {
	user, _ := c.Args["user"].(models.User)
	fmt.Println("--->>>>>", user)
	err := app.DB.Table("tokens").Where("user_id = ?", user.ID).UpdateColumns(map[string]interface{}{
		"expires_at": 0,
	}).Error
	if err != nil {
		panic(err)
	}
	return c.RenderJSON(map[string]interface{}{"result": "success"})
}

//
// ----
//

// Unauthorized returns an HTTP 401 Forbidden response whose body is the
// formatted string of msg and objs.
func Unauthorized(c *revel.Controller, msg string, objs ...interface{}) revel.Result {
	finalText := msg
	if len(objs) > 0 {
		finalText = fmt.Sprintf(msg, objs...)
	}
	c.Response.Status = http.StatusUnauthorized
	return c.RenderJSON(map[string]interface{}{"status": "Unauthorized", "message": finalText})
}

func authenticationMiddleware(c *revel.Controller) revel.Result {
	var err error
	var user models.User
	var token models.Token
	bearerToken := c.Request.Header.Get("Authorization")
	authToken := strings.Replace(bearerToken, "Bearer ", "", 1)
	err = app.DB.Where("auth_token = ? AND expires_at > ?", authToken, time.Now().Unix()).First(&token).Error
	if err != nil {
		return Unauthorized(c, "Access denied.")
	}
	fmt.Println(token)
	err = app.DB.Model(&token).Related(&user).Error
	if err != nil {
		panic(err)
	}
	c.Args["user"] = user
	return nil
}

func init() {
	revel.InterceptFunc(authenticationMiddleware, revel.BEFORE, &Private{})
}
