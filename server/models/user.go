package models

import "time"

type User struct {
	ID         int       `gorm:"type:bigint; NOT NULL"`
	Username   string    `valid:"Required;MinSize(4);MaxSize(12)"`
	Phone      string    `valid:"Mobile;Required"`
	Email      string    `valid:"Email;MaxSize(50)"`
	Password   string    `valid:"Required;MinSize(4);MaxSize(6)"`
	RePassword string    `gorm:"-" valid:"Required;MinSize(4);MaxSize(6);eqfield=Password"`
	RegTime    time.Time `gorm:"type:datetime; DEFAULT: CURRENT_TIMESTAMP"`
}
