package dao

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/mirzaakhena/admin/model"
	"github.com/mirzaakhena/common/utils"
)

// IUserDao is
type IUserDao interface {
	Create(dc model.DaoContext, bu *model.User) error
	GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.User, uint)
	GetOneByEmail(dc model.DaoContext, email string) *model.User
	GetOneByID(dc model.DaoContext, id string) *model.User
	Delete(dc model.DaoContext, ID string) error
	Update(dc model.DaoContext, obj *model.User) error
}

// UserDao is
type UserDao struct {
}

// GetAll is
func (b *UserDao) GetAll(dc model.DaoContext, req model.GetAllBasicRequest) ([]model.User, uint) {

	var objs []model.User
	var count uint

	query := dc.(*gorm.DB).Model(&model.User{})

	if req.Filter != nil {

		// filtering
		for k, v := range req.Filter {
			query = query.Where(fmt.Sprintf("%s LIKE ?", utils.SnakeCase(k)), fmt.Sprintf("%s%%", v))
		}
	}

	// count
	query.Count(&count)

	// sorting
	if req.SortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", utils.SnakeCase(req.SortBy), req.SortDir))
	}

	if req.PageNumber == 0 {
		req.PageNumber = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 1
	}

	// paging
	query = query.Offset((req.PageNumber - 1) * req.PageSize).Limit(req.PageSize)

	query.Find(&objs)
	return objs, count

}

// Create is
func (b *UserDao) Create(dc model.DaoContext, bu *model.User) error {
	return dc.(*gorm.DB).Create(bu).Error
}

// GetOneByEmail is
func (b *UserDao) GetOneByEmail(dc model.DaoContext, email string) *model.User {
	var obj model.User
	dc.(*gorm.DB).First(&obj, "email = ?", email)
	if obj.ID == "" {
		return nil
	}
	return &obj
}

// GetOneByID is
func (b *UserDao) GetOneByID(dc model.DaoContext, id string) *model.User {
	var obj model.User
	dc.(*gorm.DB).First(&obj, "id = ?", id)
	if obj.ID == "" {
		return nil
	}
	return &obj
}

// Delete is
func (b *UserDao) Delete(dc model.DaoContext, id string) error {
	return dc.(*gorm.DB).Delete(model.User{}, "id = ?", id).Error
}

// Update is
func (b *UserDao) Update(dc model.DaoContext, obj *model.User) error {
	return dc.(*gorm.DB).Save(obj).Error
}
