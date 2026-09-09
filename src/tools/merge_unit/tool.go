package merge_unit

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
	SourceUnitID       int    `json:"source_unit_id" jsonschema:"The duplicate unit ID to remove."`
	TargetUnitID       int    `json:"target_unit_id" jsonschema:"The canonical unit ID to keep."`
	ExpectedSourceName string `json:"expected_source_name" jsonschema:"Source name returned by preview_unit_merge; merge stops if it no longer matches."`
	ExpectedTargetName string `json:"expected_target_name" jsonschema:"Target name returned by preview_unit_merge; merge stops if it no longer matches."`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "merge_unit",
		Description: "Merge a duplicate unit into a canonical unit. Run preview_unit_merge first and pass its exact source and target names to guard against stale IDs. This removes the source unit.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing merge_unit. source=%d, target=%d", args.SourceUnitID, args.TargetUnitID)
		if args.SourceUnitID <= 0 || args.TargetUnitID <= 0 || args.SourceUnitID == args.TargetUnitID || args.ExpectedSourceName == "" || args.ExpectedTargetName == "" {
			return validationError("source_unit_id and target_unit_id must be different positive IDs and expected names are required"), nil, nil
		}
		source, err := unit.GetUnit(ctx, client, args.SourceUnitID)
		if err != nil {
			return apiError("retrieving source unit", err), nil, nil
		}
		target, err := unit.GetUnit(ctx, client, args.TargetUnitID)
		if err != nil {
			return apiError("retrieving target unit", err), nil, nil
		}
		if source.Name != args.ExpectedSourceName || target.Name != args.ExpectedTargetName {
			return validationError("unit names no longer match the preview; run preview_unit_merge again"), nil, nil
		}
		res, err := unit.Merge(ctx, client, source, target)
		if err != nil {
			return apiError("merging units", err), nil, nil
		}
		b, _ := json.MarshalIndent(res, "", "  ")
		return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: string(b)}}}, nil, nil
	})
}

func validationError(message string) *mcp_sdk.CallToolResult {
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: "Error: " + message}}, IsError: true}
}
func apiError(action string, err error) *mcp_sdk.CallToolResult {
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: fmt.Sprintf("Error %s: %v", action, err)}}, IsError: true}
}
