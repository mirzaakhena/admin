package model

import (
	"time"

	"github.com/gin-gonic/gin"
)

// BaseModel is
type BaseModel struct {
	ID        string     `json:"id" gorm:"primary_key"` //
	CreatedAt time.Time  `json:"-"`                     //
	UpdatedAt time.Time  `json:"-"`                     //
	DeletedAt *time.Time `json:"-" sql:"index"`         //
}

// Space is masterdata
type Space struct {
	BaseModel
	Name             string    ``         //
	Description      string    ``         //
	Expired          time.Time `json:"-"` //
	MaxUser          int       `json:"-"` //
	TotalCurrentUser int       `json:"-"` //
}

// SpaceRole is Role that exist in the space
type SpaceRole struct {
	BaseModel
	Name    string `json:"name"`    //
	SpaceID string `json:"spaceId"` //
}

// RolePermission is
type RolePermission struct {
	BaseModel
	SpaceID        string `json:"spaceId"`        //
	SpaceRoleID    string `json:"spaceRoleId"`    //
	MethodEndpoint string `json:"methodEndpoint"` //
}

// Permission is
type Permission struct {
	Method      string          `json:"method"`      //
	Endpoint    string          `json:"endpoint"`    //
	Function    gin.HandlerFunc `json:"function"`    //
	Description string          `json:"description"` //
	Category    string          `json:"category"`    //
}

// User is
type User struct {
	BaseModel
	Name               string `json:"name"`    //
	Email              string `json:"email"`   //
	Phone              string `json:"phone"`   //
	Address            string `json:"address"` //
	Password           string `json:"-"`       //
	LoginToken         string `json:"-"`       //
	ResetPasswordToken string `json:"-"`       //
	Status             string `json:"-"`       //
}

// UserSpace is
type UserSpace struct {
	BaseModel
	SpaceID string `json:"-"` //
	UserID  string `json:"-"` //
	Type    string `json:"-"` // ADMIN | USER
	Status  string `json:"-"` // ACTIVE | REQUEST | SUSPEND
}

// UserPermission is
type UserPermission struct {
	BaseModel
	MethodEndpoint string `json:"methodEndpoint"` //
	UserID         string `json:"userId"`         //
	SpaceID        string `json:"spaceId"`        //
}

// UserRole is
type UserRole struct {
	BaseModel
	UserID      string `json:"userId"`      //
	SpaceRoleID string `json:"spaceRoleId"` //
}
