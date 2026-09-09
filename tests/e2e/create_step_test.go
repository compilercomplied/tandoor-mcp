package e2e_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/ingredient"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/step"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_tandoor_step"
	"github.com/compilercomplied/tandoor-mcp/src/tools/update_tandoor_step"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestCreateStepE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath", func(t *testing.T) {
		// Arrange
		args := create_tandoor_step.Args{
			RecipeID:    1,
			Name:        "Prep step",
			Instruction: "Chop the onions",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)

		s := infra.ParseToolResponse[step.StepResponse](t, res)

		if s.ID == 0 {
			t.Errorf("expected step ID > 0, got 0")
		}
		if s.Name != "Prep step" {
			t.Errorf("expected Name='Prep step', got %q", s.Name)
		}
		if s.Instruction != "Chop the onions" {
			t.Errorf("expected Instruction='Chop the onions', got %q", s.Instruction)
		}
	})

	t.Run("HappyPath_WithIngredients", func(t *testing.T) {
		// Arrange
		tandoorClient := tandoor.NewClient(os.Getenv("TANDOOR_API_URL"), os.Getenv("TANDOOR_API_TOKEN"), false)
		nestedStep, err := step.Create(ctx, tandoorClient, 1, step.StepParam{
			Name:        "Initial Step",
			Instruction: "Mix spices",
			Ingredients: []ingredient.IngredientResponse{
				{
					Food:   ingredient.FoodRef{Name: "Coriander"},
					Unit:   ingredient.UnitRef{Name: "tbsp"},
					Amount: 1.5,
				},
			},
		})
		if err != nil {
			t.Fatalf("failed to create initial step with nested ingredient: %v", err)
		}

		if len(nestedStep.Ingredients) == 0 {
			t.Fatalf("expected created step to contain nested ingredient, got 0")
		}
		ingID := nestedStep.Ingredients[0].ID

		args := create_tandoor_step.Args{
			RecipeID:    1,
			Name:        "Prep step with ingredient",
			Instruction: "Chop the chicken",
			Ingredients: []int{ingID},
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)

		s := infra.ParseToolResponse[step.StepResponse](t, res)

		if s.ID == 0 {
			t.Errorf("expected step ID > 0, got 0")
		}
		if s.Name != "Prep step with ingredient" {
			t.Errorf("expected Name='Prep step with ingredient', got %q", s.Name)
		}
		if len(s.Ingredients) != 1 {
			t.Errorf("expected 1 ingredient associated, got %d", len(s.Ingredients))
		} else if s.Ingredients[0].ID != ingID {
			t.Errorf("expected associated ingredient ID to be %d, got %d", ingID, s.Ingredients[0].ID)
		}
	})

	t.Run("ValidationError_MissingInstructionAndName", func(t *testing.T) {
		// Arrange
		args := create_tandoor_step.Args{
			RecipeID: 1,
			// both Name and Instruction intentionally omitted
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true, got false")
		}

		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "instruction or name is required") {
			t.Errorf("expected validation message about instruction or name, got %q", errText)
		}
	})

	t.Run("ValidationError_InvalidRecipeID", func(t *testing.T) {
		// Arrange
		args := create_tandoor_step.Args{
			RecipeID:    999999,
			Instruction: "Chop the onions",
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true for non-existent recipe, got false")
		}

		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "Error creating step") {
			t.Errorf("expected error message to contain 'Error creating step', got %q", errText)
		}
	})
}

func TestUpdateTandoorStepE2E(t *testing.T) {
	// Arrange
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()
	createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_tandoor_step", create_tandoor_step.Args{
		RecipeID:    1,
		Name:        "Preparación",
		Instruction: "Corta las patatas",
	})
	infra.AssertToolSuccess(t, createRes, createErr)
	created := infra.ParseToolResponse[step.StepResponse](t, createRes)
	translatedName := "Preparation"
	translatedInstruction := "Slice the potatoes"

	// Act
	res, err := infra.CallTool(ctx, fixture.Client, "update_tandoor_step", update_tandoor_step.Args{
		StepID:      created.ID,
		Name:        &translatedName,
		Instruction: &translatedInstruction,
	})

	// Assert
	infra.AssertToolSuccess(t, res, err)
	updated := infra.ParseToolResponse[step.StepResponse](t, res)
	if updated.Name != translatedName {
		t.Errorf("expected name %q, got %q", translatedName, updated.Name)
	}
	if updated.Instruction != translatedInstruction {
		t.Errorf("expected instruction %q, got %q", translatedInstruction, updated.Instruction)
	}
}
