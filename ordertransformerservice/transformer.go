package ordertransformerservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/shadowshot-x/micro-product-go/ordertransformerservice/store"
	"go.uber.org/zap"
)

type OrderCompilation struct {
	APAC store.Orders
	EU   store.Orders
	NA   store.Orders
	SA   store.Orders
}

type RulesCompilation struct {
	rules map[string]store.Rules
}

type OrderTransformation struct {
	filteredOrders []store.Order
}

const rules_dir = "./ordertransformerservice/region_rules/"
const json_dir = "./ordertransformerservice/json_store/"

// TransformerHandler is the Transformer route handler
type TransformerController struct {
	logger           *zap.Logger
	Store_json_dir   string
	Region_rules_dir string
}

// NewTransformerController returns a frsh Transformer controller
func NewTransformerController(logger *zap.Logger) *TransformerController {
	return &TransformerController{
		logger:           logger,
		Store_json_dir:   json_dir,
		Region_rules_dir: rules_dir,
	}
}

func parser(store_directory, rules_directory string) (OrderCompilation, RulesCompilation, error) {
	files, err := ioutil.ReadDir(store_directory)
	if err != nil {
		return OrderCompilation{}, RulesCompilation{}, err
	}

	allOrders := OrderCompilation{}
	for _, file := range files {
		orderFile, err := ioutil.ReadFile(store_directory + file.Name())
		if err != nil {
			return OrderCompilation{}, RulesCompilation{}, err
		}
		contents := string(orderFile)
		orders, err := store.CreateOrdersStruct(orderFile)
		if err != nil {
			return OrderCompilation{}, RulesCompilation{}, err
		}

		if strings.Contains(contents, "APAC") {
			allOrders.APAC = orders
		} else if strings.Contains(contents, "EU") {
			allOrders.EU = orders
		} else if strings.Contains(contents, "NA") {
			allOrders.NA = orders
		} else if strings.Contains(contents, "SA") {
			allOrders.SA = orders
		} else {
			return OrderCompilation{}, RulesCompilation{}, errors.New("incorrect region provided in rules")
		}
	}

	files, err = ioutil.ReadDir(rules_directory)
	if err != nil {
		return OrderCompilation{}, RulesCompilation{}, err
	}

	allRules := RulesCompilation{
		rules: map[string]store.Rules{},
	}
	for _, file := range files {
		rulesFile, err := ioutil.ReadFile(rules_directory + file.Name())
		if err != nil {
			return OrderCompilation{}, RulesCompilation{}, err
		}
		contents := string(rulesFile)
		rules, err := store.CreateRulesStruct(rulesFile)
		if err != nil {
			return OrderCompilation{}, RulesCompilation{}, err
		}

		if strings.Contains(contents, "APAC") {
			allRules.rules["APAC"] = rules
			fmt.Println(allRules.rules["APAC"])
		} else if strings.Contains(contents, "EU") {
			allRules.rules["EU"] = rules
		} else if strings.Contains(contents, "NA") {
			allRules.rules["NA"] = rules
		} else if strings.Contains(contents, "SA") {
			allRules.rules["SA"] = rules
		} else {
			return OrderCompilation{}, RulesCompilation{}, errors.New("incorrect region provided in rules")
		}
	}
	return allOrders, allRules, nil
}

// validates a single order to all the rules corresponding to region
func validation(order store.Order, allRules RulesCompilation, region string) (bool, error) {
	regionRules := allRules.rules[region]

	for _, rule := range regionRules.RuleList {
		if rule.AmountFilter != "" {
			rg := regexp.MustCompile(`>|<|=`)
			filterAmt, err := strconv.ParseFloat(rg.Split(rule.AmountFilter, -1)[1], 64)
			if err != nil {
				return false, err
			}
			if rule.AmountFilter[0] == '>' {
				if order.Amount < filterAmt {
					return false, nil
				}
			} else if rule.AmountFilter[0] == '<' {
				if order.Amount > filterAmt {
					return false, nil
				}
			} else if rule.AmountFilter[0] == '=' {
				if order.Amount != filterAmt {
					return false, nil
				}
			}
		}
		if rule.EmailFilter != "" {
			if order.UserEmail != rule.EmailFilter {
				return false, nil
			}
		}
		if len(rule.BlacklistProduct) != 0 {
			for _, id := range rule.BlacklistProduct {
				for _, rg := range order.ProductList {
					if rg == id {
						return false, nil
					}
				}
			}
		}
	}

	return true, nil
}

func processOrderTransformation(finalFiles []store.Order) (string, error) {
	output := OrderTransformation{filteredOrders: finalFiles}
	byteOutput, err := json.Marshal(output)
	if err != nil {
		return "", err
	}
	c := &http.Client{}

	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(byteOutput))
	if err != nil {
		return "", err
	}
	response, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(responseBody), nil
}

func (ctrl *TransformerController) handleInternalError(rw http.ResponseWriter, region string, err error) {
	ctrl.logger.Error("Error in Validating", zap.Any("error", err), zap.String("region", region))
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte("Error in Validating for APAC"))
}

// adds the user to the database of users
func (ctrl *TransformerController) TransformerHandler(rw http.ResponseWriter, r *http.Request) {
	orderCompilation, ruleCompilation, err := parser(ctrl.Store_json_dir, ctrl.Region_rules_dir)
	if err != nil {
		ctrl.logger.Error("Error in Parsing orders and rules", zap.Any("error", err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error Parsing Files and rules"))
		return
	}

	ctrl.logger.Info("all orders", zap.Any("orders", orderCompilation))
	ctrl.logger.Info("all rules", zap.Any("rules", ruleCompilation))

	filteredFiles := []store.Order{}

	// Region 1
	apacOrders := orderCompilation.APAC
	for _, apacOrder := range apacOrders.OrderList {
		check, err := validation(apacOrder, ruleCompilation, "APAC")
		if err != nil {
			ctrl.handleInternalError(rw, "APAC", err)
			return
		}
		if check {
			filteredFiles = append(filteredFiles, apacOrder)
		} else {
			ctrl.logger.Info("Order Rejected", zap.Any("order", apacOrder))
		}
	}
	euOrders := orderCompilation.EU
	for _, euOrders := range euOrders.OrderList {
		check, err := validation(euOrders, ruleCompilation, "EU")
		if err != nil {
			ctrl.handleInternalError(rw, "APAC", err)
			return
		}
		if check {
			filteredFiles = append(filteredFiles, euOrders)
		} else {
			ctrl.logger.Info("Order Rejected", zap.Any("order", euOrders))
		}
	}
	naOrders := orderCompilation.NA
	for _, naOrders := range naOrders.OrderList {
		check, err := validation(naOrders, ruleCompilation, "NA")
		if err != nil {
			ctrl.handleInternalError(rw, "APAC", err)
			return
		}
		if check {
			filteredFiles = append(filteredFiles, naOrders)
		} else {
			ctrl.logger.Info("Order Rejected", zap.Any("order", naOrders))
		}
	}
	saOrders := orderCompilation.NA
	for _, saOrders := range saOrders.OrderList {
		check, err := validation(saOrders, ruleCompilation, "SA")
		if err != nil {
			ctrl.handleInternalError(rw, "APAC", err)
			return
		}
		if check {
			filteredFiles = append(filteredFiles, saOrders)
		} else {
			ctrl.logger.Info("Order Rejected", zap.Any("order", saOrders))
		}
	}

	linkinfo, err := processOrderTransformation(filteredFiles)
	if err != nil {
		ctrl.logger.Error("Error in sending file processed orders", zap.Any("error", err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error in sending file processed orders"))
	}

	ctrl.logger.Info("Link", zap.Any("Details", linkinfo))
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Files Parsed and Validated" + linkinfo + "\n" + fmt.Sprintf("%v", filteredFiles)))
}
