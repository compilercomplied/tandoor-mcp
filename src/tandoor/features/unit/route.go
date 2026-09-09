package unit

import (
	"context"
	"fmt"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
)

// CreateUnit calls POST /api/unit/
func CreateUnit(ctx context.Context, c *tandoor.Client, params UnitParam) (*UnitResponse, error) {
	res, err := tandoor.Request[UnitResponse](ctx, c, "POST", "/api/unit/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create unit: %w", err)
	}
	return res, nil
}

// ListUnits calls GET /api/unit/
func ListUnits(ctx context.Context, c *tandoor.Client, params ListUnitsParams) (*UnitListResponse, error) {
	qb := tandoor.NewQuery().
		Add("query", params.Query).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[UnitListResponse](ctx, c, "GET", "/api/unit/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list units: %w", err)
	}
	return res, nil
}

// GetUnit calls GET /api/unit/{id}/
func GetUnit(ctx context.Context, c *tandoor.Client, id int) (*UnitResponse, error) {
	path := fmt.Sprintf("/api/unit/%d/", id)
	res, err := tandoor.Request[UnitResponse](ctx, c, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit %d: %w", id, err)
	}
	return res, nil
}

// Merge combines source into target. Tandoor transfers references to target and removes source.
func Merge(ctx context.Context, c *tandoor.Client, source, target *UnitResponse) (*UnitResponse, error) {
	path := fmt.Sprintf("/api/unit/%d/merge/%d/", source.ID, target.ID)
	res, err := tandoor.Request[UnitResponse](ctx, c, "PUT", path, nil, source)
	if err != nil {
		return nil, fmt.Errorf("failed to merge unit %d into %d: %w", source.ID, target.ID, err)
	}
	return res, nil
}

func Update(ctx context.Context, c *tandoor.Client, id int, params UpdateParams) (*UnitResponse, error) {
	res, err := tandoor.Request[UnitResponse](ctx, c, "PATCH", fmt.Sprintf("/api/unit/%d/", id), nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update unit %d: %w", id, err)
	}
	return res, nil
}

// CreateUnitConversion calls POST /api/unit-conversion/
func CreateUnitConversion(ctx context.Context, c *tandoor.Client, params UnitConversionParam) (*UnitConversionResponse, error) {
	res, err := tandoor.Request[UnitConversionResponse](ctx, c, "POST", "/api/unit-conversion/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create unit conversion: %w", err)
	}
	return res, nil
}

// ListUnitConversions calls GET /api/unit-conversion/
func ListUnitConversions(ctx context.Context, c *tandoor.Client, params ListConversionsParams) (*UnitConversionListResponse, error) {
	qb := tandoor.NewQuery().
		Add("food_id", params.FoodID).
		Add("query", params.Query).
		Add("page", params.Page).
		Add("page_size", params.PageSize)

	res, err := tandoor.Request[UnitConversionListResponse](ctx, c, "GET", "/api/unit-conversion/", qb.Values(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list unit conversions: %w", err)
	}
	return res, nil
}

// GetUnitConversion calls GET /api/unit-conversion/{id}/
func GetUnitConversion(ctx context.Context, c *tandoor.Client, id int) (*UnitConversionResponse, error) {
	path := fmt.Sprintf("/api/unit-conversion/%d/", id)
	res, err := tandoor.Request[UnitConversionResponse](ctx, c, "GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit conversion %d: %w", id, err)
	}
	return res, nil
}
