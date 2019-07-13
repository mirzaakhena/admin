package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/mirzaakhena/admin/model"
)

// IUserPermissionDao is
type IUserPermissionDao interface {
	GetUserPermission(dc model.DaoContext, userID, spaceID, code string) *model.UserPermission
	Create(dc model.DaoContext, bu model.UserPermission) error
	GetOne(dc model.DaoContext, spaceID, userID string) *model.UserPermission
	GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.UserPermission, uint64)
	Delete(dc model.DaoContext, ID string) error
	Update(dc model.DaoContext, ID string, obj *model.UserPermission) error
}

// UserPermissionDao is
type UserPermissionDao struct {
}

// GetUserPermission is
func (u *UserPermissionDao) GetUserPermission(dc model.DaoContext, userID, spaceID, code string) *model.UserPermission {

	var obj model.UserPermission
	dc.(*gorm.DB).First(&obj, "user_id = ? AND space_id = ? AND method_endpoint = ?", userID, spaceID, code)
	if obj.ID == "" {
		return nil
	}
	return &obj

}

// Create is
func (u *UserPermissionDao) Create(dc model.DaoContext, bu model.UserPermission) error {
	return nil
}

// GetOne is
func (u *UserPermissionDao) GetOne(dc model.DaoContext, spaceID, userID string) *model.UserPermission {
	return nil
}

// GetAll is
func (u *UserPermissionDao) GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.UserPermission, uint64) {
	return nil, 0
}

// Delete is
func (u *UserPermissionDao) Delete(dc model.DaoContext, ID string) error {
	return nil
}

// Update is
func (u *UserPermissionDao) Update(dc model.DaoContext, ID string, obj *model.UserPermission) error {
	return nil
}
