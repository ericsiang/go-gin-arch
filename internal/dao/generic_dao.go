package dao

import (
	"fmt"
	"math"

	"gorm.io/gorm"
)

const (
	// PerPageCount 每頁數量
	PerPageCount = 15 // 每頁數量
)

// GenericDaoInterface 通用資料存取介面
type GenericDaoInterface[T any] interface {
	GetDB() *gorm.DB
	Create(model *T) (*T, error)
	GetByID(id uint64) (*T, error)
	GetAll() ([]T, error)
	// Update(model *T) error
	Delete(model *T) error
	Paginate(db *gorm.DB, result []T, option PaginateOption) (*Paginator[T], error)
	SimplePaginate(db *gorm.DB, result []T, option PaginateOption) (*Paginator[T], error)
}

// Paginator 分頁結構
type Paginator[T any] struct {
	CurrentPage  int    `json:"current_page"`   // 當前頁碼
	PerPageCount int    `json:"per_page_count"` // 每頁數量
	Total        *int64 `json:"total"`          // 總數量
	LastPage     *int   `json:"last_page"`      // 最後頁數
	Data         []T    `json:"data"`           // 當前頁數據
	HasMore      bool   `json:"has_more"`       // 是否有更多頁數
}

// PaginateOption 分頁選項
type PaginateOption struct {
	Page         int
	PerPageCount int
	OrderBy      string
	GroupBy      string
}

// GenericDAO 通用資料存取實現
type GenericDAO[T any] struct {
	db *gorm.DB
}

// NewGenericDAO 建立新的通用資料存取物件
func NewGenericDAO[T any](db *gorm.DB) GenericDaoInterface[T] {
	return &GenericDAO[T]{
		db: db,
	}
}

// GetDB 返回通用資料存取的 GORM 資料庫連接
func (g *GenericDAO[T]) GetDB() *gorm.DB {
	return g.db
}

func (g *GenericDAO[T]) Create(model *T) (*T, error) {
	err := g.db.Create(&model).Error
	return model, err
}

func (g *GenericDAO[T]) GetByID(id uint64) (*T, error) {
	var model T
	err := g.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (g *GenericDAO[T]) GetAll() ([]T, error) {
	var models []T
	err := g.db.Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// func (g *GenericDAO[T]) Update(model *T) error {
// 	return g.db.Save(model).Error
// }

func (g *GenericDAO[T]) Delete(model *T) error {
	return g.db.Delete(model).Error
}

// Paginate 分頁，COUNT 總量和最後一頁，當數據量大時，COUNT 查詢可能會很慢
func (g *GenericDAO[T]) Paginate(db *gorm.DB, result []T, option PaginateOption) (*Paginator[T], error) {
	var total int64
	var model T

	// 如果沒有指定頁碼，預設第一頁
	if option.Page < 1 {
		option.Page = 1
	}

	// 如果沒有指定每頁數量，預設 15
	if option.PerPageCount < 1 {
		option.PerPageCount = PerPageCount
	}

	// 排序
	if option.OrderBy != "" {
		db = db.Order(option.OrderBy)
	}

	if option.GroupBy != "" {
		db = db.Group(option.GroupBy)
	}

	// 計算總數
	err := db.Model(&model).Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("GenericDAO Paginate() Count() : %w", err)
	}

	// 計算偏移量
	offset := (option.Page - 1) * option.PerPageCount

	// 查詢當前頁數據
	err = db.Offset(offset).Limit(option.PerPageCount).Find(&result).Error
	if err != nil {
		logData := map[string]interface{}{
			"Offset": offset,
			"Limit":  option.PerPageCount,
			"Order":  option.OrderBy,
			"Group":  option.GroupBy,
		}
		return nil, fmt.Errorf("GenericDAO Paginate() Find() data : %+v : %w", logData, err)
	}

	// 計算最後一頁
	lastPage := int(math.Ceil(float64(total) / float64(option.PerPageCount)))

	// 判斷是否還有更多數據
	hasMore := option.Page < lastPage

	return &Paginator[T]{
		CurrentPage:  option.Page,
		PerPageCount: option.PerPageCount,
		Total:        &total,
		LastPage:     &lastPage,
		Data:         result,
		HasMore:      hasMore,
	}, nil
}

// SimplePaginate 簡單分頁，用於大數據量場景
// 這個方法不會計算總數量，只會返回當前頁的數據和是否有更多頁數
func (g *GenericDAO[T]) SimplePaginate(db *gorm.DB, result []T, option PaginateOption) (*Paginator[T], error) {

	// 基本設置
	if option.Page < 1 {
		option.Page = 1
	}
	if option.PerPageCount < 1 {
		option.PerPageCount = PerPageCount
	}

	// 排序
	if option.OrderBy != "" {
		db = db.Order(option.OrderBy)
	}

	if option.GroupBy != "" {
		db = db.Group(option.GroupBy)
	}

	// 計算偏移量
	offset := (option.Page - 1) * option.PerPageCount

	// 多額外查詢一筆資料來判斷是否有下一頁
	err := db.Offset(offset).Limit(option.PerPageCount + 1).Find(&result).Error
	if err != nil {
		logData := map[string]interface{}{
			"Offset": offset,
			"Limit":  option.PerPageCount + 1,
			"Order":  option.OrderBy,
			"Group":  option.GroupBy,
		}
		return nil, fmt.Errorf("GenericDAO SimplePaginate() Find() data : %+v : %w", logData, err)
	}

	// 判斷是否有更多數據，如果查詢的數據量大於每頁數量，則有更多數據
	hasMore := len(result) > option.PerPageCount
	if hasMore { // 有更多數據，去掉最後一筆數據
		result = result[:option.PerPageCount]
	}

	return &Paginator[T]{
		CurrentPage:  option.Page,
		PerPageCount: option.PerPageCount,
		Data:         result,
		HasMore:      hasMore,
	}, nil
}
