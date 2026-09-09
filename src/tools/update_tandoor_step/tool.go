package update_tandoor_step

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	api_ingredient "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
	api_step "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/step"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	StepID               int     `json:"step_id" jsonschema:"ID of the recipe step to update."`
	Name                 *string `json:"name,omitempty" jsonschema:"Replacement step heading."`
	Instruction          *string `json:"instruction,omitempty" jsonschema:"Replacement instruction text."`
	Time                 *int    `json:"time,omitempty" jsonschema:"Replacement time in minutes."`
	Order                *int    `json:"order,omitempty" jsonschema:"Replacement ordering value."`
	Ingredients          *[]int  `json:"ingredients,omitempty" jsonschema:"Replacement list of ingredient IDs. Omit to keep existing ingredients; use [] to clear them."`
	ShowAsHeader         *bool   `json:"show_as_header,omitempty" jsonschema:"Whether to display the step as a header."`
	ShowIngredientsTable *bool   `json:"show_ingredients_table,omitempty" jsonschema:"Whether to display its ingredients as a table."`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "update_tandoor_step",
		Description: "Partially update an existing recipe step. Omitted fields remain unchanged; use this to replace translated step text without replacing its ingredients.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing update_tandoor_step. step_id=%d", args.StepID)
		if args.StepID <= 0 {
			return validationError("step_id must be a positive ID"), nil, nil
		}
		if args.Name == nil && args.Instruction == nil && args.Time == nil && args.Order == nil && args.Ingredients == nil && args.ShowAsHeader == nil && args.ShowIngredientsTable == nil {
			return validationError("at least one field to update is required"), nil, nil
		}

		var ingredients *[]api_ingredient.IngredientResponse
		if args.Ingredients != nil {
			resolved := make([]api_ingredient.IngredientResponse, len(*args.Ingredients))
			for i, ingredientID := range *args.Ingredients {
				ingredient, err := api_ingredient.Get(ctx, client, ingredientID)
				if err != nil {
					return apiError(fmt.Errorf("fetching ingredient %d: %w", ingredientID, err)), nil, nil
				}
				resolved[i] = *ingredient
			}
			ingredients = &resolved
		}

		res, err := api_step.Update(ctx, client, args.StepID, api_step.UpdateStepParam{
			Name: args.Name, Instruction: args.Instruction, Time: args.Time, Order: args.Order, Ingredients: ingredients,
			ShowAsHeader: args.ShowAsHeader, ShowIngredientsTable: args.ShowIngredientsTable,
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
	return &mcp_sdk.CallToolResult{Content: []mcp_sdk.Content{&mcp_sdk.TextContent{Text: fmt.Sprintf("Error updating step: %v", err)}}, IsError: true}
}
