package dao

import "github.com/mirzaakhena/admin/model"

// IUserSpaceDao is
type IUserSpaceDao interface {
	Create(dc model.DaoContext, bu model.UserSpace) error
	GetOne(dc model.DaoContext, spaceID, userID string) *model.UserSpace
	GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.UserSpace, uint)
	Delete(dc model.DaoContext, ID string) error
	Update(dc model.DaoContext, ID string, obj *model.UserSpace) error
}

// UserSpaceDao is
type UserSpaceDao struct {
}

func (u *UserSpaceDao) Create(dc model.DaoContext, bu model.UserSpace) error {
	return nil
}

func (u *UserSpaceDao) GetOne(dc model.DaoContext, spaceID, userID string) *model.UserSpace {
	return nil
}

func (u *UserSpaceDao) GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.UserSpace, uint) {
	return nil, 0
}

func (u *UserSpaceDao) Delete(dc model.DaoContext, ID string) error {
	return nil
}

func (u *UserSpaceDao) Update(dc model.DaoContext, ID string, obj *model.UserSpace) error {
	return nil
}
