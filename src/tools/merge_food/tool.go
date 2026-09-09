package merge_food

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/food"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	SourceFoodID       int    `json:"source_food_id" jsonschema:"The duplicate food ID to remove."`
	TargetFoodID       int    `json:"target_food_id" jsonschema:"The canonical food ID to keep."`
	ExpectedSourceName string `json:"expected_source_name" jsonschema:"Source name returned by preview_food_merge; merge stops if it no longer matches."`
	ExpectedTargetName string `json:"expected_target_name" jsonschema:"Target name returned by preview_food_merge; merge stops if it no longer matches."`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "merge_food",
		Description: "Merge a duplicate food into a canonical food. Run preview_food_merge first and pass its exact source and target names to guard against stale IDs. This removes the source food.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing merge_food. source=%d, target=%d", args.SourceFoodID, args.TargetFoodID)
		if args.SourceFoodID <= 0 || args.TargetFoodID <= 0 || args.SourceFoodID == args.TargetFoodID || args.ExpectedSourceName == "" || args.ExpectedTargetName == "" {
			return validationError("source_food_id and target_food_id must be different positive IDs and expected names are required"), nil, nil
		}
		source, err := food.Get(ctx, client, args.SourceFoodID)
		if err != nil {
			return apiError("retrieving source food", err), nil, nil
		}
		target, err := food.Get(ctx, client, args.TargetFoodID)
		if err != nil {
			return apiError("retrieving target food", err), nil, nil
		}
		if source.Name != args.ExpectedSourceName || target.Name != args.ExpectedTargetName {
			return validationError("food names no longer match the preview; run preview_food_merge again"), nil, nil
		}
		res, err := food.Merge(ctx, client, source, target)
		if err != nil {
			return apiError("merging foods", err), nil, nil
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
