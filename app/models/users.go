package models

import (
	"bytes"
	"fmt"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"regexp"
	"time"
)

type User struct {
	// Friends        []User `json:friends db:"friends"`
	Staff          bool   `json:is_staff db:"is_staff"`
	Active         bool   `json:is_active db:"is_active"`
	Id             int    `json:id db:"id, primarykey, autoincrement"`
	Name           string `json:name db:"name"`
	Username       string `json:username db:"username, primarykey"`
	Phone          string `json:phone db:"phone"`
	Email          string `json:email db:"email"`
	Password       string `json:password db:"-"`
	Verify         string `json:verify db:"-"`
	HashedPassword []byte `json:"-" db:"hashed_password"`
}

type RegistrationProfile struct {
	Id            string `json:id db:"id, primarykey, autoincrement"`
	Username      string `json:username db:"username, primarykey"`
	Phone         string `json:phone db:"phone"`
	Email         string `json:email db:"email"`
	ActivationKey string `json:activation_key db:"activation_key"`
	Expires       int64  `json:expires_at db:"expires_at"`
	Activated     bool   `json:is_active db:"is_active"`
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var (
	userRegex   = regexp.MustCompile("^\\w*$")
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func (user *User) Validate(v *revel.Validation) map[string]*revel.ValidationError {
	v.Check(user.Password,
		revel.Required{},
		revel.MinSize{5},
	).Key("user.Password")
	v.Check(user.Verify, revel.Required{}).Key("user.Verify")
	v.Required(user.Password == user.Verify).MessageKey("Passwords don't match").Key("user.Verify")

	if v.HasErrors() {
		return v.ErrorMap()
	}
	return nil
}

func (user *User) SavePassword() {
	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
}

func (user *User) CheckPassword(password string) bool {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return bytes.Equal(hashedPassword, user.HashedPassword)
}

func (u *User) Activate() {
	return
}

func (p *RegistrationProfile) Validate(v *revel.Validation) map[string]*revel.ValidationError {
	v.Check(p.Username, revel.Required{}, revel.MinSize{4}, revel.MaxSize{100}).Key("registration.Username")
	v.Check(p.Email, revel.Required{}).Key("registration.Email") // TODO: Add Regexp.
	if v.HasErrors() {
		return v.ErrorMap()
	}
	return nil
}

func (p *RegistrationProfile) GenerateKey() string {
	starting_time := time.Now()
	rand.Seed(starting_time.UnixNano())
	buff := make([]rune, 32)
	for idx := range buff {
		buff[idx] = letterRunes[rand.Intn(len(letterRunes))]
	}
	p.ActivationKey = string(buff)
	p.Expires = starting_time.AddDate(0, 0, 12).UnixNano() // 12 Days
	return string(buff)
}

func (p *RegistrationProfile) Activate() {
	p.Activated = true
}

func (p *RegistrationProfile) IsExpired() bool {
	return time.Now().UnixNano() >= p.Expires
}
