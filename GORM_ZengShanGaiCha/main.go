package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //加个下划线是代表没有直接用到他
)

//https://www.liwenzhou.com/posts/Go/gorm-crud/
//GORM增删改查
type User struct {
	gorm.Model
	Name string
	Age int64
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
	db.AutoMigrate(&User{})
	//创建数据行
	u1 := User{Name:"zhiqing", Age:26}
	u2 := User{Name:"xining", Age:25}
	db.Create(&u1)
	db.Create(&u2)

	//一般查询
	var u User
	// 查询所有的记录
	db.Find(&u)
	fmt.Printf("查询所有的记录u:%v\n", u)
	// 查询指定的某条记录(仅当主键为整型时可用)
	db.First(&u, 0)
	fmt.Printf("查询指定的某条记录u:%v\n", u)
	//这里是查询表中第一条数据保存到u中，后面会介绍复杂的查询方法
	db.First(&u)
	fmt.Printf("第一条数据u:%v\n", u)
	// 随机获取一条记录
	db.Take(&u)
	fmt.Printf("随机获取一条记录u:%v\n", u)
	// 根据主键查询最后一条记录
	db.Last(&u)
	fmt.Printf("最后一条记录u:%v\n", u)

	//where条件查询
	//按条件查询返回第一条记录
	db.Where("name = ?", "zhiqing").First(&u)
	//按条件查询返回所有记录
	db.Where("name = ?", "zhiqing").Find(&u)
	// <> 不等于
	db.Where("age <> ?", 30).Find(&u)
	//范围条件
	db.Where("name IN (?)", []string{"zhiqing", "xining"}).Find(&u)
	// 模糊查询
	db.Where("name LIKE ?", "%zhi%").Find(&u)
	// AND
	db.Where("name = ? AND age >= ?", "zhiqing", 26).Find(&u)
	// Time
	//db.Where("updated_at > ?", lastWeek).Find(&u)
	// BETWEEN
	//db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&u)

	//Struct & Map查询
	// Struct
	db.Where(&User{Name: "zhiqing", Age: 26}).First(&u)
	// Map
	db.Where(map[string]interface{}{"name": "zhiqing", "age": 26}).Find(&u)
	// 主键的切片
	db.Where([]int64{20, 21, 22}).Find(&u)

	//Not条件查询
	db.Not("name", "zhiqing").First(&u)
	// Not In
	db.Not("name", []string{"zhiqing", "zhiqing 2"}).Find(&u)
	// Not In slice of primary keys
	db.Not([]int64{1,2,3}).First(&u)
	db.Not([]int64{}).First(&u)
	// Plain SQL
	db.Not("name = ?", "zhiqing").First(&u)
	// Struct
	db.Not(User{Name: "zhiqing"}).First(&u)

	//Or条件
	db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&u)
	// Struct
	db.Where("name = 'zhiqing'").Or(User{Name: "zhiqing 2"}).Find(&u)
	// Map
	db.Where("name = 'zhiqing'").Or(map[string]interface{}{"name": "zhiqing 2"}).Find(&u)

	//高级查询
	//子查询
	db.Where("amount > ?", db.Table("orders").Select("AVG(amount)").Where("state = ?", "paid").SubQuery()).Find(&u)
	// SELECT * FROM "orders"  WHERE "orders"."deleted_at" IS NULL AND (amount > (SELECT AVG(amount) FROM "orders"  WHERE (state = 'paid')));
	//选择字段
	//Select，指定你想从数据库中检索出的字段，默认会选择全部字段。
	db.Select("name, age").Find(&u)
	//// SELECT name, age FROM users;
	db.Select([]string{"name", "age"}).Find(&u)
	//// SELECT name, age FROM users;
	db.Table("users").Select("COALESCE(age,?)", 42).Rows()
	//// SELECT COALESCE(age,'42') FROM users;

	//更新
	//更新所有字段
	//Save()默认会更新该对象的所有字段，即使你没有赋值。
	db.First(&u)
	u.Name = "七米"
	u.Age = 99
	db.Save(&u)
	////  UPDATE `users` SET `created_at` = '2020-02-16 12:52:20', `updated_at` = '2020-02-16 12:54:55', `deleted_at` = NULL, `name` = '七米', `age` = 99, `active` = true  WHERE `users`.`deleted_at` IS NULL AND `users`.`id` = 1
	//更新修改字段
	//如果你只希望更新指定字段，可以使用Update或者Updates
	// 更新单个属性，如果它有变化
	db.Model(&u).Update("name", "hello")
	//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;
	// 根据给定的条件更新单个属性
	db.Model(&u).Where("active = ?", true).Update("name", "hello")
	//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
	// 使用 map 更新多个属性，只会更新其中有变化的属性
	db.Model(&u).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	//// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
	// 使用 struct 更新多个属性，只会更新其中有变化且为非零值的字段
	db.Model(&u).Updates(User{Name: "hello", Age: 18})
	//// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;
	// 警告：当使用 struct 更新时，GORM只会更新那些非零值的字段
	// 对于下面的操作，不会发生任何更新，"", 0, false 都是其类型的零值
	db.Model(&u).Updates(User{Name: "", Age: 0, Active: false})
	//更新选定字段
	//如果你想更新或忽略某些字段，你可以使用 Select，Omit
	db.Model(&u).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;
	db.Model(&u).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
	//// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

	//删除记录
	//警告 删除记录时，请确保主键字段有值，GORM 会通过主键去删除记录，如果主键为空，GORM 会删除该 model 的所有记录。
	// 删除现有记录
	db.Delete(&u)
	//// DELETE from emails where id=10;
	// 为删除 SQL 添加额外的 SQL 操作
	db.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Delete(&u)
	//// DELETE from emails where id=10 OPTION (OPTIMIZE FOR UNKNOWN);
	//批量删除
	//删除全部匹配的记录
	db.Where("email LIKE ?", "%jinzhu%").Delete(User{})
	//// DELETE from emails where email LIKE "%jinzhu%";
	db.Delete(User{}, "email LIKE ?", "%jinzhu%")
	//// DELETE from emails where email LIKE "%jinzhu%";
	//软删除
	//如果一个 model 有 DeletedAt 字段，他将自动获得软删除的功能！ 当调用 Delete 方法时， 记录不会真正的从数据库中被删除， 只会将DeletedAt 字段的值会被设置为当前时间
	db.Delete(&u)
	//// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;
	// 批量删除
	db.Where("age = ?", 20).Delete(&User{})
	//// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;
	// 查询记录时会忽略被软删除的记录
	db.Where("age = 20").Find(&u)
	//// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
	// Unscoped 方法可以查询被软删除的记录
	db.Unscoped().Where("age = 20").Find(&u)
	//// SELECT * FROM users WHERE age = 20;
	//物理删除
	// Unscoped 方法可以物理删除记录
	db.Unscoped().Delete(&u)
	//// DELETE FROM orders WHERE id=10;

}

