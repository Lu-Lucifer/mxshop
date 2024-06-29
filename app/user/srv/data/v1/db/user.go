package db

import (
	"context"
	"gorm.io/gorm"
	"mxshop/app/pkg/code"
	dv1 "mxshop/app/user/srv/data/v1"
	code2 "mxshop/gmicro/code"
	metav1 "mxshop/pkg/common/meta/v1"
	"mxshop/pkg/errors"
)

type users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *users {
	return &users{
		db: db,
	}
}

var _ dv1.UserStore = &users{}

// List
//
//	@Description: 获取用户列表，都需要分页
//	@receiver u
//	@param ctx
//	@param orderby
//	@param opts
//	@return *dv1.UserDOList
//	@return error
func (u users) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*dv1.UserDOList, error) {
	ret := dv1.UserDOList{}
	var limit, offset int
	if opts.PageSize == 0 {
		limit = 10
	} else {
		limit = opts.PageSize
	}
	if opts.Page > 0 {
		offset = (opts.Page - 1) * limit
	}
	//排序
	query := u.db
	for _, value := range orderby {
		query = query.Order(value)
	}
	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return &ret, nil
}

// GetByMobile
//
//	@Description:根据手机号查询用户
//	@receiver u
//	@param ctx
//	@param mobile
//	@return *dv1.UserDO
//	@return error
func (u users) GetByMobile(ctx context.Context, mobile string) (*dv1.UserDO, error) {
	user := dv1.UserDO{}
	err := u.db.Where("mobile =?", mobile).First(&user).Error
	// err是gorm的err，这种error不建议往上抛
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

// GetByID
//
//	@Description:根据id查询用户
//	@receiver u
//	@param ctx
//	@param id
//	@return *dv1.UserDO
//	@return error
func (u users) GetByID(ctx context.Context, id uint64) (*dv1.UserDO, error) {
	user := dv1.UserDO{}
	err := u.db.First(&user, id).Error
	// err是gorm的err，这种error不建议往上抛
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &user, nil
}

// Create
//
//	@Description: 创建用户
//	@receiver u
//	@param ctx
//	@param user
//	@return error
func (u users) Create(ctx context.Context, user *dv1.UserDO) error {
	tx := u.db.Create(user)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil

}

func (u users) Update(ctx context.Context, user *dv1.UserDO) error {
	tx := u.db.Save(user)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}
