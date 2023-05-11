package main

//什么是ORM，对象关系映射  关系数据库MySQL
//数据表对应结构体，数据行对应结构体实例，字段对应结构体字段

//连接不同的数据库都需要导入对应数据的驱动程序，
//GORM已经贴心的为我们包装了一些驱动程序，只需要按如下方式导入需要的数据库驱动即可：
import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"//加个下划线是代表没有直接用到他
)
// import _ "github.com/jinzhu/gorm/dialects/postgres"
// import _ "github.com/jinzhu/gorm/dialects/sqlite"
// import _ "github.com/jinzhu/gorm/dialects/mssql"

type UserInfo struct {
	ID uint
	Name string
	Gender string
	Hobby string
}

func main() {
	//链接MySQL
	db, err := gorm.Open("mysql", "root:123456@(localhost)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	}
	//一定记得要关闭数据库 defer延时关闭即可
	defer db.Close()

	//创建表 自动迁移（把结构体和数据表进行对应）
	db.AutoMigrate(&UserInfo{})

	//创建数据行
	u1 := UserInfo{519, "zhiqing", "女", "追剧"}
	u2 := UserInfo{1123, "xining", "男", "游戏"}
	db.Create(&u1)
	db.Create(&u2)
	//查询
	var u UserInfo
	db.First(&u)   //这里是查询表中第一条数据保存到u中，后面会介绍复杂的查询方法
	fmt.Printf("u:%v\n", u)
	//条件查找
	var uu UserInfo
	db.Find(&uu, "hobby=?", "游戏")
	fmt.Printf("u:%#v\n", uu)
	//更新
	db.Model(&u).Update("hobby", "看赵丽颖的剧")
	//删除
	db.Delete(&u)
}
