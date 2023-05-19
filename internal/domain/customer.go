package domain

import "time"

type Customer struct {
	ID          int       `gorm:"primarykey;autoIncrement:true" qsearch:"-"`
	Name        string    `gorm:"type:varchar(50);column:name" qsearch:"name"`
	Email       string    `gorm:"type:varchar(100);column:email" qsearch:"email"`
	MobilePhone string    `gorm:"type:varchar(14);column:mobile_phone" qsearch:"mobile_phone"`
	Password    string    `gorm:"type:varchar(100);column:password" qsearch:"-"`
	CreatedAt   time.Time `gorm:"column:created_at" qsearch:"-"`
	UpdatedAt   time.Time `gorm:"column:updated_at" qsearch:"-"`
}

// TableName name of table
func (r Customer) TableName() string {
	return "customers"
}
