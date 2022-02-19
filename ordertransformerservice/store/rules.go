package store

import "gopkg.in/yaml.v2"

type Rules struct {
	Region   string
	RuleList []Rule
}

type Rule struct {
	AmountFilter     string
	BlacklistProduct []string
	EmailFilter      string
}

func CreateRulesStruct(inp []byte) (Rules, error) {
	output := Rules{}
	err := yaml.Unmarshal(inp, &output)
	if err != nil {
		return Rules{}, err
	}
	return output, nil
}
