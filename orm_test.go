package gorm

import "testing"

type User struct {
	Id   int64
	Name string
}

var db DB

func init() {
	db, _ = Open("postgres", "user=gorm dbname=gorm sslmode=disable")
	db.Exec("drop table users;")
}

func TestCreateTable(t *testing.T) {
	orm := db.CreateTable(&User{})
	if orm.Error != nil {
		t.Errorf("No error should raise when create table, but got %+v", orm.Error)
	}
}

func TestSaveAndFind(t *testing.T) {
	name := "save_and_find"
	u := &User{Name: name}
	db.Save(u)
	if u.Id == 0 {
		t.Errorf("Should have ID after create record")
	}

	user := &User{}
	db.First(user)
	if user.Name != name {
		t.Errorf("User should be saved and fetched correctly")
	}

	users := []User{}
	db.Find(&users)
}

func TestUpdate(t *testing.T) {
	name := "update"
	user := User{Name: name}
	db.Save(&user)

	user_id := user.Id
	if user_id == 0 {
		t.Errorf("User Id should exist after create")
	}

	orm := db.Where("name = ?", "update").First(&User{})
	if orm.Error != nil {
		t.Errorf("No error should raise when looking for a exiting user")
	}

	user.Name = "update2"
	db.Save(&user)
	orm = db.Where("name = ?", "update").First(&User{})
	if orm.Error == nil {
		t.Errorf("Should raise error when looking for a existing user with an outdated name")
	}

	orm = db.Where("name = ?", "update2").First(&User{})
	if orm.Error != nil {
		t.Errorf("Shouldn't raise error when looking for a existing user with the new name")
	}
}

func TestDelete(t *testing.T) {
	name, name2 := "delete", "delete2"
	user := User{Name: name}
	db.Save(&user)
	db.Save(&User{Name: name2})
	orm := db.Delete(&user)

	orm = db.Where("name = ?", name).First(&User{})
	if orm.Error == nil {
		t.Errorf("User should be deleted successfully")
	}

	orm = db.Where("name = ?", name2).First(&User{})
	if orm.Error != nil {
		t.Errorf("User2 should not be deleted")
	}
}

func TestWhere(t *testing.T) {
	name := "where"
	db.Save(&User{Name: name})

	user := &User{}
	db.Where("Name = ?", name).First(user)
	if user.Name != name {
		t.Errorf("Should found out user with name '%v'", name)
	}

	user = &User{}
	orm := db.Where("Name = ?", "noexisting-user").First(user)
	if orm.Error == nil {
		t.Errorf("Should return error when looking for none existing record, %+v", user)
	}

	users := []User{}
	orm = db.Where("Name = ?", "none-noexisting").Find(&users)
	if orm.Error != nil {
		t.Errorf("Shouldn't return error when looking for none existing records, %+v", users)
	}
	if len(users) != 0 {
		t.Errorf("Shouldn't find anything when looking for none existing records, %+v", users)
	}
}