package models

import (
	"github.com/revel/revel"
)

type Credentials struct {
	Username string `json:username`
	Password string `json:password`
}

type PasswordSubmition struct {
	Password string `json:password`
	Verify   string `json:verify`
}

func (c *Credentials) Validate(v *revel.Validation) map[string]*revel.ValidationError {
	v.Check(c.Username, revel.Required{}).Key("username")
	v.Check(c.Password, revel.Required{}).Key("password")
	if v.HasErrors() {
		return v.ErrorMap()
	}
	return nil
}

func (p *PasswordSubmition) Validate(v *revel.Validation) map[string]*revel.ValidationError {
	v.Check(p.Password,
		revel.Required{},
		revel.MinSize{5},
	).Key("user.Password")
	v.Check(p.Verify, revel.Required{}).Key("user.Verify")
	v.Required(p.Password == p.Verify).MessageKey("Passwords don't match").Key("user.Verify")

	if v.HasErrors() {
		return v.ErrorMap()
	}
	return nil
}
