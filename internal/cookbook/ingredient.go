package cookbook

import (
	"fmt"
	"strconv"
	"strings"
)

type Ingredient struct {
	Amount  int    `json:"amount"`
	Measure string `json:"measure"`
	Name    string `json:"name"`
}

func NewIngredient(name string, amount int, measure string) Ingredient {
	return Ingredient{
		Name:    name,
		Amount:  amount,
		Measure: measure,
	}
}

func ParseIngredients(ingredients string) ([]Ingredient, error) {
	ingredientsList := strings.Split(ingredients, "\n")
	if len(ingredientsList) <= 1 {
		if ingredientsList[0] == "" {
			return nil, fmt.Errorf("must provide ingredients list")
		}
	}
	ret := make([]Ingredient, len(ingredientsList))

	for i, ingredient := range ingredientsList {
		ingredientSplit := strings.SplitN(ingredient, " ", 3)
		if len(ingredientSplit) != 3 {
			return []Ingredient{}, fmt.Errorf("malformed ingredient: must be number measure name")
		}
		intAmount, err := strconv.Atoi(ingredientSplit[0])
		if err != nil {
			return []Ingredient{}, fmt.Errorf("start or line must be a number")
		}
		ret[i] = NewIngredient(ingredientSplit[2], intAmount, ingredientSplit[1])
	}

	return ret, nil
}
