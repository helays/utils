package userDb

import (
	"errors"
	"fmt"
	"github.com/helays/utils/db/dbErrors/errTools"
	"gorm.io/gorm"
	"regexp"
)

type Model interface {
	Valid(tx *gorm.DB) error // 验证
}

// Create 通用创建函数（使用泛型）
func Create[T Model](tx *gorm.DB, src T) error {
	// 调用 Valid 方法进行验证
	if err := src.Valid(tx); err != nil {
		return fmt.Errorf("验证失败：%s", err.Error())
	}
	// 使用 GORM 的 Create 方法插入数据
	if err := tx.Session(&gorm.Session{}).Create(src).Error; err != nil {
		return fmt.Errorf("插入数据失败：%s", errTools.Error(err).Error())
	}

	return nil
}

type QueryOpt struct {
	Query string
	Args  []any
}

// Curd 通用curd函数，配置结构
type Curd struct {
	Table         string
	Alias         string
	Select        QueryOpt
	Where         QueryOpt
	Preload       []QueryOpt
	Omit          []string
	Pk            string                    // 主键字段 id row_Id
	MustField     map[string]*regexp.Regexp // 必填字段，正则校验
	ValidExist    bool                      // 存在校验,true 需要校验目标是否存在，false 忽略校验
	Update        []any                     // 更新
	Updates       any                       // 更新
	EnableDefault bool
	Joins         []QueryOpt
}

// Update 通用更新函数，使用泛型
func Update[T Model](tx *gorm.DB, src T, c Curd) error {
	// 调用 Valid 方法进行验证
	err := src.Valid(tx.Session(&gorm.Session{}))
	if err != nil {
		return fmt.Errorf("验证失败：%s", err.Error())
	}
	if c.ValidExist {
		// 搜索数据是否存在
		if _, err = FindOne[T](tx, Curd{
			Where: c.Where,
		}); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("数据不存在")
			}
			return err
		}
	}

	_tx := tx.Session(&gorm.Session{}).Where(c.Where.Query, c.Where.Args...)
	if c.Omit != nil && len(c.Omit) > 0 {
		_tx.Omit(c.Omit...)
	}
	if c.Select.Query != "" {
		_tx.Select(c.Select.Query, c.Select.Args...)
	}
	// 使用 GORM 的 Update 方法更新数据
	if err = _tx.Save(src).Error; err != nil {
		return fmt.Errorf("更新数据失败：%w", errTools.Error(err))
	}
	return nil
}

// UpdateWithoutValid 通过更新状态字段，实现通用软删除
func UpdateWithoutValid[T any](tx *gorm.DB, opt Curd) (err error) {
	_tx := tx.Session(&gorm.Session{}).Model(new(T)).Where(opt.Where.Query, opt.Where.Args...)
	if len(opt.Omit) > 0 {
		_tx.Omit(opt.Omit...)
	}
	if opt.Select.Query != "" {
		_tx.Select(opt.Select.Query, opt.Select.Args...)
	}
	if len(opt.Update) == 2 {
		err = _tx.Update(opt.Update[0].(string), opt.Update[1]).Error
	} else {
		err = _tx.Updates(opt.Updates).Error
	}
	if err != nil {
		return errTools.Error(err)
	}
	return
}

// FindOne 查询某个对象
func FindOne[T any](tx *gorm.DB, opts Curd) (T, error) {
	_tx := tx.Session(&gorm.Session{}).Where(opts.Where.Query, opts.Where.Args...)
	if opts.Table != "" {
		_tx.Table(opts.Table)
	}
	if opts.Select.Query != "" {
		_tx.Select(opts.Select.Query, opts.Select.Args...)
	}
	if len(opts.Omit) > 0 {
		_tx.Omit(opts.Omit...)
	}
	for _, item := range opts.Preload {
		_tx.Preload(item.Query)
	}
	var data T
	err := _tx.Take(&data).Error
	if err != nil {
		err = errTools.Error(err)
	}
	return data, err
}

// Query 查询多个对象
func Query[T any](tx *gorm.DB, opts Curd) (res []T, err error) {
	_tx := tx.Session(&gorm.Session{}).Where(opts.Where.Query, opts.Where.Args...)
	if opts.Table != "" {
		_tx.Table(opts.Table)
	}
	if opts.Select.Query != "" {
		_tx.Select(opts.Select.Query, opts.Select.Args...)
	}
	if len(opts.Omit) > 0 {
		_tx.Omit(opts.Omit...)
	}
	for _, item := range opts.Preload {
		_tx.Preload(item.Query)
	}
	err = _tx.Find(&res).Error
	if err != nil {
		err = errTools.Error(err)
	}
	return
}
