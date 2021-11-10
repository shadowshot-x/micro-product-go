package productservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shadowshot-x/micro-product-go/productservice/store"
	"go.uber.org/zap"
)

var db *gorm.DB

func GetSecret() string {
	return os.Getenv("MYSQL_SECRET")
}

// ProductController is the getproduct route handler
type ProductController struct {
	logger *zap.Logger
}

// NewProductController returns a frsh Upload controller
func NewProductController(logger *zap.Logger) *ProductController {
	return &ProductController{
		logger: logger,
	}
}

func (ctrl *ProductController) InitGormConnection() {
	// database configuration for mysql
	// first we fetch the mysql secret string stored in environment variables
	sqlsecret := GetSecret()
	// if secret is empty, we want to warn the user
	if sqlsecret == "" {
		ctrl.logger.Warn("Unable to get mysql secret")
		return
	}
	var err error
	// lets open the conncection
	db, err = gorm.Open("mysql", GetSecret())
	if err != nil {
		ctrl.logger.Warn("Connection Failed to Open", zap.Error(err))
	} else {
		ctrl.logger.Info("Connection Established")
	}

	//We have the database name in our Environment secret.
	// Auto Migrate creates a table named products in that Database
	db.AutoMigrate(&store.Product{})
}

func (ctrl *ProductController) GetAllProductsHandler(rw http.ResponseWriter, r *http.Request) {
	// we know we will get a list of all products.
	AllProducts := []store.Product{}
	// here db.Find fetches all the existing Elements in Products and stores them in AllProducts
	db.Find(&AllProducts)
	// We can Send back all values to the ResponseWriter by jsonencoding the results
	json.NewEncoder(rw).Encode(AllProducts)
}

func handleNotInHeader(rw http.ResponseWriter, r *http.Request, param string) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(fmt.Sprintf("%s Missing", param)))
}
func (ctrl *ProductController) AddProductHandler(rw http.ResponseWriter, r *http.Request) {
	//validate the request first
	if _, ok := r.Header["Productname"]; !ok {
		ctrl.logger.Warn("Name was not found in the header")
		handleNotInHeader(rw, r, "name")
		return
	}
	if _, ok := r.Header["Productvendor"]; !ok {
		ctrl.logger.Warn("Vendor was not found in the header")
		handleNotInHeader(rw, r, "Vendor")
		return
	}
	if _, ok := r.Header["Productinventory"]; !ok {
		ctrl.logger.Warn("Inventory was not found in the header")
		handleNotInHeader(rw, r, "Inventory")
		return
	}
	if _, ok := r.Header["Productdescription"]; !ok {
		ctrl.logger.Warn("Description was not found in the header")
		handleNotInHeader(rw, r, "Description")
		return
	}
	// We want to get the details of the Product first. So these have to be in the request
	inventory, err := strconv.Atoi(r.Header["Productinventory"][0])
	if err != nil {
		ctrl.logger.Error("Error converting string to integer in inventory", zap.Error(err))
	}
	newProduct := store.Product{
		Name:        r.Header["Productname"][0],
		VendorName:  r.Header["Productvendor"][0],
		Inventory:   inventory,
		Description: r.Header["Productdescription"][0],
		CreateAt:    time.Now(),
	}
	ctrl.logger.Info("1 Product was Added")
	db.Omit("Id").Create(&newProduct)
	ctrl.logger.Info("Product was Added")
}
