package e2e_test

import (
	"context"
	"testing"

	api_create_recipe "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/create_recipe"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
	api_update_recipe "github.com/compilercomplied/tandoor-mcp/src/tandoor/features/update_recipe"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_ingredient"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_recipe"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_recipe"
	"github.com/compilercomplied/tandoor-mcp/src/tools/update_recipe"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestCreateRecipeE2E(t *testing.T) {
	// Arrange
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	// Act
	res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_recipe", create_tandoor_recipe.Args{
		Name:        "Test Recipe 1",
		Description: "A tasty recipe",
	})

	// Assert
	infra.AssertToolSuccess(t, res, err)

	recipe := infra.ParseToolResponse[api_create_recipe.RecipeResponse](t, res)

	if recipe.Name != "Test Recipe 1" {
		t.Errorf("expected name 'Test Recipe 1', got %q", recipe.Name)
	}
	if recipe.Description != "A tasty recipe" {
		t.Errorf("expected description 'A tasty recipe', got %q", recipe.Description)
	}
}

func TestUpdateRecipeE2E(t *testing.T) {
	// Arrange
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()
	createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_tandoor_recipe", create_tandoor_recipe.Args{
		Name:        "Tortilla española",
		Description: "Patatas y huevo",
	})
	infra.AssertToolSuccess(t, createRes, createErr)
	created := infra.ParseToolResponse[api_create_recipe.RecipeResponse](t, createRes)
	translatedName := "Spanish omelette"
	translatedDescription := "Potatoes and eggs"
	sourceURL := "https://example.com/tortilla"
	private := true

	// Act
	res, err := infra.CallTool(ctx, fixture.Client, "update_recipe", update_recipe.Args{
		RecipeID:    created.ID,
		Name:        &translatedName,
		Description: &translatedDescription,
		SourceURL:   &sourceURL,
		Private:     &private,
	})

	// Assert
	infra.AssertToolSuccess(t, res, err)
	updated := infra.ParseToolResponse[api_update_recipe.Response](t, res)
	if updated.Name != translatedName {
		t.Errorf("expected name %q, got %q", translatedName, updated.Name)
	}
	if updated.Description != translatedDescription {
		t.Errorf("expected description %q, got %q", translatedDescription, updated.Description)
	}
	if updated.SourceURL == nil || *updated.SourceURL != sourceURL || updated.Private != private {
		t.Errorf("expected source URL %q and private %v, got %v and %v", sourceURL, private, updated.SourceURL, updated.Private)
	}
}

func TestCreateRecipeWithStepsAndIngredientsE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	recipeID := 1

	// 1. Create two ingredients associated with recipe 1 first (so they are not floating)
	ingRes1, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", create_ingredient.Args{
		FoodName: "Flour",
		UnitName: "grams",
		Amount:   500,
		RecipeID: &recipeID,
	})
	infra.AssertToolSuccess(t, ingRes1, err)
	ing1 := infra.ParseToolResponse[ingredient.IngredientResponse](t, ingRes1)

	ingRes2, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", create_ingredient.Args{
		FoodName: "Water",
		UnitName: "ml",
		Amount:   300,
		RecipeID: &recipeID,
	})
	infra.AssertToolSuccess(t, ingRes2, err)
	ing2 := infra.ParseToolResponse[ingredient.IngredientResponse](t, ingRes2)

	// 2. Call create_recipe with steps linking these ingredients
	recipeArgs := create_recipe.Args{
		Name:        "Sourdough Bread",
		Description: "Simple homemade sourdough",
		Steps: []create_recipe.StepParam{
			{
				Name:        "Mix",
				Instruction: "Mix flour and water together",
				Ingredients: []int{ing1.ID, ing2.ID},
			},
		},
	}

	res, err := infra.CallTool(ctx, fixture.Client, "create_recipe", recipeArgs)
	infra.AssertToolSuccess(t, res, err)

	recipe := infra.ParseToolResponse[api_create_recipe.RecipeResponse](t, res)

	if recipe.Name != "Sourdough Bread" {
		t.Errorf("expected name 'Sourdough Bread', got %q", recipe.Name)
	}
	if len(recipe.Steps) != 1 {
		t.Fatalf("expected 1 step in response, got %d", len(recipe.Steps))
	}

	step := recipe.Steps[0]
	if step.Name != "Mix" {
		t.Errorf("expected step name 'Mix', got %q", step.Name)
	}
	if len(step.Ingredients) != 2 {
		t.Fatalf("expected 2 ingredients in step, got %d", len(step.Ingredients))
	}
	if step.Ingredients[0].ID != ing1.ID && step.Ingredients[1].ID != ing1.ID {
		t.Errorf("expected ingredient ID %d to be in the step", ing1.ID)
	}
}
