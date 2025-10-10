package cookbook

import (
	"fmt"
	"strconv"
	"strings"
)

type Ingredient struct {
	Amount  int
	Measure string
	Name    string
}

func NewIngredient(name string, amount int, measure string) Ingredient {
	return Ingredient{
		Name:    name,
		Amount:  amount,
		Measure: measure,
	}
}

func ParseIngredients(ingredients string) ([]Ingredient, error) {
	ret := make([]Ingredient, 0)
	ingredientsList := strings.Split(ingredients, "\n")

	for _, ingredient := range ingredientsList {
		ingredientSplit := strings.SplitN(ingredient, " ", 3)
		if len(ingredientSplit) != 3 {
			return []Ingredient{}, fmt.Errorf("malformed ingredient")
		}
		intAmount, err := strconv.Atoi(ingredientSplit[0])
		if err != nil {
			return []Ingredient{}, err
		}
		ret = append(ret, NewIngredient(ingredientSplit[2], intAmount, ingredientSplit[1]))
	}

	return ret, nil
}
