package main

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Brand struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Products []Product `gorm:"foreignKey:BrandID"`
}

type Product struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `json:"name"`
	BrandID uint   `json:"brand_id"`
	Brand   Brand
	Tags    []Tag `gorm:"many2many:product_tags;"`
}

type Tag struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Products []Product `gorm:"many2many:product_tags;"`
}

var db *gorm.DB

func InitDB() {
	var err error

	dsn := "host=localhost user=postgres password=؟ dbname=GoTest port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect to the database: %v", err)
	}

	err = db.AutoMigrate(&Brand{}, &Product{}, &Tag{})

	if err != nil {
		fmt.Println("Faild To Migrate. ", err)
	}
}

func main() {
	InitDB()
	app := iris.New()

	app.Get("/products", getAllProducts)
	app.Post("/products", createProduct)
	app.Patch("/products/{id:int}", UpdateProduct)
	app.Delete("/products/{id:int}", DeleteProduct)
	app.Get("/brands", getALlBrands)
	app.Post("/brands", createBrand)
	app.Get("/tags", getAllTags)
	app.Post("/tags", createTag)

	app.Listen(":8080")
}

// products

func createProduct(ctx iris.Context) {
	var product Product

	if err := ctx.ReadJSON(&product); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	if err := db.Create(&product).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(product)
}

func getAllProducts(ctx iris.Context) {
	var products []Product
	db.Preload("Brand").Preload("Tags").Find(&products)
	ctx.JSON(products)
}

func UpdateProduct(ctx iris.Context) {
	var product Product
	productID, _ := ctx.Params().GetInt("id")

	if err := db.First(&product, productID).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Product not found"})
		return
	}

	if err := ctx.ReadJSON(&product); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid input"})
		return
	}

	if err := db.Save(&product).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Error updating product"})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(product)
}

func DeleteProduct(ctx iris.Context) {
	var product Product
	productID, _ := ctx.Params().GetInt("id")

	if err := db.First(&product, productID).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Product not found"})
		return
	}

	if err := db.Delete(&product).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Error updating product"})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Product deleted successfully"})
}

// brand

func createBrand(ctx iris.Context) {
	var brand Brand

	if err := ctx.ReadJSON(&brand); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	if err := db.Create(&brand).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(brand)
}

func getALlBrands(ctx iris.Context) {
	var brands []Brand
	db.Find(&brands)
	ctx.JSON(brands)
}

// tag
func createTag(ctx iris.Context) {
	var tag Tag

	if err := ctx.ReadJSON(&tag); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	if err := db.Create(&tag).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(tag)
}

func getAllTags(ctx iris.Context) {
	var tags []Tag
	db.Find(&tags)
	ctx.JSON(tags)
}
