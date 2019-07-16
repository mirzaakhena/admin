package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/mirzaakhena/admin/model"
)

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

// NewUserSpaceDao is
func NewUserSpaceDao() *UserSpaceDao {
	return &UserSpaceDao{}
}

// Create is
func (u *UserSpaceDao) Create(dc model.DaoContext, bu model.UserSpace) error {
	return dc.(*gorm.DB).Create(&bu).Error
}

// GetOne is
func (u *UserSpaceDao) GetOne(dc model.DaoContext, spaceID, userID string) *model.UserSpace {
	return nil
}

// GetAll is
func (u *UserSpaceDao) GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.UserSpace, uint) {
	var objs []model.UserSpace
	var count uint

	query := dc.(*gorm.DB).Model(&model.UserSpace{})

	// // filtering
	// for k, v := range req.Filter {
	// 	query = query.Where(fmt.Sprintf("%s LIKE ?", utils.SnakeCase(k)), fmt.Sprintf("%s%%", v))
	// }

	// count
	query.Count(&count)

	// sorting
	// if req.SortBy != "" {
	// 	query = query.Order(fmt.Sprintf("%s %s", utils.SnakeCase(req.SortBy), req.SortDir))
	// }

	// paging
	query = query.Offset((req.PageNumber - 1) * req.PageSize).Limit(req.PageSize)

	query.Preload("Space").Find(&objs)
	return objs, count
}

// Delete is
func (u *UserSpaceDao) Delete(dc model.DaoContext, ID string) error {
	return nil
}

// Update is
func (u *UserSpaceDao) Update(dc model.DaoContext, ID string, obj *model.UserSpace) error {
	return nil
}
