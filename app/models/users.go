package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"regexp"
	"time"
)

type User struct {
	gorm.Model
	Staff          bool   `json:is_staff gorm:"default:false;"`
	Active         bool   `json:is_active gorm:"default:true;"`
	Name           string `json:name`
	Username       string `json:username gorm:"size:100;unique_index;not null"`
	Phone          string `json:phone`
	Email          string `json:email gorm:"not null"`
	HashedPassword string `json:"-" gorm:"size:255"`
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
	User      User   `json:"-"`
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var (
	userRegex        = regexp.MustCompile("^\\w*$")
	letterRunes      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	expirationPeriod = int64(7200) // TODO: Move to the config.
)

func (token *Token) Make() {
	token.ExpiresAt = time.Now().Unix() + expirationPeriod
	token.AuthToken = (&RegistrationProfile{}).GenerateKey()
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
		token = &Token{User: *user}
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
	starting_time := time.Now()
	rand.Seed(starting_time.Unix())
	buff := make([]rune, 32)
	for idx := range buff {
		buff[idx] = letterRunes[rand.Intn(len(letterRunes))]
	}
	p.ActivationKey = string(buff)
	p.Expires = starting_time.AddDate(0, 0, 12).Unix() // 12 Days TODO: Move to the config.
	return string(buff)
}

func (p *RegistrationProfile) Activate() {
	p.Activated = true
}

func (p *RegistrationProfile) IsExpired() bool {
	return time.Now().Unix() >= p.Expires
}
