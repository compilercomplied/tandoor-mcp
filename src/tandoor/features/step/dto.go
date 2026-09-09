package step

import (
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
)

type StepParam struct {
	Name        string                          `json:"name,omitempty"`
	Instruction string                          `json:"instruction,omitempty"`
	Time        *int                            `json:"time,omitempty"`
	Order       *int                            `json:"order,omitempty"`
	Ingredients []ingredient.IngredientResponse `json:"ingredients"`
}

type StepResponse struct {
	ID          int                             `json:"id"`
	Name        string                          `json:"name"`
	Instruction string                          `json:"instruction"`
	Time        int                             `json:"time"`
	Order       int                             `json:"order"`
	Ingredients []ingredient.IngredientResponse `json:"ingredients"`
}

// UpdateStepParam contains the step fields that may be updated. Nil fields are left unchanged.
type UpdateStepParam struct {
	Name                 *string                          `json:"name,omitempty"`
	Instruction          *string                          `json:"instruction,omitempty"`
	Time                 *int                             `json:"time,omitempty"`
	Order                *int                             `json:"order,omitempty"`
	Ingredients          *[]ingredient.IngredientResponse `json:"ingredients,omitempty"`
	ShowAsHeader         *bool                            `json:"show_as_header,omitempty"`
	ShowIngredientsTable *bool                            `json:"show_ingredients_table,omitempty"`
}
