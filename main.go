package main

import (
	"log"
	"os"
	"strings"
	"time"

	"database/sql/driver"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        UUID      `gorm:"column:id" json:"id"`
	Status    string    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Product struct {
	BaseModel
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

func GenNewBaseModel() BaseModel {
	now := time.Now().UTC()
	newUUID, _ := uuid.NewV7()

	return BaseModel{
		ID:        UUID(newUUID),
		Status:    "activated",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (Product) TableName() string {
	return "products"
}

type UUID uuid.UUID

func (id *UUID) Scan(src interface{}) error {
	var sqlID uuid.UUID

	if err := sqlID.Scan(src); err != nil {
		return err
	}

	*id = UUID(sqlID)

	return nil
}

func (id UUID) Value() (driver.Value, error) {
	return uuid.UUID(id).MarshalBinary()
}

// Method chuyển UUID từ byte slice thành chuỗi
func (id UUID) String() string {
	return uuid.UUID(id).String()
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
	products := v1.Group("products")
	products.POST("", func(c *gin.Context) {
		// Parse product data from body
		var productData Product

		if err := c.Bind(&productData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Business's logic
		productData.Name = strings.TrimSpace(productData.Name)
		if productData.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "product's name can not be blank",
			})
			return
		}

		productData.BaseModel = GenNewBaseModel()

		// Save to database
		if err := db.Table(Product{}.TableName()).Create(&productData).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "can not create new product",
			})
			return
		} else {
			log.Println("Product created:", productData)
			c.JSON(http.StatusCreated, gin.H{
				"data": productData.ID.String(),
			})
		}
	})

	r.Run(":3000")
}
