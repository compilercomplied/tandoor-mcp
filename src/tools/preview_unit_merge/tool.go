package preview_unit_merge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/unit"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	SourceUnitID int `json:"source_unit_id" jsonschema:"The duplicate unit ID to remove."`
	TargetUnitID int `json:"target_unit_id" jsonschema:"The canonical unit ID to keep."`
}

type Response struct {
	Source unit.UnitResponse `json:"source"`
	Target unit.UnitResponse `json:"target"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "preview_unit_merge",
		Description: "Inspect a proposed unit merge without changing data. Source is the duplicate to remove; target is the canonical unit to keep.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing preview_unit_merge. source=%d, target=%d", args.SourceUnitID, args.TargetUnitID)
		if args.SourceUnitID <= 0 || args.TargetUnitID <= 0 || args.SourceUnitID == args.TargetUnitID {
			return validationError("source_unit_id and target_unit_id must be different positive IDs"), nil, nil
		}
		source, err := unit.GetUnit(ctx, client, args.SourceUnitID)
		if err != nil {
			return apiError("retrieving source unit", err), nil, nil
		}
		target, err := unit.GetUnit(ctx, client, args.TargetUnitID)
		if err != nil {
			return apiError("retrieving target unit", err), nil, nil
		}
		return success(Response{Source: *source, Target: *target}), nil, nil
	})
}

func validationError(message string) *mcp_sdk.CallToolResult {
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: "Error: " + message}}, IsError: true}
}
func apiError(action string, err error) *mcp_sdk.CallToolResult {
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: fmt.Sprintf("Error %s: %v", action, err)}}, IsError: true}
}
func success(response Response) *mcp_sdk.CallToolResult {
	b, _ := json.MarshalIndent(response, "", "  ")
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: string(b)}}}
}
