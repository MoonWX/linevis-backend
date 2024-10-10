package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model                   // ID, CreatedAt, UpdatedAt, DeletedAt
	Name          string         `json:"name"`           // 商品名称
	MainBarcode   string         `json:"main_barcode"`   // 主条码
	ProductModel  string         `json:"model"`          // 商品型号
	Weight        int            `json:"weight"`         // 重量
	Specification string         `json:"specification"`  // 规格
	TargetAddress string         `json:"target_address"` // 目标地址
	Manual        string         `json:"manual"`         // 说明书
	SubMaterials  SubMaterialArr `json:"sub_materials"`  // 子物料，实际存储为 JSON 字符串
}

type SubMaterial struct {
	Name       string `json:"name"`
	SubBarcode string `json:"sub_barcode"`
}

type SubMaterialArr []SubMaterial

func (s SubMaterialArr) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *SubMaterialArr) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to scan Array value:", value))
	}
	if len(bytes) > 0 {
		return json.Unmarshal(bytes, s)
	}
	*s = make([]SubMaterial, 0)
	return nil
}

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	fmt.Println("InitDB")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})
	return db
}
