package models

import "errors"

type User struct {
	First string
	Last  string
}

func (u *User) Validate() error {
	if u.First == "" {
		return errors.New("first name can't be blank")
	}
	if u.Last == "" {
		return errors.New("last name can't be blank")
	}
	return nil
}
