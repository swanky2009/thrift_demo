package models

import (
	"errors"
	"fmt"
	//	"strconv"
	//	"time"
	//	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var o orm.Ormer

func init() {

	// 需要在init中注册定义的model
	orm.RegisterModel(new(User), new(Profile))

	o = orm.NewOrm()

	//createTable()
}

//自动建表
func createTable() {
	name := "default"                          //数据库别名
	force := false                             //不强制建数据库
	verbose := true                            //打印建表过程
	err := orm.RunSyncdb(name, force, verbose) //建表
	if err != nil {
		//fmt.P(err)
	}
}

type User struct {
	Uid     int `orm:"pk;auto"`
	Name    string
	Profile *Profile `orm:"rel(one)"` // OneToOne relation
}

type Profile struct {
	Uid int `orm:"pk;auto"`
	Age int16
}

func AddUser(name string, age int16) int64 {
	profile := new(Profile)
	profile.Age = age

	user := new(User)
	user.Profile = profile
	user.Name = name

	id, err := o.Insert(profile)
	if err != nil {
		fmt.Println(err)
	}
	id, err = o.Insert(user)
	if err == nil {
		return id
	}
	fmt.Println(err)
	return 0
}

func GetUser(uid int) (u *User, err error) {
	user := &User{Uid: uid}
	err = o.Read(user)
	if err == nil {
		if user.Profile != nil {
			o.Read(user.Profile)
		}
		return user, nil
	}
	return nil, errors.New("User not exists")
}

func GetAllUsers(rows int, page int) []*User {
	var users []*User
	_, err := o.QueryTable("User").RelatedSel().Offset((page - 1) * rows).Limit(rows).All(&users)
	if err == nil {
		return users
	}
	return nil
}

func UpdateUser(uid int, name string, age int16) (a *User, err error) {
	user := User{Uid: uid}
	if o.Read(&user) == nil {
		user.Name = name
		//fmt.Println(user.Profile.Uid)
		p, _ := o.Raw("UPDATE profile SET age = ? WHERE Uid = ?").Prepare()
		p.Exec(age, user.Profile.Uid)
		p.Close() // 别忘记关闭 statement

		if _, err := o.Update(&user); err == nil {
			o.Read(user.Profile)
			return &user, nil
		}
	}
	return nil, errors.New("User Not Exist")
}

func DeleteUser(uid int) int64 {
	user := User{Uid: uid}
	if num, err := o.Delete(&user); err == nil {
		return num
	}
	return 0
}
