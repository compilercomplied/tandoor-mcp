package step

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

type RecipeSteps struct {
	Steps []StepParam `json:"steps"`
}

// RecipeStepsResponse is used to decode PATCH /api/recipe/<id>/ responses
// which include the assigned step IDs from Tandoor.
type RecipeStepsResponse struct {
	Steps []StepResponse `json:"steps"`
}

// Create executes a create step request by fetching the existing recipe steps and
// appending the new one via PATCH /api/recipe/<recipeID>/.
// Note: step_recipe (recipe-in-recipe nesting) is NOT set here — recipeID is only
// used to determine which recipe endpoint to PATCH.
func Create(ctx context.Context, c *tandoor.Client, recipeID int, params StepParam) (*StepResponse, error) {
	recipeEndpoint := fmt.Sprintf("/api/recipe/%d/", recipeID)

	// Fetch the existing steps so we can append without clobbering them.
	existingRecipe, err := tandoor.Request[RecipeSteps](ctx, c, "GET", recipeEndpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recipe: %w", err)
	}

	// Append the new step and PATCH the recipe.
	existingRecipe.Steps = append(existingRecipe.Steps, params)
	updatedRecipe, err := tandoor.Request[RecipeStepsResponse](ctx, c, "PATCH", recipeEndpoint, nil, existingRecipe)
	if err != nil {
		return nil, fmt.Errorf("failed to update recipe with new step: %w", err)
	}

	// Return the newly created step (the last one in the response array).
	if len(updatedRecipe.Steps) > 0 {
		lastStep := updatedRecipe.Steps[len(updatedRecipe.Steps)-1]
		return &lastStep, nil
	}

	return nil, fmt.Errorf("failed to extract created step from Tandoor response")
}

// Update partially updates a step without replacing fields that were not supplied.
func Update(ctx context.Context, c *tandoor.Client, stepID int, params UpdateStepParam) (*StepResponse, error) {
	path := fmt.Sprintf("/api/step/%d/", stepID)
	res, err := tandoor.Request[StepResponse](ctx, c, "PATCH", path, nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update step %d: %w", stepID, err)
	}
	return res, nil
}
