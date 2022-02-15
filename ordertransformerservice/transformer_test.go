package ordertransformerservice

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/shadowshot-x/micro-product-go/ordertransformerservice/store"
)

// this is a normal testing scenario where a simple code workflow is tested.
// It is tested based on possible errors.
func TestValidation(t *testing.T) {
	order := store.Order{
		OrderId:     "1",
		ProductList: []string{"2", "3"},
		Amount:      100,
		UserEmail:   "abc@example.com",
	}
	allRules := RulesCompilation{
		rules: map[string]store.Rules{
			"APAC": {
				Region: "APAC",
				RuleList: []store.Rule{
					{
						AmountFilter: "<200",
						EmailFilter:  "abc@example.com",
					},
				},
			},
		},
	}
	t.Run("All Good", func(t *testing.T) {
		check, err := validation(order, allRules, "APAC")
		if err != nil {
			t.Fatalf("Got unexpected error: %v", err)
		}
		if !check {
			t.Fatalf("Correct Rule Invalidated!")
		}
	})
	t.Run("Amount Filter Check", func(t *testing.T) {
		allRules.rules["APAC"].RuleList[0].AmountFilter = ">200"
		check, err := validation(order, allRules, "APAC")
		if err != nil {
			t.Fatalf("Got unexpected error: %v", err)
		}
		if check {
			t.Fatalf("Incorrect Rule Validated!")
		}

		allRules.rules["APAC"].RuleList[0].AmountFilter = "=200"
		check, err = validation(order, allRules, "APAC")
		if err != nil {
			t.Fatalf("Got unexpected error: %v", err)
		}
		if check {
			t.Fatalf("Incorrect Rule Validated!")
		}

		allRules.rules["APAC"].RuleList[0].AmountFilter = "<10"
		check, err = validation(order, allRules, "APAC")
		if err != nil {
			t.Fatalf("Got unexpected error: %v", err)
		}
		if check {
			t.Fatalf("Incorrect Rule Validated!")
		}

		allRules.rules["APAC"].RuleList[0].AmountFilter = "<200"
	})
	t.Run("Email Filter Check", func(t *testing.T) {
		allRules.rules["APAC"].RuleList[0].EmailFilter = "incorrect@example.com"
		check, err := validation(order, allRules, "APAC")
		if err != nil {
			t.Fatalf("Got unexpected error: %v", err)
		}
		if check {
			t.Fatalf("Incorrect Rule Validated!")
		}
		allRules.rules["APAC"].RuleList[0].EmailFilter = "abc@example.com"
	})
	t.Run("Blacklist Check", func(t *testing.T) {
		allRules.rules["APAC"].RuleList[0].BlacklistProduct = []string{"2"}
		check, err := validation(order, allRules, "APAC")
		if err != nil {
			t.Fatalf("Got unexpected error: %v", err)
		}
		if check {
			t.Fatalf("Incorrect Rule Validated!")
		}
	})

	t.Run("Error: Incorrect Filter Provided", func(t *testing.T) {
		allRules.rules["APAC"].RuleList[0].AmountFilter = ">>>>"
		_, err := validation(order, allRules, "APAC")
		if err == nil {
			t.Fatalf("Did not get error when expected: %v", err)
		}
	})
}

// in this test, we need to build temporary files
func TestParser(t *testing.T) {
	// we need to set up temporary directories for this case.
	// golang provides this in ioutil.TempDir
	dir1, err := ioutil.TempDir("./", "")
	if err != nil {
		t.Fatalf("Unable to create Temp Dir 1 : %s", dir1)
	}
	dir2, err := ioutil.TempDir("./", "")
	if err != nil {
		t.Fatalf("Unable to create Temp Dir 1 : %s", dir2)
	}
	// we need to make sure the temp directories are removed
	defer os.RemoveAll(dir1)
	defer os.RemoveAll(dir2)
	// now we can set up temporary files
	oTempName := "apacOrder.json"
	f1, err := os.Create(fmt.Sprintf("./%s/%s", dir1, oTempName))
	if err != nil {
		t.Fatalf("Unable to create temporary file %s", oTempName)
	}
	contents := `{
		"region": "APAC",
		"orderlist": [{
			"orderid" : "1",
		"amount":     33.44,
		"useremail":  "abc@example.com",
		"create_at":   ""
		}]
	}`
	f1.WriteString(contents)
	rTempName := "apacDirective.yaml"
	f2, err := os.Create(fmt.Sprintf("./%s/%s", dir2, rTempName))
	if err != nil {
		t.Fatalf("Unable to create temporary file %s", rTempName)
	}
	contents = `
region: APAC
rulelist:
  - amountfilter: "<18000"
`
	f2.WriteString(contents)
	t.Run("Good: All Pass", func(t *testing.T) {
		oc, rc, err := parser(fmt.Sprintf("./%s/", dir1), fmt.Sprintf("./%s/", dir2))
		if err != nil {
			t.Fatalf("Got unexpected error %v", err)
		}
		expectedOc := OrderCompilation{
			APAC: store.Orders{
				Region: "APAC",
				OrderList: []store.Order{
					{
						OrderId:   "1",
						Amount:    33.44,
						UserEmail: "abc@example.com",
					},
				},
			},
		}
		expectedRc := RulesCompilation{
			rules: map[string]store.Rules{
				"APAC": {
					Region: "APAC",
					RuleList: []store.Rule{
						{
							AmountFilter: "<18000",
						},
					},
				},
			},
		}
		if !reflect.DeepEqual(oc, expectedOc) {
			t.Fatalf("Incorrect OrderCompilation Generated. %v,  %v", oc, expectedOc)
		}
		if !reflect.DeepEqual(rc, expectedRc) {
			t.Fatalf("Incorrect ruleCompilation Generated. %v,  %v", rc, expectedRc)
		}
	})

}
