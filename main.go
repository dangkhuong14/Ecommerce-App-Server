package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// `id` int NOT NULL AUTO_INCREMENT,
// `category_id` int DEFAULT NULL,
// `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
// `image` json DEFAULT NULL,
// `type` enum('drink','food','topping') NOT NULL DEFAULT 'drink',
// `status` enum('activated','deactivated','out_of_stock') DEFAULT 'activated',
// `description` text,
// `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
// `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

type BaseModel struct {
	ID int `gorm:"column:id"`
	Status string `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type Product struct {
	BaseModel
	CategoryID int `gorm:"column:category_id"`
	Name string `gorm:"column:name"`
	Image *string `gorm:"column:image"`
	Type string `gorm:"column:type"`
	Description string `gorm:"column:description"`
}

type ProductUpdate struct {
	CategoryID *int `gorm:"column:category_id"`
	Name *string `gorm:"column:name"`
	Image *string `gorm:"column:image"`
	Type *string `gorm:"column:type"`
	Description *string `gorm:"column:description"`
	Status *string `gorm:"column:status"`
}

func (Product) TableName() string {
	return "products"
}

func main() {
	// Checking that an environment variable is present or not.
	mysqlConnStr, ok := os.LookupEnv("MYSQL_CONNECTION")

	if !ok {
		log.Fatalln("Missing MySQL connection string")
	}

	dsn := mysqlConnStr
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Cannot connect to MySQL:", err)
	}

	log.Println("Connected to MySQL:", db)

	// now:= time.Now().UTC()

	// Create new Product instance
	// newProduct := Product{
	// 	BaseModel: BaseModel{
	// 		Status: "activated",
	// 		CreatedAt: now,
	// 		UpdatedAt: now,
	// 	},
	// 	CategoryID: 1,
	// 	Name: "Americano",
	// 	Image: nil,
	// 	Type: "drink",
	// }
	
	// if err := db.Table(Product{}.TableName()).Create(&newProduct).Error; err != nil {
	// 	log.Println("Cannot create new Product:", err)
	// } else {
	// 	log.Println("Product created:", newProduct)
	// }

	// Read Product instance
	var existingProduct Product

	if err:=db.
	Table(Product{}.TableName()).
	Where("id= ?", 1).
	First(&existingProduct).Error; err !=nil {
		log.Println("Cannot read from Product:", err)
	} else {
		log.Println("Product read:", existingProduct)
	}

	//List products
	var products []Product

	if err:=db.
	Table(Product{}.TableName()).
	Where("status not in (?)", []string{"deactivated"}).
	Limit(3).
	Order("id desc").
	Offset(2).
	Find(&products).Error; err !=nil {
		log.Println("Cannot list from Product:", err)
	} else {
		log.Println("Products listed:", products)
	}

	// Update product 1: use existing product
	existingProduct.Name = "Latte"

	if err:=db.
	Table(Product{}.TableName()).
	Updates(existingProduct).Error; err !=nil {
		log.Println("Cannot update existing Product:", err)
	} else {
		log.Println("Products updated:", existingProduct)
	}

	// Update product 2: use where condition
	if err:=db.Table(Product{}.TableName()).
	Where("id = ?", 2).Updates(map [string]interface{}{"name": "Espresso"}).
	Error; err !=nil {
		log.Println("Cannot update Product with ID 2:", err)
	} else {
		log.Println("Products ID 2 updated:")
	}
	
	// Update product 3: empty string
	emptyString := ""
	if err := db.
	Table(Product{}.TableName()).
	Where("id = ?", 3).
	Updates(ProductUpdate{Name: &emptyString}).Error; err !=nil {
		log.Println("Cannot update Product with ID 3:", err)
	} else {
		log.Println("Products ID 3 updated:")
	}	
}