package schema

import (
	"reflect"
	"time"

	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
)

// NOTE: We need this function to initialize the ID field of tables with CUID
func RegisterCUIDCallback(db *gorm.DB) {
	db.Callback().Create().Before("gorm:before_create").Register("set_cuid_if_empty", func(tx *gorm.DB) {
		if tx.Statement.Schema != nil {
			idField := tx.Statement.Schema.LookUpField("ID")
			if idField != nil && idField.FieldType.Kind() == reflect.String {
				val, _ := idField.ValueOf(tx.Statement.Statement.Context, tx.Statement.ReflectValue)
				if str, ok := val.(string); ok && str == "" {
					tx.Statement.SetColumn("ID", cuid2.Generate())
				}
			}
		}
	})
}

type User struct {
	ID           string `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time `gorm:"index"`
	Email        string
	Name         string
	PasswordHash string
}
