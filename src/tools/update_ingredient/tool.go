package update_ingredient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	IngredientID int      `json:"ingredient_id"`
	FoodID       *int     `json:"food_id,omitempty"`
	UnitID       *int     `json:"unit_id,omitempty"`
	Amount       *float64 `json:"amount,omitempty"`
	Note         *string  `json:"note,omitempty"`
	Order        *int     `json:"order,omitempty"`
	NoAmount     *bool    `json:"no_amount,omitempty"`
}

func Register(s *mcp.Server, c *tandoor.Client) {
	mcp.AddTool(s, &mcp.Tool{Name: "update_ingredient", Description: "Partially update a recipe ingredient's food, unit, amount, note, order, or no-amount setting."}, func(ctx context.Context, r *mcp.CallToolRequest, a Args) (*mcp.CallToolResult, any, error) {
		if a.IngredientID <= 0 {
			return fail("ingredient_id must be a positive ID"), nil, nil
		}
		if a.FoodID == nil && a.UnitID == nil && a.Amount == nil && a.Note == nil && a.Order == nil && a.NoAmount == nil {
			return fail("at least one field to update is required"), nil, nil
		}
		var f *ingredient.FoodRef
		if a.FoodID != nil {
			f = &ingredient.FoodRef{ID: *a.FoodID}
		}
		var u *ingredient.UnitRef
		if a.UnitID != nil {
			u = &ingredient.UnitRef{ID: *a.UnitID}
		}
		out, e := ingredient.Update(ctx, c, a.IngredientID, ingredient.UpdateParams{Food: f, Unit: u, Amount: a.Amount, Note: a.Note, Order: a.Order, NoAmount: a.NoAmount})
		if e != nil {
			return api(e), nil, nil
		}
		return ok(out), nil, nil
	})
}
func fail(x string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Error: " + x}}, IsError: true}
}
func api(e error) *mcp.CallToolResult { return fail(fmt.Sprintf("updating ingredient: %v", e)) }
func ok(v any) *mcp.CallToolResult {
	b, _ := json.MarshalIndent(v, "", "  ")
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}
}
