package models

type Role string

const (
	RolePlayer     Role = "player"
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "superadmin"
)