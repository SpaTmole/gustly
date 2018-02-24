package models

import (
	"fmt"
	"github.com/SpaTmole/gustly/app/utils"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	gorm.Model
	Staff          bool     `json:is_staff gorm:"default:false;"`
	Active         bool     `json:is_active gorm:"default:true;"`
	Name           string   `json:name`
	Username       string   `json:username gorm:"size:100;unique_index;not null"`
	Phone          string   `json:phone`
	Email          string   `json:email gorm:"not null"`
	HashedPassword string   `json:"-" gorm:"size:255"`
	Tokens         []Token  `json:"-"`
	Stories        []Story  `json:"-"`
	Friends        []*User  `json:"-" gorm:"many2many:friendships;association_jointable_foreignkey:friend_id"`
	Views          []*Story `json:"-" gorm:"many2many:views;"`
}

type RegistrationProfile struct {
	gorm.Model
	Username      string `json:username;not null`
	Phone         string `json:phone`
	Email         string `json:email;not null`
	ActivationKey string `json:activation_key`
	Expires       int64  `json:expires_at`
	Activated     bool   `json:is_active gorm:"default:false;"`
}

type Token struct {
	gorm.Model
	AuthToken string `json:auth_token`
	ExpiresAt int64  `json:expires_at`
	UserID    uint   // IT'S A DAMN! Documentation doesn't tell it, but it MUST be specified explicitly.
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var (
	expirationPeriod = int64(7200) // TODO: Move to the config.
)

func (token *Token) Make() {
	token.ExpiresAt = time.Now().Unix() + expirationPeriod
	token.AuthToken = utils.MakeUniqueKey(32)
}

func (user *User) SavePassword(password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.HashedPassword = string(hash)
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

func (user *User) Login(credentials *Credentials) (token *Token, expires int64) {
	if user.CheckPassword(credentials.Password) {
		token = &Token{}
		token.Make()
		expires = expirationPeriod
	}
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
	p.ActivationKey = utils.MakeUniqueKey(32)
	p.Expires = time.Now().AddDate(0, 0, 12).Unix() // 12 Days TODO: Move to the config.
	return p.ActivationKey
}

func (p *RegistrationProfile) Activate() {
	p.Activated = true
}

func (p *RegistrationProfile) IsExpired() bool {
	return time.Now().Unix() >= p.Expires
}
