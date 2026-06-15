package main

import (
	"fmt"
	"slices"
)

const (
	RoleBasic     = "basic"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)

type User interface {
	GetUsername() string
	HasPermission(permission string) bool
	GetRole() string
}

type BasicUser struct {
	username    string
	permissions []string
	role        string
}

func (b *BasicUser) GetUsername() string {
	return b.username
}

func (b *BasicUser) HasPermission(permission string) bool {
	return slices.Contains(b.permissions, permission)
}

func (b *BasicUser) GetRole() string {
	return b.role
}

func NewBasicUser(username string) *BasicUser {
	return &BasicUser{
		username:    username,
		permissions: []string{"read"},
		role:        RoleBasic,
	}
}

type Moderator struct {
	BasicUser
}

func NewModerator(username string) *Moderator {
	basic := NewBasicUser(username)
	basic.role = RoleModerator
	basic.permissions = append(basic.permissions, "edit", "ban_user")
	return &Moderator{BasicUser: *basic}
}

type Admin struct {
	Moderator
}

func NewAdmin(username string) *Admin {
	mod := NewModerator(username)
	mod.role = RoleAdmin
	mod.permissions = append(mod.permissions, "delete", "manage_roles")
	return &Admin{Moderator: *mod}
}

func main() {
	fmt.Println("Hello, World!")
}
