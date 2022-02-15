package main

import (
	"net/http"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shadowshot-x/micro-product-go/authservice"
	"github.com/shadowshot-x/micro-product-go/authservice/middleware"
	"github.com/shadowshot-x/micro-product-go/clientclaims"
	"github.com/shadowshot-x/micro-product-go/couponservice"
	"github.com/shadowshot-x/micro-product-go/ordertransformerservice"
	retrospectiveservice "github.com/shadowshot-x/micro-product-go/privateretrospectiveservice"
	"github.com/shadowshot-x/micro-product-go/productservice"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"
)

func PingHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Your App seems Healthy"))
}

// lets set up prometheus custom metrics
var retrospectiveAvailSuccess = promauto.NewCounter(prometheus.CounterOpts{
	Name: "retrospective_avail_success",
	Help: "Successful availing of Retrospective",
})
var retrospectiveAvailFail = promauto.NewCounter(prometheus.CounterOpts{
	Name: "retrospective_avail_fail",
	Help: "Failure to avail Retrospective due to Mutex lock",
})
var retrospectiveUpdateSuccess = promauto.NewCounter(prometheus.CounterOpts{
	Name: "retrospective_update_success",
	Help: "Successful updation of Retrospective",
})
var retrospectiveUpdateFail = promauto.NewCounter(prometheus.CounterOpts{
	Name: "retrospective_update_fail",
	Help: "Failure to update Retrospective due to Mutex lock",
})
var retrospectiveReleaseSuccess = promauto.NewCounter(prometheus.CounterOpts{
	Name: "retrospective_release_success",
	Help: "Successful release of Retrospective mutex",
})
var retrospectiveReleaseFail = promauto.NewCounter(prometheus.CounterOpts{
	Name: "retrospective_release_fail",
	Help: "Failure to release Retrospective due to Authorization",
})

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
	transc := ordertransformerservice.NewTransformerController(log)

	redisInstance := couponservice.RedisInstanceGenerator(log)
	cc := couponservice.NewCouponStreamController(log, redisInstance)

	melodyInstance := melody.New()
	rc := retrospectiveservice.NewRetrospectiveController(log, melodyInstance, retrospectiveAvailSuccess, retrospectiveAvailFail,
		retrospectiveUpdateSuccess, retrospectiveUpdateFail, retrospectiveReleaseSuccess, retrospectiveReleaseFail)

	// ping function
	mainRouter.HandleFunc("/ping", PingHandler)

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

	// Retrospective Service SubRouter
	retrospectiveRouter := mainRouter.PathPrefix("/retrospective").Subrouter()
	retrospectiveRouter.HandleFunc("/wc", rc.AssignSocket).Methods("POST")
	retrospectiveRouter.HandleFunc("/avail", rc.AvailRetrospective).Methods("GET")
	retrospectiveRouter.HandleFunc("/release", rc.ReleaseRetrospective).Methods("POST")
	retrospectiveRouter.HandleFunc("/check", rc.CheckAccess).Methods("GET")
	retrospectiveRouter.HandleFunc("/checkstring", rc.BroadcastMessage).Methods("GET")
	retrospectiveRouter.HandleFunc("/change", rc.ChangeRetrospective).Methods("POST")

	// Transformer Service SubRouter
	transformerOrderRouter := mainRouter.PathPrefix("/transformer").Subrouter()
	transformerOrderRouter.HandleFunc("/transform", transc.TransformerHandler).Methods("GET")

	// CORS Header
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	// Adding Prometheus http handler to expose the metrics
	// this will display our metrics as well as some standard metrics
	mainRouter.Path("/metrics").Handler(promhttp.Handler())
	// Add the Middleware to different subrouter
	// HTTP Server
	// Add Time outs
	server := &http.Server{
		Addr:    ":9090",
		Handler: ch(mainRouter),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Error("Error Booting the Server", zap.Error(err))
	}
}
