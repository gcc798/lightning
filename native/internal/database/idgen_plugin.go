package database

import (
	"fmt"
	"reflect"

	"github.com/gcc798/microservice-kit/internal/utils/idgen"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IDGenPlugin GORM ID生成插件
type IDGenPlugin struct{}

// Name 执行业务逻辑。
func (p *IDGenPlugin) Name() string {
	return "idgen_plugin"
}

// Initialize 执行业务逻辑。
func (p *IDGenPlugin) Initialize(db *gorm.DB) error {
	return db.Callback().Create().Before("gorm:create").Register("idgen:before_create", p.beforeCreate)
}

// beforeCreate 在创建记录前自动生成ID
func (p *IDGenPlugin) beforeCreate(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}

	reflectValue := db.Statement.ReflectValue
	if reflectValue.Kind() == reflect.Slice {
		for i := 0; i < reflectValue.Len(); i++ {
			elem := reflectValue.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			p.processRecord(db, elem)
		}
	} else {
		p.processRecord(db, reflectValue)
	}
}

// processRecord 处理单条记录的ID生成
func (p *IDGenPlugin) processRecord(db *gorm.DB, reflectValue reflect.Value) {
	for _, field := range db.Statement.Schema.Fields {
		autogenTag := field.Tag.Get("autogen")
		if autogenTag == "" {
			continue
		}

		_, isZero := field.ValueOf(db.Statement.Context, reflectValue)
		if !isZero {
			continue
		}

		switch autogenTag {
		case "int64":
			if field.FieldType.Kind() != reflect.Int64 {
				fmt.Printf("[IDGen] Field %s has autogen:int64 but type is %v, skipping\n", field.Name, field.FieldType.Kind())
				continue
			}
			newID, err := idgen.NextID()
			if err != nil {
				db.AddError(fmt.Errorf("failed to generate int64 ID for field %s: %w", field.Name, err))
				return
			}
			if err := field.Set(db.Statement.Context, reflectValue, newID); err != nil {
				db.AddError(fmt.Errorf("failed to set int64 ID for field %s: %w", field.Name, err))
				return
			}

		case "string":
			if field.FieldType.Kind() != reflect.String {
				fmt.Printf("[IDGen] Field %s has autogen:string but type is %v, skipping\n", field.Name, field.FieldType.Kind())
				continue
			}
			newID := uuid.New().String()
			if err := field.Set(db.Statement.Context, reflectValue, newID); err != nil {
				db.AddError(fmt.Errorf("failed to set string ID for field %s: %w", field.Name, err))
				return
			}
		}
	}
}
