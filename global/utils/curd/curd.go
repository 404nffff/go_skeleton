package curd

import (
	"fmt"
	"tool/pkg/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 配置
type Config[T any] struct {
	db    *gorm.DB // 数据库连接
	model T        // 模型
	Conn  string   // 数据库名称
}

// 初始化
func New[T any](model T) *Config[T] {

	Config := &Config[T]{model: model, Conn: "Local"}

	//链接数据库
	db := mysql.NewClient(Config.Conn)

	Config.db = db.Model(model)

	return Config
}

// 设置连接数据库
func (b *Config[T]) SetDb(Conn string) *Config[T] {

	//链接数据库
	db := mysql.NewClient(Conn)

	b.db = db.Model(b.model)
	return b
}

// 设置查询条件
func (b *Config[T]) Where(query interface{}, args ...interface{}) *Config[T] {
	b.db.Where(query, args...)
	return b
}

// 重置where条件
func (b *Config[T]) ResetWhere() *Config[T] {

	b.db.Statement.Clauses = map[string]clause.Clause{}

	return b
}

// 查询单条数据
func (b *Config[T]) First() (T, error) {

	var data T
	err := b.db.First(&data).Error
	return data, err
}

// 查询多条数据
func (b *Config[T]) Find() (T, error) {

	var data T
	err := b.db.Find(&data).Error
	return data, err
}

// 分页数据
func (b *Config[T]) PageList(page int, pageSize int) ([]T, error) {

	defer b.ResetWhere()

	var list []T
	if err := b.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil

}

// 统计
func (b *Config[T]) Count() (int64, error) {

	defer b.ResetWhere()

	var count int64
	db := b.db.Model(b.model)

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// 创建
func (b *Config[T]) Add() (int64, error) {

	defer b.ResetWhere()

	err := b.db.Create(b.model).Error

	if err != nil {
		return 0, err
	}

	//获取主键名称
	primaryKey := b.db.Statement.Schema.PrioritizedPrimaryField.Name

	//获取主键值
	primaryKeyValue := b.db.Statement.ReflectValue.FieldByName(primaryKey).Int()

	return primaryKeyValue, nil
}

// 修改
func (b *Config[T]) Update(data T) error {

	defer b.ResetWhere()

	return b.db.Updates(data).Error
}

// 新增或修改
func (b *Config[T]) Save(data T) (int16, error) {

	defer b.ResetWhere()

	// //判断是否有where
	if _, ok := b.db.Statement.Clauses["WHERE"]; ok {

		//判断数据是否存在
		_, err := b.First()
		if err != nil {
			return 0, err
		}

		//获取主键名称
		primaryKey := b.db.Statement.Schema.PrioritizedPrimaryField.Name

		//获取当前主键值
		primaryKeyValue := b.db.Statement.ReflectValue.FieldByName(primaryKey).Int()

		//如果是0 则没有数据
		if primaryKeyValue == 0 {
			return 0, fmt.Errorf("数据不存在")
		}

		//修改
		err = b.Update(data)

		if err != nil {
			return 0, err
		}

		return int16(primaryKeyValue), nil

	}

	//新增
	err := b.db.Create(&data).Error

	if err != nil {
		return 0, err
	}

	// 获取主键值
	primaryKey := b.db.Statement.Schema.PrioritizedPrimaryField.Name

	//获取主键值
	primaryKeyValue := b.db.Statement.ReflectValue.FieldByName(primaryKey).Int()

	return int16(primaryKeyValue), nil

}

// 删除
func (b *Config[T]) Delete() error {

	defer b.ResetWhere()

	//判断是否有where
	if _, ok := b.db.Statement.Clauses["WHERE"]; !ok {
		return fmt.Errorf("删除条件不能为空")
	}

	return b.db.Delete(b.model).Error
}
