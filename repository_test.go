package gorm_generics_test

import (
	"context"
	gorm_generics "github.com/ompluscator/gorm-generics"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	return gorm.Open(sqlite.Open("file:test?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
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
		ID:          8,
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

	_, err := repository.FindByID(ctx, 8)

	if err != nil {
		panic(err)
	}
}

func TestGormRepository_Count(t *testing.T) {
	db, _ := getDB()
	repository := gorm_generics.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	nb, err := repository.Count(ctx)

	if err != nil {
		panic(err)
	}
	if nb != 1 {
		panic("not good count")
	}
}

func TestGormRepository_DeleteByID(t *testing.T) {
	db, _ := getDB()
	repository := gorm_generics.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()
	err := repository.DeleteById(ctx, 8)
	if err != nil {
		panic(err)
	}
	_, err = repository.FindByID(ctx, 8)
	if err == nil {
		panic("supposed to be deleted")
	}
}

func TestGormRepository_Find(t *testing.T) {
	db, _ := getDB()
	repository := gorm_generics.NewRepository[ProductGorm, Product](db)
	ctx := context.Background()

	product := Product{
		ID:          1,
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	repository.Insert(ctx, &product)
	product2 := Product{
		ID:          2,
		Name:        "product2",
		Weight:      50,
		IsAvailable: true,
	}
	repository.Insert(ctx, &product2)
	many, err := repository.Find(ctx, gorm_generics.GreaterOrEqual("weight", 50))
	if err != nil {
		panic(err)
	}
	if len(many) != 2 {
		panic("should be 2")
	}

	repository.Insert(ctx, &Product{
		ID:          3,
		Name:        "product3",
		Weight:      250,
		IsAvailable: false,
	})

	many, err = repository.Find(ctx, gorm_generics.GreaterOrEqual("weight", 90))
	if err != nil {
		panic(err)
	}
	if len(many) != 2 {
		panic("should be 2")
	}

	many, err = repository.Find(ctx, gorm_generics.And(
		gorm_generics.GreaterOrEqual("weight", 90),
		gorm_generics.Equal("is_available", true)),
	)
	if err != nil {
		panic(err)
	}
	if len(many) != 1 {
		panic("should be 1")
	}
}

/*
TODO
Delete (by item)
Update
Find (with sql cond)
*/
