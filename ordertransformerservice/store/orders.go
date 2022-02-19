package store

import "encoding/json"

type Orders struct {
	Region    string
	OrderList []Order
}

type Order struct {
	OrderId     string
	ProductList []string
	Amount      float64
	UserEmail   string
	UserAddress string
	Create_At   string
}

func CreateOrdersStruct(inp []byte) (Orders, error) {
	output := Orders{}
	err := json.Unmarshal(inp, &output)
	if err != nil {
		return Orders{}, err
	}
	return output, nil
}
