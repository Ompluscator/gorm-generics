# PoC for Go generics with GORM

## Introduction
This repository represents a small PoC for using 
[Go generics](https://levelup.gitconnected.com/generics-in-go-viva-la-revolution-e27898bf5495) together 
with [GORM](https://gorm.io/index.html), an Object-relational 
mapping library for Golang.

At this stage it emphasizes possibilities, and it is not stable implementation.
In this stage, tt is not meant to be used for production system.

Future development is the intention for this project,
and any contribution is more than welcome.

## Example
```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ompluscator/gorm-generics"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
func (g ProductGorm) ToEntity() Product {
	return Product{
		ID:          g.ID,
		Name:        g.Name,
		Weight:      g.Weight,
		IsAvailable: g.IsAvailable,
	}
}

// NewProductGorm respects the factory method gorm_generics.ModelFactoryMethod
func NewProductGorm(product Product) ProductGorm {
	return ProductGorm{
		ID:          product.ID,
		Name:        product.Name,
		Weight:      product.Weight,
		IsAvailable: product.IsAvailable,
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("file:test?mode=memory&cache=shared&_fk=1"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(ProductGorm{})
	if err != nil {
		log.Fatal(err)
	}

	repository := gorm_generics.NewRepository(db, NewProductGorm)

	ctx := context.Background()

	product := Product{
		Name:        "product1",
		Weight:      100,
		IsAvailable: true,
	}
	err = repository.Insert(ctx, &product)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(product)
	// Out:
	// {1 product1 100 true}

	single, err := repository.FindByID(ctx, product.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(single)
	// Out:
	// {1 product1 100 true}

	err = repository.Insert(ctx, &Product{
		Name:        "product2",
		Weight:      50,
		IsAvailable: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	many, err := repository.Find(ctx, gorm_generics.GreaterOrEqual("weight", 50))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(many)
	// Out:
	// [{1 product1 100 true} {2 product2 50 true}]

	err = repository.Insert(ctx, &Product{
		Name:        "product3",
		Weight:      250,
		IsAvailable: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	many, err = repository.Find(ctx, gorm_generics.GreaterOrEqual("weight", 90))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(many)
	// Out:
	// [{1 product1 100 true} {3 product3 250 false}]

	many, err = repository.Find(ctx, gorm_generics.And(
		gorm_generics.GreaterOrEqual("weight", 90),
		gorm_generics.Equal("is_available", true)),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(many)
	// Out:
	// [{1 product1 100 true}]
}
```