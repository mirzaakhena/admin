package dao

import "github.com/mirzaakhena/admin/model"

// ISpaceDao is
type ISpaceDao interface {
	IsUniqueNamePerUserID(dc model.DaoContext, spaceName, userID string) bool
	Create(dc model.DaoContext, bu model.Space) error
	GetOne(dc model.DaoContext, ID string) *model.Space
	GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.Space, string)
	Delete(dc model.DaoContext, ID string) error
	Update(dc model.DaoContext, ID string, obj *model.Space) error
}

// SpaceDao is
type SpaceDao struct {
}
