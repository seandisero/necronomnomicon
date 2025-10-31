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

func ParseIngredientsFromForm(ingredients string) ([]Ingredient, error) {
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

func ParseIngredientsFromDB(text string) ([]Ingredient, error) {
	ingredientList := strings.Split(text, "?")
	ingredients := make([]Ingredient, len(ingredientList))
	for i, ingstr := range ingredientList {
		parts := strings.SplitN(ingstr, ":", 3)
		if len(parts) < 3 {
			return nil, fmt.Errorf(
				"malformed ingredient in list for ingredient: %s",
				ingstr,
			)
		}

		amount, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("malformed ingredient amount %v", err)
		}
		ingredient := Ingredient{
			Amount:  amount,
			Measure: parts[1],
			Name:    parts[2],
		}
		ingredients[i] = ingredient
	}
	return ingredients, nil
}

func ComposeIngredientsForDB(ingredients []Ingredient) string {
	ingredientList := make([]string, len(ingredients))
	for i, ingredient := range ingredients {
		parts := make([]string, 3)
		amount := strconv.Itoa(ingredient.Amount)
		parts[0] = amount
		parts[1] = ingredient.Measure
		parts[2] = ingredient.Name
		joined := strings.Join(parts, ":")
		ingredientList[i] = joined
	}
	return strings.Join(ingredientList, "?")
}
