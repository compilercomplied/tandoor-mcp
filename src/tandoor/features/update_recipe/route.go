package update_recipe

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Update partially updates a recipe without replacing fields that were not supplied.
func Update(ctx context.Context, c *tandoor.Client, recipeID int, params Params) (*Response, error) {
	path := fmt.Sprintf("/api/recipe/%d/", recipeID)
	res, err := tandoor.Request[Response](ctx, c, "PATCH", path, nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update recipe %d: %w", recipeID, err)
	}
	return res, nil
}
