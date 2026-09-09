package e2e_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/step"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_ingredient"
	"github.com/compilercomplied/tandoor-mcp/src/tools/update_ingredient"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestCreateIngredientE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			FoodName: "Onion",
			UnitName: "pieces",
			Amount:   2,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)

		ing := infra.ParseToolResponse[ingredient.IngredientResponse](t, res)

		if ing.ID == 0 {
			t.Errorf("expected ingredient ID > 0, got 0")
		}
		if ing.Food.Name != "Onion" {
			t.Errorf("expected Food.Name='Onion', got %q", ing.Food.Name)
		}
		if ing.Unit.Name != "pieces" {
			t.Errorf("expected Unit.Name='pieces', got %q", ing.Unit.Name)
		}
		if ing.Amount != 2 {
			t.Errorf("expected Amount=2, got %v", ing.Amount)
		}
	})

	t.Run("HappyPath_UpdateIngredient", func(t *testing.T) {
		// Arrange
		recipeID := 1
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_ingredient", create_ingredient.Args{FoodName: "Tomato", UnitName: "g", Amount: 100, RecipeID: &recipeID})
		infra.AssertToolSuccess(t, createRes, createErr)
		created := infra.ParseToolResponse[ingredient.IngredientResponse](t, createRes)
		amount := 250.0
		note := "diced"
		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "update_ingredient", update_ingredient.Args{IngredientID: created.ID, Amount: &amount, Note: &note})
		// Assert
		infra.AssertToolSuccess(t, res, err)
		updated := infra.ParseToolResponse[ingredient.IngredientResponse](t, res)
		if updated.Amount != amount || updated.Note != note {
			t.Errorf("expected amount %v and note %q, got %v and %q", amount, note, updated.Amount, updated.Note)
		}
	})

	t.Run("HappyPath_WithNote", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			FoodName: "Garlic",
			UnitName: "cloves",
			Amount:   3,
			Note:     "finely chopped",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)

		ing := infra.ParseToolResponse[ingredient.IngredientResponse](t, res)

		if ing.Note != "finely chopped" {
			t.Errorf("expected Note='finely chopped', got %q", ing.Note)
		}
	})

	t.Run("ValidationError_MissingFoodName", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			UnitName: "pieces",
			Amount:   1,
			// FoodName intentionally omitted
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "food_name is required") {
			t.Errorf("expected validation message about food_name, got %q", errText)
		}
	})

	t.Run("ValidationError_MissingUnitName", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			FoodName: "Salt",
			Amount:   1,
			// UnitName intentionally omitted
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "unit_name is required") {
			t.Errorf("expected validation message about unit_name, got %q", errText)
		}
	})

	t.Run("ValidationError_ZeroAmount", func(t *testing.T) {
		// Arrange
		args := create_ingredient.Args{
			FoodName: "Pepper",
			UnitName: "grams",
			Amount:   0, // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true for zero amount, got false")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "amount must be greater than 0") {
			t.Errorf("expected validation message about amount, got %q", errText)
		}
	})

	t.Run("HappyPath_WithRecipeID", func(t *testing.T) {
		recipeID := 1
		// Arrange
		args := create_ingredient.Args{
			FoodName: "Carrot",
			UnitName: "pieces",
			Amount:   4,
			RecipeID: &recipeID,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_ingredient", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)

		ing := infra.ParseToolResponse[ingredient.IngredientResponse](t, res)

		if ing.ID == 0 {
			t.Errorf("expected ingredient ID > 0, got 0")
		}
		if ing.Food.Name != "Carrot" {
			t.Errorf("expected Food.Name='Carrot', got %q", ing.Food.Name)
		}

		// 1. Verify we can GET the ingredient (since it's associated with a recipe, it won't 404)
		tandoorClient := tandoor.NewClient(os.Getenv("TANDOOR_API_URL"), os.Getenv("TANDOOR_API_TOKEN"), false)
		fetchedIng, err := ingredient.Get(ctx, tandoorClient, ing.ID)
		if err != nil {
			t.Fatalf("failed to GET non-floating ingredient: %v", err)
		}
		if fetchedIng.ID != ing.ID {
			t.Errorf("expected fetched ingredient ID %d, got %d", ing.ID, fetchedIng.ID)
		}

		// 2. Fetch recipe details and verify the ingredient is in the first step
		recipeStepsRes, err := tandoor.Request[step.RecipeSteps](ctx, tandoorClient, "GET", "/api/recipe/1/", nil, nil)
		if err != nil {
			t.Fatalf("failed to GET recipe steps: %v", err)
		}
		if len(recipeStepsRes.Steps) == 0 {
			t.Fatalf("expected recipe to have at least one step")
		}

		found := false
		for _, stepIng := range recipeStepsRes.Steps[0].Ingredients {
			if stepIng.ID == ing.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected ingredient %d to be in the first step of recipe %d", ing.ID, recipeID)
		}
	})
}
