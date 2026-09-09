package update_recipe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_update_recipe "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/update_recipe"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	RecipeID    int     `json:"recipe_id" jsonschema:"ID of the recipe to update."`
	Name        *string `json:"name,omitempty" jsonschema:"Replacement recipe name."`
	Description *string `json:"description,omitempty" jsonschema:"Replacement recipe description."`
	Servings    *int    `json:"servings,omitempty" jsonschema:"Replacement number of servings."`
	WorkingTime *int    `json:"working_time,omitempty" jsonschema:"Replacement working time in minutes."`
	WaitingTime *int    `json:"waiting_time,omitempty" jsonschema:"Replacement waiting time in minutes."`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "update_recipe",
		Description: "Partially update an existing recipe. Omitted fields remain unchanged; use this to replace translated recipe text without replacing steps or ingredients.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing update_recipe. recipe_id=%d", args.RecipeID)
		if args.RecipeID <= 0 {
			return validationError("recipe_id must be a positive ID"), nil, nil
		}
		if args.Name == nil && args.Description == nil && args.Servings == nil && args.WorkingTime == nil && args.WaitingTime == nil {
			return validationError("at least one field to update is required"), nil, nil
		}
		res, err := api_update_recipe.Update(ctx, client, args.RecipeID, api_update_recipe.Params{
			Name: args.Name, Description: args.Description, Servings: args.Servings, WorkingTime: args.WorkingTime, WaitingTime: args.WaitingTime,
		})
		if err != nil {
			return apiError(err), nil, nil
		}
		b, _ := json.MarshalIndent(res, "", "  ")
		return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: string(b)}}}, nil, nil
	})
}

func validationError(message string) *mcp_sdk.CallToolResult {
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: "Error: " + message}}, IsError: true}
}

func apiError(err error) *mcp_sdk.CallToolResult {
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: fmt.Sprintf("Error updating recipe: %v", err)}}, IsError: true}
}
