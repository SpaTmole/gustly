package models

import (
	"bytes"
	"fmt"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type User struct {
	Active         bool   `json:is_active db:"is_active"`
	Id             int    `json:id db:"id, primarykey, autoincrement"`
	Friends        []User `json:friends db:"friends"`
	Name           string `json:name db:"name"`
	Username       string `json:username db:"username, primarykey"`
	Phone          string `json:phone db:"phone"`
	Email          string `json:email db:"email"`
	Password       string `json:password db:"-"`
	Verify         string `json:verify db:"-"`
	hashedPassword []byte `json:"-" db:"hashed_password"`
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
	ValidatePassword(v, user.Password).
		Key("user.password")

	v.Check(user.Username,
		revel.Required{},
		revel.MinSize{4},
		revel.MaxSize{100},
	).Key("user.username")
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{5},
	)
}

func (user *User) SavePassword() {
	user.hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
}

func (user *User) CheckPassword(password string) bool {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return bytes.Equal(hashedPassword, user.hashedPassword)
}

func (u *User) Register() {
	return
}
