package main

import (
	"net/http"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/shadowshot-x/micro-product-go/authservice"
	"github.com/shadowshot-x/micro-product-go/authservice/middleware"
	"github.com/shadowshot-x/micro-product-go/clientclaims"
	"github.com/shadowshot-x/micro-product-go/productservice"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Starting...")

	err := godotenv.Load(".env")

	if err != nil {
		log.Error("Error loading .env file", zap.Error(err))
	}

	mainRouter := mux.NewRouter()

	suc := authservice.NewSignupController(log)
	sic := authservice.NewSigninController(log)
	uc := clientclaims.NewUploadController(log)
	dc := clientclaims.NewDownloadController(log)
	tm := middleware.NewTokenMiddleware(log)
	pc := productservice.NewProductController(log)

	// We will create a Subrouter for Authentication service
	// route for sign up and signin. The Function will come from auth-service package
	// checks if given header params already exists. If not,it adds the user
	authRouter := mainRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", suc.SignupHandler).Methods("POST")

	// The Signin will send the JWT Token back as we are making microservices.
	// JWT token will make sure that other services are protected.
	// So, ultimately, we would need a middleware
	authRouter.HandleFunc("/signin", sic.SigninHandler).Methods("GET")

	// File Upload SubRouter
	claimsRouter := mainRouter.PathPrefix("/claims").Subrouter()
	claimsRouter.HandleFunc("/upload", uc.UploadFile)
	claimsRouter.HandleFunc("/download", dc.DownloadFile)
	claimsRouter.Use(tm.TokenValidationMiddleware)

	//Initialize the Gorm connection
	pc.InitGormConnection()
	productRouter := mainRouter.PathPrefix("/product").Subrouter()
	productRouter.HandleFunc("/getprods", pc.GetAllProductsHandler).Methods("GET")
	productRouter.HandleFunc("/addprod", pc.AddProductHandler).Methods("POST")
	productRouter.HandleFunc("/getprodbyid", pc.GetAllProductByIdHandler).Methods("GET")
	productRouter.HandleFunc("/deletebyid", pc.DeleteProductHandler).Methods("DELETE")
	productRouter.HandleFunc("/customquery", pc.CustomQueryHandler).Methods("GET", "POST")

	// CORS Header
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))
	// Add the Middleware to different subrouter
	// HTTP Server
	// Add Time outs
	server := &http.Server{
		Addr:    "127.0.0.1:9090",
		Handler: ch(mainRouter),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Error("Error Booting the Server", zap.Error(err))
	}
}
