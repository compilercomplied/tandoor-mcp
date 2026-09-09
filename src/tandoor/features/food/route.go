package food

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// Create calls POST /api/food/
func Create(ctx context.Context, c *tandoor.Client, params FoodParam) (*FoodResponse, error) {
	res, err := tandoor.Request[FoodResponse](ctx, c, "POST", "/api/food/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create food: %w", err)
	}
	return res, nil
}

// List calls GET /api/food/
func List(ctx context.Context, c *tandoor.Client, params ListParams) (*FoodListResponse, error) {
	qb := tandoor.NewQuery().
		Add("query", params.Query).
		Add("root", params.Root).
		Add("tree", params.Tree).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[FoodListResponse](ctx, c, "GET", "/api/food/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list foods: %w", err)
	}
	return res, nil
}

// ListInheritFields calls GET /api/food-inherit-field/
func ListInheritFields(ctx context.Context, c *tandoor.Client) (*[]FoodInheritFieldResponse, error) {
	res, err := tandoor.Request[[]FoodInheritFieldResponse](ctx, c, "GET", "/api/food-inherit-field/", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list food inherit fields: %w", err)
	}
	return res, nil
}

// Get calls GET /api/food/{id}/
func Get(ctx context.Context, c *tandoor.Client, id int) (*FoodResponse, error) {
	path := fmt.Sprintf("/api/food/%d/", id)
	res, err := tandoor.Request[FoodResponse](ctx, c, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get food %d: %w", id, err)
	}
	return res, nil
}

// Merge combines source into target. Tandoor transfers references to target and removes source.
func Merge(ctx context.Context, c *tandoor.Client, source, target *FoodResponse) (*FoodResponse, error) {
	path := fmt.Sprintf("/api/food/%d/merge/%d/", source.ID, target.ID)
	res, err := tandoor.Request[FoodResponse](ctx, c, "PUT", path, nil, source)
	if err != nil {
		return nil, fmt.Errorf("failed to merge food %d into %d: %w", source.ID, target.ID, err)
	}
	return res, nil
}
