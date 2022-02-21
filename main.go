package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shadowshot-x/micro-product-go/authservice"
	"github.com/shadowshot-x/micro-product-go/authservice/middleware"
	"github.com/shadowshot-x/micro-product-go/clientclaims"
	"github.com/shadowshot-x/micro-product-go/couponservice"
	"github.com/shadowshot-x/micro-product-go/monitormodule"
	"github.com/shadowshot-x/micro-product-go/ordertransformerservice"
	"github.com/shadowshot-x/micro-product-go/productservice"
	"go.uber.org/zap"
)

func PingHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Your App seems Healthy"))
}

func simplePostHandler(rw http.ResponseWriter, r *http.Request) {
	fileName, err := os.OpenFile("./metricDetails.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("error in file ops", zap.Error(err))
	}
	defer fileName.Close()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fileName.Write([]byte(reqBody))
	fileName.Write([]byte("\n"))

	rw.Write([]byte("Post Request Recieved for the Success"))
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Starting...")

	err := godotenv.Load(".env")

	if err != nil {
		log.Error("Error loading .env file", zap.Error(err))
	}

	err = monitormodule.MonitorBinder(log)
	if err != nil {
		fmt.Println(err)
		return
	}

	mainRouter := mux.NewRouter()

	suc := authservice.NewSignupController(log)
	sic := authservice.NewSigninController(log)
	uc := clientclaims.NewUploadController(log)
	dc := clientclaims.NewDownloadController(log)
	tm := middleware.NewTokenMiddleware(log)
	pc := productservice.NewProductController(log)
	transc := ordertransformerservice.NewTransformerController(log)

	redisInstance := couponservice.RedisInstanceGenerator(log)
	cc := couponservice.NewCouponStreamController(log, redisInstance)

	// ping function
	mainRouter.HandleFunc("/ping", PingHandler)
	mainRouter.HandleFunc("/checkRoutine", simplePostHandler).Methods("POST")

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
	// pc.InitGormConnection()
	productRouter := mainRouter.PathPrefix("/product").Subrouter()
	productRouter.HandleFunc("/getprods", pc.GetAllProductsHandler).Methods("GET")
	productRouter.HandleFunc("/addprod", pc.AddProductHandler).Methods("POST")
	productRouter.HandleFunc("/getprodbyid", pc.GetAllProductByIdHandler).Methods("GET")
	productRouter.HandleFunc("/deletebyid", pc.DeleteProductHandler).Methods("DELETE")
	productRouter.HandleFunc("/customquery", pc.CustomQueryHandler).Methods("GET", "POST")

	//Coupon Service SubRouter
	couponRouter := mainRouter.PathPrefix("/coupon").Subrouter()
	couponRouter.HandleFunc("/addcoupon", cc.AddCouponList).Methods("POST")
	couponRouter.HandleFunc("/getvendorcoupons", cc.GetCouponForInternalValidation).Methods("GET")
	couponRouter.HandleFunc("/delregionstream", cc.PurgeStream).Methods("DELETE")

	// Transformer Service SubRouter
	transformerOrderRouter := mainRouter.PathPrefix("/transformer").Subrouter()
	transformerOrderRouter.HandleFunc("/transform", transc.TransformerHandler).Methods("GET")

	// CORS Header
	cors := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	// Adding Prometheus http handler to expose the metrics
	// this will display our metrics as well as some standard metrics
	mainRouter.Path("/metrics").Handler(promhttp.Handler())
	// Add the Middleware to different subrouter
	// HTTP Server
	// Add Time outs
	server := &http.Server{
		Addr:    ":9090",
		Handler: cors(mainRouter),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Error("Error Booting the Server", zap.Error(err))
	}
}
