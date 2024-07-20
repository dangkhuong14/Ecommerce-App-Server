package main

import (
	"log"
	"os"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ecommerce/common"
	"ecommerce/module/product/controller"
)



type Product struct {
	common.BaseModel
	CategoryID  int     `gorm:"column:category_id" json:"category_id"`
	Name        string  `gorm:"column:name" json:"name"`
	Image       *string `gorm:"column:image" json:"image"`
	Type        string  `gorm:"column:type" json:"type"`
	Description string  `gorm:"column:description" json:"description"`
}

type ProductUpdate struct {
	CategoryID  *int    `gorm:"column:category_id" json:"category_id"`
	Name        *string `gorm:"column:name" json:"name"`
	Image       *string `gorm:"column:image" json:"image"`
	Type        *string `gorm:"column:type" json:"type"`
	Description *string `gorm:"column:description" json:"description"`
	Status      *string `gorm:"column:status" json:"status"`
}



func (Product) TableName() string {
	return "products"
}



// func main() {
// 	// Checking that an environment variable is present or not.
// 	mysqlConnStr, ok := os.LookupEnv("MYSQL_CONNECTION")

// 	if !ok {
// 		log.Fatalln("Missing MySQL connection string")
// 	}

// 	dsn := mysqlConnStr
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		log.Fatalln("Cannot connect to MySQL:", err)
// 	}

// 	log.Println("Connected to MySQL:", db)

// 	// Create new Product instance
// 	now:= time.Now().UTC()
// 	productUUID, _ := uuid.NewV7()

// 	newProduct := Product{
// 		BaseModel: BaseModel{
// 			ID: UUID(productUUID),
// 			Status: "activated",
// 			CreatedAt: now,
// 			UpdatedAt: now,
// 		},
// 		CategoryID: 1,
// 		Name: "Americano",
// 		Image: nil,
// 		Type: "drink",
// 	}

// 	if err := db.Table(Product{}.TableName()).Create(&newProduct).Error; err != nil {
// 		log.Println("Cannot create new Product:", err)
// 	} else {
// 		log.Println("Product created:", newProduct)
// 	}

// 	// // Read Product instance
// 	// var existingProduct Product

// 	// if err:=db.
// 	// Table(Product{}.TableName()).
// 	// Where("id= ?", 1).
// 	// First(&existingProduct).Error; err !=nil {
// 	// 	log.Println("Cannot read from Product:", err)
// 	// } else {
// 	// 	log.Println("Product read:", existingProduct)
// 	// }

// 	//List products
// 	var products []Product

// 	if err:=db.
// 	Table(Product{}.TableName()).
// 	Where("status not in (?)", []string{"deactivated"}).
// 	Limit(3).
// 	Order("id desc").
// 	Offset(0).
// 	Find(&products).Error; err !=nil {
// 		log.Println("Cannot list from Product:", err)
// 	} else {
// 		log.Println("Products listed:", products[0].ID)
// 		log.Println("Products listed:", products[0].ID.String())
// 	}

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

	// Gin API Ping
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Create Product
	v1 := r.Group("/v1")
	{
		products := v1.Group("products")
		{
			products.POST("", controller.CreateProductAPI(db))
		}
	}

	r.Run(":3000")
}
