package controllers

import (
	"fmt"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) SignUp() revel.Result {
	var jsonData map[string]interface{}
	c.Params.BindJSON(&jsonData)
	fmt.Println(jsonData)
	return c.RenderJSON(map[string]string{"result": "success"})
}
