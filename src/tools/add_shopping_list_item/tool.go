package add_shopping_list_item

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/shoppinglist"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Args struct {
	FoodNameOrID string    `json:"food_name_or_id"`
	Amount       string `json:"amount"`
	UnitNameOrID string    `json:"unit_name_or_id"`
	Note         string `json:"note,omitempty"`
}

func Register(server *mcp_sdk.Server, client *tandoor.Client) {
	mcp_sdk.AddTool(server, &mcp_sdk.Tool{
		Name:        "add_shopping_list_item",
		Description: "Add an item to the shopping list, allowing food/unit names or IDs.",
	}, func(ctx context.Context, req *mcp_sdk.CallToolRequest, args Args) (*mcp_sdk.CallToolResult, any, error) {
		log.Printf("Executing add_shopping_list_item")

		amt, err := parseAmount(args.Amount)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error parsing amount: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		foodID, foodName := parseNameOrID(args.FoodNameOrID)
		if foodID == 0 && foodName == "" {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: "Error: food_name_or_id is required and cannot be empty"},
				},
				IsError: true,
			}, nil, nil
		}
		var foodParam shoppinglist.FoodShopping
		if foodID > 0 {
			foodParam.ID = foodID
		} else {
			foodParam.Name = foodName
		}

		unitID, unitName := parseNameOrID(args.UnitNameOrID)
		var unitParam *shoppinglist.UnitResponse
		if unitID > 0 || unitName != "" {
			unitParam = &shoppinglist.UnitResponse{}
			if unitID > 0 {
				unitParam.ID = unitID
			} else {
				unitParam.Name = unitName
			}
		}

		list, err := shoppinglist.GetOrCreateDefaultList(ctx, client)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error checking/creating shopping list: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		entryParams := shoppinglist.CreateEntryParam{
			ShoppingLists: []shoppinglist.ShoppingListID{
				{ID: list.ID},
			},
			Food:   foodParam,
			Unit:   unitParam,
			Amount: amt,
		}

		res, err := shoppinglist.CreateEntry(ctx, client, entryParams)
		if err != nil {
			return &mcp_sdk.CallToolResult{
				Content: []mcp_sdk.Content{
					&mcp_sdk.TextContent{Text: fmt.Sprintf("Error creating shopping list item: %v", err)},
				},
				IsError: true,
			}, nil, nil
		}

		b, _ := json.MarshalIndent(res, "", "  ")
		return &mcp_sdk.CallToolResult{
			Content: []mcp_sdk.Content{
				&mcp_sdk.TextContent{Text: string(b)},
			},
		}, nil, nil
	})
}

func parseAmount(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("amount is empty")
	}
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val, nil
	}
	if parts := strings.Split(s, "/"); len(parts) == 2 {
		num, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		den, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err1 == nil && err2 == nil && den != 0 {
			return num / den, nil
		}
	}
	return 0, fmt.Errorf("invalid amount format: %q", s)
}

func parseNameOrID(val string) (int, string) {
	s := strings.TrimSpace(val)
	if s == "" {
		return 0, ""
	}
	if id, err := strconv.Atoi(s); err == nil {
		return id, ""
	}
	return 0, s
}
