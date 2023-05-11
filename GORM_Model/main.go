package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"//加个下划线是代表没有直接用到他
	"time"
)

//在使用ORM工具时，通常我们需要在代码中定义模型（Models）与数据库中的数据表进行映射，
//在GORM中模型（Models）通常是正常定义的结构体、基本的go类型或它们的指针。 同时也支持sql.Scanner及driver.Valuer接口（interfaces）。
//为了方便模型定义，GORM内置了一个gorm.Model结构体。
//gorm.Model是一个包含了ID, CreatedAt, UpdatedAt, DeletedAt四个字段的Golang结构体。
// gorm.Model 定义
//时间戳跟踪
//如果模型有 CreatedAt字段，该字段的值将会是初次创建记录的时间。
//如果模型有UpdatedAt字段，该字段的值将会是每次更新记录的时间。
//如果模型有DeletedAt字段，调用Delete删除该记录时，将会设置DeletedAt字段为当前时间，而不是直接将记录从数据库中删除。
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

//下面是自己定义的结构体模型
type User struct {    //表名默认就是结构体名称的复数Users
	gorm.Model   //内嵌gorm.Model
	Name         string
	Age          sql.NullInt64  //零值类型
	Birthday     *time.Time
	Email        string  `gorm:"type:varchar(100);unique_index"`//unique_index代表建立唯一的索引，不重复的
	Role         string  `gorm:"size:255"` // 设置字段大小为255
	MemberNumber *string `gorm:"unique;not null"` // 设置会员号（member number）唯一并且不为空
	Num          int     `gorm:"AUTO_INCREMENT"` // 设置 num 为自增类型
	Address      string  `gorm:"index:addr"` // 给address字段创建名为addr的索引
	IgnoreMe     int     `gorm:"-"` // 忽略本字段
}
//可以自定义设置表名  为ProUser或者其他  但其实
//func (User) TableName() string {
//	return "ProUser"
//}
func main() {
	//链接MySQL
	db, err := gorm.Open("mysql", "root:123456@(localhost)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	}
	//创建表 自动迁移（把结构体和数据表进行对应）
	//GORM 默认会使用名为ID的字段作为表的主键。或者在结构体里使用`gorm:"primary_key"`指定主键
	//表名默认就是结构体名称的复数// 将 User 的表名设置为 `profiles`

	db.AutoMigrate(&User{})
}
