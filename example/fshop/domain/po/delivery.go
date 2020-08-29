// Code generated by 'freedom new-po'
package po

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Delivery struct {
	changes        map[string]interface{}
	Id             int       `gorm:"primary_key;column:id"`
	AdminId        int       `gorm:"column:admin_id"` // 管理员id
	OrderNo        string    `gorm:"column:order_no"`
	TrackingNumber string    `gorm:"column:tracking_number"` // 快递单号
	Created        time.Time `gorm:"column:created"`
	Updated        time.Time `gorm:"column:updated"`
}

func (obj *Delivery) TableName() string {
	return "delivery"
}

// TakeChanges .
func (obj *Delivery) TakeChanges() map[string]interface{} {
	if obj.changes == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range obj.changes {
		result[k] = v
	}
	obj.changes = nil
	return result
}

// updateChanges .
func (obj *Delivery) setChanges(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetAdminId .
func (obj *Delivery) SetAdminId(adminId int) {
	obj.AdminId = adminId
	obj.setChanges("admin_id", adminId)
}

// SetOrderNo .
func (obj *Delivery) SetOrderNo(orderNo string) {
	obj.OrderNo = orderNo
	obj.setChanges("order_no", orderNo)
}

// SetTrackingNumber .
func (obj *Delivery) SetTrackingNumber(trackingNumber string) {
	obj.TrackingNumber = trackingNumber
	obj.setChanges("tracking_number", trackingNumber)
}

// SetCreated .
func (obj *Delivery) SetCreated(created time.Time) {
	obj.Created = created
	obj.setChanges("created", created)
}

// SetUpdated .
func (obj *Delivery) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.setChanges("updated", updated)
}

// AddAdminId .
func (obj *Delivery) AddAdminId(adminId int) {
	obj.AdminId += adminId
	obj.setChanges("admin_id", gorm.Expr("admin_id + ?", adminId))
}
