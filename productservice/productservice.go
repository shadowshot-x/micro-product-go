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

func handleNotInHeader(rw http.ResponseWriter, r *http.Request, param string) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(fmt.Sprintf("%s Missing", param)))
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

func (ctrl *ProductController) GetAllProductByIdHandler(rw http.ResponseWriter, r *http.Request) {
	// we know we will get a list of all products with a certain id.
	Products := store.Product{}
	if _, ok := r.Header["Id"]; !ok {
		ctrl.logger.Warn("Id was not found in the header")
		handleNotInHeader(rw, r, "Id")
		return
	}

	// here db.First fetches the first existing Elements in Products and stores them in Product
	// we need to record errors because if none exist, that is an error.
	err := db.First(&Products, r.Header["Id"][0]).Error

	if err != nil {
		ctrl.logger.Error("The stated record was not found")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Record not found"))
		return
	}
	// We can Send back all values to the ResponseWriter by jsonencoding the results
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(Products)
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
		// if we get an error, we dont want to execute any further
		ctrl.logger.Error("Error converting string to integer in inventory", zap.Error(err))
		return
	}
	// define the object we want to add
	newProduct := store.Product{
		Name:        r.Header["Productname"][0],
		VendorName:  r.Header["Productvendor"][0],
		Inventory:   inventory,
		Description: r.Header["Productdescription"][0],
		CreateAt:    time.Now(),
	}
	db.Omit("Id").Create(&newProduct)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Record was added"))
}

func (ctrl *ProductController) DeleteProductHandler(rw http.ResponseWriter, r *http.Request) {
	// This is the ideal case where we want to make sure the user is authenticated and has a certain role
	// however, right now I am keeping it simple
	// lets create an issue for this

	// first we see the request has the id for the product to be deleted.
	if _, ok := r.Header["Id"]; !ok {
		ctrl.logger.Warn("Id was not found in the header")
		handleNotInHeader(rw, r, "Id")
		return
	}

	// Now we know that the request has the parameter. Lets see how the gorm handles deletion
	// We can use the Where clause in gorm to query this.
	// db.Where(fmt.Sprintf("Id = %s", r.Header["Id"][0])).Delete(&store.Product{})
	// There is another easy way. As we know that Id is primary Key, we can do the following :-
	err := db.Delete(&store.Product{}, r.Header["Id"][0]).Error

	// this would mean there is an internal error
	if err != nil {
		ctrl.logger.Error("Could not delete the Given Product")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Record could not be deleted"))
		return
	}

	// succesful request
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Record deleted"))
}

// However, this is becoming a bit strict.
// We dont want to add a function everytime we get a new SQL query required.
// Doing this in gorm will be difficult.
// We would need a wrapper to provide user with this support, but this will take significant effort
// However something that executes queries as is can do the job.
// So, we can have an Exec Query of Gorm to deal with this. Custom Queries!

func (ctrl *ProductController) CustomQueryHandler(rw http.ResponseWriter, r *http.Request) {

	// to Query ie. to use SELECT we have .Raw in Gorm
	// to Execute like delete, add and update, we have .Exec in Gorm
	if _, ok := r.Header["Type"]; !ok {
		ctrl.logger.Warn("Type was not found in the header")
		handleNotInHeader(rw, r, "Type")
		return
	}

	if _, ok := r.Header["Query"]; !ok {
		ctrl.logger.Warn("Query was not found in the header")
		handleNotInHeader(rw, r, "Query")
		return
	}

	if r.Header["Type"][0] == "get" {
		// we know that we can only get an array of Products.
		var Products []store.Product
		err := db.Raw(r.Header["Query"][0]).Scan(&Products).Error

		if err != nil {
			ctrl.logger.Error("Could not Execute your Query")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Query not executed"))
			return
		}

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(Products)
		return
	} else if r.Header["Type"][0] == "exec" {
		err := db.Exec(r.Header["Query"][0]).Error

		if err != nil {
			ctrl.logger.Error("Could not Execute your Query")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Query not executed"))
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Query Executed"))
		return
	}

	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte("Incorrect Query Type"))
}
