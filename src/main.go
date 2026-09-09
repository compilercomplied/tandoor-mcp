package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tools/add_category_to_supermarket"
	"github.com/compilercomplied/tandoor-mcp/src/tools/add_shopping_list_item"
	"github.com/compilercomplied/tandoor-mcp/src/tools/auto_plan"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_bookmarklet_import"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_cook_log"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_food"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_ingredient"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_inventory_entry"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_inventory_location"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_keyword"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_meal_plan"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_meal_type"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_property"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_property_type"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_recipe"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_recipe_import"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_storage"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_supermarket"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_supermarket_category"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_recipe"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_step"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_unit"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_unit_conversion"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_view_log"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_bookmarklet_imports"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_cook_logs"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_food_inherit_fields"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_foods"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_inventory_entries"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_inventory_locations"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_inventory_logs"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_keywords"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_meal_plans"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_meal_types"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_properties"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_property_types"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_recipe_details"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_recipe_imports"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_recipes"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_shopping_list"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_storages"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_supermarket_categories"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_supermarkets"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_unit_conversions"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_units"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_view_logs"
	"github.com/compilercomplied/tandoor-mcp/src/tools/merge_food"
	"github.com/compilercomplied/tandoor-mcp/src/tools/merge_unit"
	"github.com/compilercomplied/tandoor-mcp/src/tools/parse_ingredients"
	"github.com/compilercomplied/tandoor-mcp/src/tools/preview_food_merge"
	"github.com/compilercomplied/tandoor-mcp/src/tools/preview_unit_merge"
	"github.com/compilercomplied/tandoor-mcp/src/tools/remove_shopping_list_item"
	"github.com/compilercomplied/tandoor-mcp/src/tools/update_inventory_entry"
	"github.com/compilercomplied/tandoor-mcp/src/tools/update_shopping_list_item"
	mcp_sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	LoadEnv()

	logFormat := GetEnv("LOG_FORMAT")
	SetupLogging(logFormat)

	transportFlag := flag.String("transport", "sse", "Transport to use: 'sse' (HTTP) or 'stdio'")
	hostFlag := flag.String("host", "0.0.0.0", "Host to listen on (only for SSE transport)")
	portFlag := flag.String("port", "8080", "Port to listen on (only for SSE transport)")
	flag.Parse()

	apiURL := GetEnv("TANDOOR_API_URL")
	apiToken := GetEnv("TANDOOR_API_TOKEN")
	logHTTPBody := GetEnv("LOG_HTTP_BODY") == "true"

	client := tandoor.NewClient(apiURL, apiToken, logHTTPBody)

	server := mcp_sdk.NewServer(
		&mcp_sdk.Implementation{
			Name:    "tandoor-recipes-mcp",
			Version: "1.0.0",
		},
		nil,
	)

	get_recipes.Register(server, client)
	create_recipe.Register(server, client)
	create_tandoor_recipe.Register(server, client)
	get_recipe_details.Register(server, client)
	create_tandoor_step.Register(server, client)
	create_ingredient.Register(server, client)
	create_cook_log.Register(server, client)
	get_cook_logs.Register(server, client)
	create_view_log.Register(server, client)
	get_view_logs.Register(server, client)
	create_bookmarklet_import.Register(server, client)
	get_bookmarklet_imports.Register(server, client)
	parse_ingredients.Register(server, client)
	create_recipe_import.Register(server, client)
	get_recipe_imports.Register(server, client)
	create_meal_plan.Register(server, client)
	get_meal_plans.Register(server, client)
	auto_plan.Register(server, client)
	create_meal_type.Register(server, client)
	get_meal_types.Register(server, client)
	get_shopping_list.Register(server, client)
	add_shopping_list_item.Register(server, client)
	update_shopping_list_item.Register(server, client)
	remove_shopping_list_item.Register(server, client)
	get_supermarkets.Register(server, client)
	create_supermarket.Register(server, client)
	get_supermarket_categories.Register(server, client)
	create_supermarket_category.Register(server, client)
	add_category_to_supermarket.Register(server, client)
	create_storage.Register(server, client)
	get_storages.Register(server, client)
	get_inventory_locations.Register(server, client)
	create_inventory_location.Register(server, client)
	get_inventory_entries.Register(server, client)
	create_inventory_entry.Register(server, client)
	update_inventory_entry.Register(server, client)
	get_inventory_logs.Register(server, client)
	get_foods.Register(server, client)
	create_food.Register(server, client)
	get_food_inherit_fields.Register(server, client)
	get_keywords.Register(server, client)
	create_keyword.Register(server, client)
	get_units.Register(server, client)
	create_unit.Register(server, client)
	get_unit_conversions.Register(server, client)
	create_unit_conversion.Register(server, client)
	preview_food_merge.Register(server, client)
	merge_food.Register(server, client)
	preview_unit_merge.Register(server, client)
	merge_unit.Register(server, client)
	get_property_types.Register(server, client)
	create_property_type.Register(server, client)
	get_properties.Register(server, client)
	create_property.Register(server, client)

	switch *transportFlag {
	case "stdio":
		log.Println("Starting MCP server on stdio transport...")
		if err := server.Run(context.Background(), &mcp_sdk.StdioTransport{}); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	case "sse":
		addr := fmt.Sprintf("%s:%s", *hostFlag, *portFlag)
		log.Printf("Starting MCP server on SSE transport at http://%s ...", addr)
		log.Println("SSE Endpoint: /sse")
		log.Println("Message Endpoint: /sse?sessionid=<id> (automatically handled)")

		handler := mcp_sdk.NewSSEHandler(func(request *http.Request) *mcp_sdk.Server {
			return server
		}, nil)

		http.Handle("/sse", handler)
		http.Handle("/", handler)

		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	default:
		log.Fatalf("Unknown transport %q. Supported transports: stdio, sse", *transportFlag)
	}
}

// SetupLogging configures the logging format.
// If format is "json", it sets slog.JSONHandler as the default and redirects standard logs to it.
func SetupLogging(format string) {
	if strings.ToLower(format) == "json" {
		handler := slog.NewJSONHandler(os.Stderr, nil)
		slog.SetDefault(slog.New(handler))
		log.SetFlags(0)
		log.SetOutput(slog.NewLogLogger(handler, slog.LevelInfo).Writer())
	}
}
