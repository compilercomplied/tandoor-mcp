package ingredient

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

type recipeSteps struct {
	Steps []recipeStep `json:"steps"`
}

type recipeStep struct {
	Name        string               `json:"name,omitempty"`
	Instruction string               `json:"instruction,omitempty"`
	Time        *int                 `json:"time,omitempty"`
	Order       *int                 `json:"order,omitempty"`
	Ingredients []IngredientResponse `json:"ingredients"`
}

// Create calls POST /api/ingredient/ to create a standalone ingredient.
// If params.RecipeID is provided, it associates the new ingredient with the specified recipe.
func Create(ctx context.Context, c *tandoor.Client, params IngredientParam) (*IngredientResponse, error) {
	res, err := tandoor.Request[IngredientResponse](ctx, c, "POST", "/api/ingredient/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create ingredient: %w", err)
	}

	if params.RecipeID != nil && *params.RecipeID > 0 {
		recipeID := *params.RecipeID
		recipeEndpoint := fmt.Sprintf("/api/recipe/%d/", recipeID)

		// Fetch existing steps to append or create a new step
		existingRecipe, err := tandoor.Request[recipeSteps](ctx, c, "GET", recipeEndpoint, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch recipe %d: %w", recipeID, err)
		}

		if len(existingRecipe.Steps) > 0 {
			// Add to the first step's ingredients
			existingRecipe.Steps[0].Ingredients = append(existingRecipe.Steps[0].Ingredients, *res)
		} else {
			// Create a default step with the ingredient
			defaultStep := recipeStep{
				Name:        "Ingredients",
				Instruction: "Prepare ingredients",
				Ingredients: []IngredientResponse{*res},
			}
			existingRecipe.Steps = append(existingRecipe.Steps, defaultStep)
		}

		// PATCH the recipe steps
		_, err = tandoor.Request[any](ctx, c, "PATCH", recipeEndpoint, nil, existingRecipe)
		if err != nil {
			return nil, fmt.Errorf("failed to associate ingredient with recipe %d: %w", recipeID, err)
		}
	}

	return res, nil
}

// Get calls GET /api/ingredient/<id>/ to retrieve an ingredient by ID.
func Get(ctx context.Context, c *tandoor.Client, id int) (*IngredientResponse, error) {
	res, err := tandoor.Request[IngredientResponse](ctx, c, "GET", fmt.Sprintf("/api/ingredient/%d/", id), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredient: %w", err)
	}
	return res, nil
}

func Update(ctx context.Context, c *tandoor.Client, id int, params UpdateParams) (*IngredientResponse, error) {
	res, err := tandoor.Request[IngredientResponse](ctx, c, "PATCH", fmt.Sprintf("/api/ingredient/%d/", id), nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update ingredient %d: %w", id, err)
	}
	return res, nil
}
