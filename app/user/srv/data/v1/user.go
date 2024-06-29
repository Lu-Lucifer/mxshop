package v1

import (
	"context"
	"gorm.io/gorm"
	metav1 "mxshop/pkg/common/meta/v1"
	"time"
)

// 基本模型字段
type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	//软删除
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type UserDO struct {
	BaseModel `json:"base_model"`
	Mobile    string `gorm:"index:idx_mobile;unique;type:varchar(11);not null" json:"mobile,omitempty"`
	Password  string `gorm:"type:varchar(100);not null" json:"password,omitempty"`
	NickName  string `gorm:"type:varchar(20)" json:"nick_name,omitempty"`
	// 日期保存容易报错，这里用指针类型
	Birthday *time.Time `gorm:"type:datetime" json:"birthday,omitempty"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'" json:"gender,omitempty"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'" json:"role,omitempty"`
}

func (u *UserDO) TableName() string {
	return "user"
}

type UserDOList struct {
	TotalCount int64     //总数
	Items      []*UserDO //数据
}

//	func List(ctx context.Context,opts metav1.ListMeta)(*UserDOList,error){
//		return &UserDOList{},nil
//	}

type UserStore interface {
	/*
		有数据访问的方法，一定要有error
		参数中最好有ctx，可以cancel，另外可以传入value值，或链路追踪
	*/
	//用户列表
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDOList, error)
	//通过手机号查询用户
	GetByMobile(ctx context.Context, mobile string) (*UserDO, error)
	GetByID(ctx context.Context, id uint64) (*UserDO, error)
	Create(ctx context.Context, user *UserDO) error
	Update(ctx context.Context, user *UserDO) error
}
