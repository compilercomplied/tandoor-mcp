package update_unit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/unit"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	UnitID       int     `json:"unit_id"`
	Name         *string `json:"name,omitempty"`
	PluralName   *string `json:"plural_name,omitempty"`
	Description  *string `json:"description,omitempty"`
	BaseUnit     *string `json:"base_unit,omitempty"`
	OpenDataSlug *string `json:"open_data_slug,omitempty"`
}

func Register(s *mcp.Server, c *tandoor.Client) {
	mcp.AddTool(s, &mcp.Tool{Name: "update_unit", Description: "Partially update a canonical unit's name, plural name, description, base unit, or Open Data slug."}, func(ctx context.Context, r *mcp.CallToolRequest, a Args) (*mcp.CallToolResult, any, error) {
		if a.UnitID <= 0 {
			return fail("unit_id must be a positive ID"), nil, nil
		}
		if a.Name == nil && a.PluralName == nil && a.Description == nil && a.BaseUnit == nil && a.OpenDataSlug == nil {
			return fail("at least one field to update is required"), nil, nil
		}
		out, e := unit.Update(ctx, c, a.UnitID, unit.UpdateParams{Name: a.Name, PluralName: a.PluralName, Description: a.Description, BaseUnit: a.BaseUnit, OpenDataSlug: a.OpenDataSlug})
		if e != nil {
			return api(e), nil, nil
		}
		return ok(out), nil, nil
	})
}
func fail(x string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: "Error: " + x}}, IsError: true}
}
func api(e error) *mcp.CallToolResult { return fail(fmt.Sprintf("updating unit: %v", e)) }
func ok(v any) *mcp.CallToolResult {
	b, _ := json.MarshalIndent(v, "", "  ")
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}
}
