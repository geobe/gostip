package model

import (
	"fmt"
	"golang.org/x/crypto/scrypt"
)

type User struct {
	Model
	Login    string
	Fullname string
	Password string
	Role     int
}

const U_ANY int = 1
const U_ENROL int = 2
const U_POFF int = 4
const U_UADMIN int = 8
const U_FULLADMIN int = 128

func NewUser(l, n, pw string, roles ...int) *User {
	var r int
	if len(roles) == 0 {
		r = U_ANY
	} else {
		for _, v := range roles {
			r |= v
		}
	}
	return &User{
		Login:    l,
		Fullname: n,
		Password: encrypt(pw + l),
		Role:     r,
	}
}

func (u *User) ValidatePw(pw string) bool {
	return u.Password == encrypt(pw+u.Login)
}

func (u *User) ChangePw(oldpw, newpw string) bool {
	if u.ValidatePw(oldpw) {
		u.Password = encrypt(newpw)
		return true
	} else {
		return false
	}
}

func (u *User) Authorize(role int) bool {
	return u.Role|role != 0
}

func (u *User) String() string {
	return fmt.Sprintf("%s: %s [%s] %d", u.Login, u.Fullname, u.Password, u.Role)
}

func encrypt(pw string) string {
	dk, err := scrypt.Key([]byte(pw), []byte("#K?@1"), 16384, 4, 1, 32)
	if err != nil {
		panic(fmt.Errorf("Encryption error %s \n", err))
	}
	return fmt.Sprintf("%x", dk)
}

func Encrypt(pw string) string {
	return encrypt(pw)
}
