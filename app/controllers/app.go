package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	activation := c.Params.Get("activation")
	return c.Render(activation)
}
