package gorm_generics_test

import (
	"context"
	gorm_generics "github.com/philiphil/gorm-generics"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

// Product is a domain entity
type Product struct {
	ID          uint
	Name        string
	Weight      uint
	IsAvailable bool
}

// ProductGorm is DTO used to map Product entity to database
type ProductGorm struct {
	ID          uint   `gorm:"primaryKey;column:id"`
	Name        string `gorm:"column:name"`
	Weight      uint   `gorm:"column:weight"`
	IsAvailable bool   `gorm:"column:is_available"`
}

// ToEntity respects the gorm_generics.GormModel interface
// Creates new Entity from GORM model.
func (g ProductGorm) ToEntity() Product {
	return Product{
		ID:          g.ID,
		Name:        g.Name,
		Weight:      g.Weight,
		IsAvailable: g.IsAvailable,
	}
}

// FromEntity respects the gorm_generics.GormModel interface
// Creates new GORM model from Entity.
func (g ProductGorm) FromEntity(product Product) interface{} {
	return ProductGorm{
		ID:          product.ID,
		Name:        product.Name,
		Weight:      product.Weight,
		IsAvailable: product.IsAvailable,
	}
}

func getDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file:test?mode=memory&cache=shared&_fk=1"), &gorm.Config{})
}
func TestMain(m *testing.M) {
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(ProductGorm{})
	if err != nil {
		log.Fatal(err)
	}
	ret := m.Run()
	os.Exit(ret)
}
func TestGormRepository_Insert(t *testing.T) {
	db, _ := getDB()
	repository := gorm_generics.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	product := Product{
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err := repository.Insert(ctx, &product)
	if err != nil {
		panic(err)
	}
}

func TestGormRepository_FindByID(t *testing.T) {
	db, _ := getDB()
	repository := gorm_generics.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	product := Product{
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err := repository.Insert(ctx, &product)
	single, err := repository.FindByID(ctx, product.ID)
	if err != nil || single.ID != product.ID {
		panic(err)
	}
}
