package update_food

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/food"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	FoodID         int     `json:"food_id"`
	ExpectedName   string  `json:"expected_name"`
	Name           *string `json:"name,omitempty"`
	PluralName     *string `json:"plural_name,omitempty"`
	Description    *string `json:"description,omitempty"`
	IgnoreShopping *bool   `json:"ignore_shopping,omitempty"`
	ParentID       *int    `json:"parent_id,omitempty"`
}

func Register(s *mcp.Server, c *tandoor.Client) {
	mcp.AddTool(s, &mcp.Tool{Name: "update_food", Description: "Partially update a shared food. expected_name must match its current name to prevent stale-ID updates."}, func(ctx context.Context, r *mcp.CallToolRequest, a Args) (*mcp.CallToolResult, any, error) {
		if a.FoodID <= 0 || a.ExpectedName == "" {
			return fail("food_id and expected_name are required"), nil, nil
		}
		if a.Name == nil && a.PluralName == nil && a.Description == nil && a.IgnoreShopping == nil && a.ParentID == nil {
			return fail("at least one field to update is required"), nil, nil
		}
		existing, e := food.Get(ctx, c, a.FoodID)
		if e != nil {
			return api(e), nil, nil
		}
		if existing.Name != a.ExpectedName {
			return fail("food name no longer matches; retrieve it again before updating"), nil, nil
		}
		out, e := food.Update(ctx, c, a.FoodID, food.UpdateParams{Name: a.Name, PluralName: a.PluralName, Description: a.Description, IgnoreShopping: a.IgnoreShopping, Parent: a.ParentID})
		if e != nil {
			return api(e), nil, nil
		}
		return ok(out), nil, nil
	})
}
func fail(x string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Error: " + x}}, IsError: true}
}
func api(e error) *mcp.CallToolResult { return fail(fmt.Sprintf("updating food: %v", e)) }
func ok(v any) *mcp.CallToolResult {
	b, _ := json.MarshalIndent(v, "", "  ")
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}
}
