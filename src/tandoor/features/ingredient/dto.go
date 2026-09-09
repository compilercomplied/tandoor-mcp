package ingredient

// FoodRef is used to reference or create a Food by name or ID.
type FoodRef struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// UnitRef is used to reference or create a Unit by name or ID.
type UnitRef struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// IngredientParam is the request body for POST /api/ingredient/.
type IngredientParam struct {
	Food     FoodRef `json:"food"`
	Unit     UnitRef `json:"unit"`
	Amount   float64 `json:"amount"`
	Note     string  `json:"note,omitempty"`
	Order    *int    `json:"order,omitempty"`
	NoAmount bool    `json:"no_amount,omitempty"`
	RecipeID *int    `json:"-"`
}

// IngredientResponse is the response from POST /api/ingredient/.
type IngredientResponse struct {
	ID     int     `json:"id"`
	Food   FoodRef `json:"food"`
	Unit   UnitRef `json:"unit"`
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
	Order  int     `json:"order"`
}

type UpdateParams struct {
	Food     *FoodRef `json:"food,omitempty"`
	Unit     *UnitRef `json:"unit,omitempty"`
	Amount   *float64 `json:"amount,omitempty"`
	Note     *string  `json:"note,omitempty"`
	Order    *int     `json:"order,omitempty"`
	NoAmount *bool    `json:"no_amount,omitempty"`
}
