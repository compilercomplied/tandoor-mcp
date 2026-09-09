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
	SourceURL   *string `json:"source_url,omitempty" jsonschema:"Replacement source URL."`
	Private     *bool   `json:"private,omitempty" jsonschema:"Replacement privacy setting."`
	KeywordIDs  *[]int  `json:"keyword_ids,omitempty" jsonschema:"Replacement keyword IDs; omit to keep existing keywords."`
	PropertyIDs *[]int  `json:"property_ids,omitempty" jsonschema:"Replacement property IDs; omit to keep existing properties."`
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
		if args.Name == nil && args.Description == nil && args.Servings == nil && args.WorkingTime == nil && args.WaitingTime == nil && args.SourceURL == nil && args.Private == nil && args.KeywordIDs == nil && args.PropertyIDs == nil {
			return validationError("at least one field to update is required"), nil, nil
		}
		var keywords, properties *[]api_update_recipe.IDRef
		if args.KeywordIDs != nil {
			values := make([]api_update_recipe.IDRef, len(*args.KeywordIDs))
			for i, id := range *args.KeywordIDs {
				values[i] = api_update_recipe.IDRef{ID: id}
			}
			keywords = &values
		}
		if args.PropertyIDs != nil {
			values := make([]api_update_recipe.IDRef, len(*args.PropertyIDs))
			for i, id := range *args.PropertyIDs {
				values[i] = api_update_recipe.IDRef{ID: id}
			}
			properties = &values
		}
		res, err := api_update_recipe.Update(ctx, client, args.RecipeID, api_update_recipe.Params{Name: args.Name, Description: args.Description, Servings: args.Servings, WorkingTime: args.WorkingTime, WaitingTime: args.WaitingTime, SourceURL: args.SourceURL, Private: args.Private, Keywords: keywords, Properties: properties})
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
