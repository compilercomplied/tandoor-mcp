package preview_food_merge

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
	SourceFoodID int `json:"source_food_id" jsonschema:"The duplicate food ID to remove."`
	TargetFoodID int `json:"target_food_id" jsonschema:"The canonical food ID to keep."`
}

type Response struct {
	Source food.FoodResponse `json:"source"`
	Target food.FoodResponse `json:"target"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "preview_food_merge",
		Description: "Inspect a proposed food merge without changing data. Source is the duplicate to remove; target is the canonical food to keep.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing preview_food_merge. source=%d, target=%d", args.SourceFoodID, args.TargetFoodID)
		if args.SourceFoodID <= 0 || args.TargetFoodID <= 0 || args.SourceFoodID == args.TargetFoodID {
			return validationError("source_food_id and target_food_id must be different positive IDs"), nil, nil
		}
		source, err := food.Get(ctx, client, args.SourceFoodID)
		if err != nil {
			return apiError("retrieving source food", err), nil, nil
		}
		target, err := food.Get(ctx, client, args.TargetFoodID)
		if err != nil {
			return apiError("retrieving target food", err), nil, nil
		}
		b, _ := json.MarshalIndent(Response{Source: *source, Target: *target}, "", "  ")
		return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: string(b)}}}, nil, nil
	})
}

func validationError(message string) *mcp_sdk.CallToolResult {
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: "Error: " + message}}, IsError: true}
}
func apiError(action string, err error) *mcp_sdk.CallToolResult {
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: fmt.Sprintf("Error %s: %v", action, err)}}, IsError: true}
}
