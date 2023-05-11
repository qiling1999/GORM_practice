package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"//加个下划线是代表没有直接用到他
)

//定义结构体模型
type User struct {
	ID int64
	Name string `gorm:"default:'zhiqing'"`
	Age int64
}
//通过tag定义字段的默认值，在创建记录时候生成的 SQL 语句会排除没有值或值为 零值 的字段。
//在将记录插入到数据库后，Gorm会从数据库加载那些字段的默认值
//注意：所有字段的零值, 比如0, "",false或者其它零值，都不会保存到数据库内，但会使用他们的默认值。
//如果你想避免这种情况，可以考虑使用指针或实现 Scanner/Valuer接口，比如：
//Name *string `gorm:"default:'zhiqing'"`
//u := User{ Name:new(string), Age:26} //空值传入数据库方式
func main() {
	//链接MySQL
	db, err := gorm.Open("mysql", "root:123456@(localhost)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	}
	defer db.Close()

	//创建表 自动迁移（把结构体和数据表进行对应）
	db.AutoMigrate(&User{})

	//创建数据行，
	//要创建记录，要先创建结构体实例
	//u := User{ Name:new(string), Age:26} //空值传入数据库方式
	u := User{ Name:"", Age:26}
	fmt.Println(db.NewRecord(&u))//判断主键是否为空 主键为ID
	db.Create(&u) //这里写u也是可以的
	fmt.Println(db.NewRecord(&u))//判断主键是否为空
}
