package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/food"
	"github.com/compilercomplied/tandoor-mcp/src/tandoor/features/unit"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_food"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_unit"
	"github.com/compilercomplied/tandoor-mcp/src/tools/create_unit_conversion"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_unit_conversions"
	"github.com/compilercomplied/tandoor-mcp/src/tools/get_units"
	"github.com/compilercomplied/tandoor-mcp/src/tools/merge_unit"
	"github.com/compilercomplied/tandoor-mcp/src/tools/preview_unit_merge"
	"github.com/compilercomplied/tandoor-mcp/tests/e2e/infra"
)

func TestUnitE2E(t *testing.T) {
	ctx := context.Background()
	defer infra.PurgeAndSeedDatabase()

	t.Run("HappyPath_CreateUnit", func(t *testing.T) {
		// Arrange
		name := "Gram"
		pluralName := "Grams"
		desc := "Metric weight unit"

		args := create_unit.Args{
			Name:        name,
			PluralName:  &pluralName,
			Description: &desc,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_unit", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		u := infra.ParseToolResponse[unit.UnitResponse](t, res)

		if u.ID == 0 {
			t.Errorf("expected unit ID > 0, got 0")
		}
		if u.Name != name {
			t.Errorf("expected name %q, got %q", name, u.Name)
		}
		if u.PluralName == nil {
			t.Errorf("expected plural name to be non-nil")
		} else if *u.PluralName != pluralName && *u.PluralName != name {
			t.Errorf("expected plural name %q or %q, got %q", pluralName, name, *u.PluralName)
		}

		if u.Description == nil {
			t.Errorf("expected description %q, got nil", desc)
		} else if *u.Description != desc {
			t.Errorf("expected description %q, got %q", desc, *u.Description)
		}
	})

	t.Run("HappyPath_GetUnits", func(t *testing.T) {
		// Arrange: Create a unit first
		name := "Kilogram"
		args := create_unit.Args{
			Name: name,
		}
		createRes, createErr := infra.CallTool(ctx, fixture.Client, "create_unit", args)
		infra.AssertToolSuccess(t, createRes, createErr)
		u := infra.ParseToolResponse[unit.UnitResponse](t, createRes)

		// Act
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_units", get_units.Args{
			Query: &name,
		})

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		list := infra.ParseToolResponse[unit.UnitListResponse](t, listRes)

		found := false
		for _, item := range list.Results {
			if item.ID == u.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find unit ID %d in get_units response", u.ID)
		}
	})

	t.Run("HappyPath_PreviewAndMergeUnit", func(t *testing.T) {
		// Arrange
		targetRes, targetErr := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "Tablespoon"})
		infra.AssertToolSuccess(t, targetRes, targetErr)
		target := infra.ParseToolResponse[unit.UnitResponse](t, targetRes)
		sourceRes, sourceErr := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "tbsp."})
		infra.AssertToolSuccess(t, sourceRes, sourceErr)
		source := infra.ParseToolResponse[unit.UnitResponse](t, sourceRes)

		// Act
		previewRes, previewErr := infra.CallTool(ctx, fixture.Client, "preview_unit_merge", preview_unit_merge.Args{
			SourceUnitID: source.ID,
			TargetUnitID: target.ID,
		})

		// Assert
		infra.AssertToolSuccess(t, previewRes, previewErr)
		preview := infra.ParseToolResponse[preview_unit_merge.Response](t, previewRes)
		if preview.Source.Name != source.Name || preview.Target.Name != target.Name {
			t.Fatalf("expected preview %q -> %q, got %q -> %q", source.Name, target.Name, preview.Source.Name, preview.Target.Name)
		}

		// Act
		mergeRes, mergeErr := infra.CallTool(ctx, fixture.Client, "merge_unit", merge_unit.Args{
			SourceUnitID:       source.ID,
			TargetUnitID:       target.ID,
			ExpectedSourceName: source.Name,
			ExpectedTargetName: target.Name,
		})

		// Assert
		infra.AssertToolSuccess(t, mergeRes, mergeErr)
	})

	t.Run("ValidationError_StaleMergePreview", func(t *testing.T) {
		// Arrange
		targetRes, targetErr := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "Canonical tablespoon"})
		infra.AssertToolSuccess(t, targetRes, targetErr)
		target := infra.ParseToolResponse[unit.UnitResponse](t, targetRes)
		sourceRes, sourceErr := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "stale tbsp"})
		infra.AssertToolSuccess(t, sourceRes, sourceErr)
		source := infra.ParseToolResponse[unit.UnitResponse](t, sourceRes)

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "merge_unit", merge_unit.Args{
			SourceUnitID:       source.ID,
			TargetUnitID:       target.ID,
			ExpectedSourceName: "outdated unit name",
			ExpectedTargetName: target.Name,
		})

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatal("expected a stale merge preview to be rejected")
		}
		if !strings.Contains(infra.ExtractErrorText(t, res), "run preview_unit_merge again") {
			t.Fatalf("expected stale preview error, got %q", infra.ExtractErrorText(t, res))
		}
	})

	t.Run("HappyPath_CreateUnitConversionWithoutFood", func(t *testing.T) {
		// Arrange: Create base and converted units
		baseRes, err1 := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "Milliliter"})
		infra.AssertToolSuccess(t, baseRes, err1)
		baseUnit := infra.ParseToolResponse[unit.UnitResponse](t, baseRes)

		convRes, err2 := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "Liter"})
		infra.AssertToolSuccess(t, convRes, err2)
		convUnit := infra.ParseToolResponse[unit.UnitResponse](t, convRes)

		args := create_unit_conversion.Args{
			BaseAmount:      1000,
			BaseUnitID:      baseUnit.ID,
			ConvertedAmount: 1,
			ConvertedUnitID: convUnit.ID,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_unit_conversion", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		uc := infra.ParseToolResponse[unit.UnitConversionResponse](t, res)

		if uc.ID == 0 {
			t.Errorf("expected unit conversion ID > 0, got 0")
		}
		if uc.BaseAmount != 1000 {
			t.Errorf("expected base amount 1000, got %f", uc.BaseAmount)
		}
		if uc.BaseUnit.ID != baseUnit.ID {
			t.Errorf("expected base unit ID %d, got %d", baseUnit.ID, uc.BaseUnit.ID)
		}
		if uc.ConvertedAmount != 1 {
			t.Errorf("expected converted amount 1, got %f", uc.ConvertedAmount)
		}
		if uc.ConvertedUnit.ID != convUnit.ID {
			t.Errorf("expected converted unit ID %d, got %d", convUnit.ID, uc.ConvertedUnit.ID)
		}
	})

	t.Run("HappyPath_CreateUnitConversionWithFood", func(t *testing.T) {
		// Arrange: Create units and food
		baseRes, err1 := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "Cup"})
		infra.AssertToolSuccess(t, baseRes, err1)
		baseUnit := infra.ParseToolResponse[unit.UnitResponse](t, baseRes)

		convRes, err2 := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "Ounce"})
		infra.AssertToolSuccess(t, convRes, err2)
		convUnit := infra.ParseToolResponse[unit.UnitResponse](t, convRes)

		foodRes, err3 := infra.CallTool(ctx, fixture.Client, "create_food", create_food.Args{Name: "Flour"})
		infra.AssertToolSuccess(t, foodRes, err3)
		fd := infra.ParseToolResponse[food.FoodResponse](t, foodRes)

		args := create_unit_conversion.Args{
			BaseAmount:      1,
			BaseUnitID:      baseUnit.ID,
			ConvertedAmount: 4.25,
			ConvertedUnitID: convUnit.ID,
			FoodID:          &fd.ID,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_unit_conversion", args)

		// Assert
		infra.AssertToolSuccess(t, res, err)
		uc := infra.ParseToolResponse[unit.UnitConversionResponse](t, res)

		if uc.ID == 0 {
			t.Errorf("expected unit conversion ID > 0, got 0")
		}
		if uc.Food == nil || uc.Food.ID != fd.ID {
			t.Errorf("expected food ID %d, got %v", fd.ID, uc.Food)
		}
	})

	t.Run("HappyPath_GetUnitConversions", func(t *testing.T) {
		// Arrange: Create base and converted units
		baseRes, err1 := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "Pinch"})
		infra.AssertToolSuccess(t, baseRes, err1)
		baseUnit := infra.ParseToolResponse[unit.UnitResponse](t, baseRes)

		convRes, err2 := infra.CallTool(ctx, fixture.Client, "create_unit", create_unit.Args{Name: "Teaspoon"})
		infra.AssertToolSuccess(t, convRes, err2)
		convUnit := infra.ParseToolResponse[unit.UnitResponse](t, convRes)

		args := create_unit_conversion.Args{
			BaseAmount:      16,
			BaseUnitID:      baseUnit.ID,
			ConvertedAmount: 1,
			ConvertedUnitID: convUnit.ID,
		}
		createRes, err3 := infra.CallTool(ctx, fixture.Client, "create_unit_conversion", args)
		infra.AssertToolSuccess(t, createRes, err3)
		uc := infra.ParseToolResponse[unit.UnitConversionResponse](t, createRes)

		// Act
		listRes, listErr := infra.CallTool(ctx, fixture.Client, "get_unit_conversions", get_unit_conversions.Args{})

		// Assert
		infra.AssertToolSuccess(t, listRes, listErr)
		list := infra.ParseToolResponse[unit.UnitConversionListResponse](t, listRes)

		found := false
		for _, item := range list.Results {
			if item.ID == uc.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find unit conversion ID %d in get_unit_conversions response", uc.ID)
		}
	})

	t.Run("ValidationError_MissingName", func(t *testing.T) {
		// Arrange
		args := create_unit.Args{
			Name: "", // invalid
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_unit", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true for empty name")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "name is required") {
			t.Errorf("expected error message to contain 'name is required', got %q", errText)
		}
	})

	t.Run("ValidationError_ZeroAmount", func(t *testing.T) {
		// Arrange
		args := create_unit_conversion.Args{
			BaseAmount:      0, // invalid
			BaseUnitID:      1,
			ConvertedAmount: 1,
			ConvertedUnitID: 2,
		}

		// Act
		res, err := infra.CallTool(ctx, fixture.Client, "create_unit_conversion", args)

		// Assert
		if err != nil {
			t.Fatalf("unexpected transport error: %v", err)
		}
		if !res.IsError {
			t.Fatalf("expected IsError=true for zero amount")
		}
		errText := infra.ExtractErrorText(t, res)
		if !strings.Contains(errText, "base_amount must be greater than 0") {
			t.Errorf("expected error message to contain 'base_amount must be greater than 0', got %q", errText)
		}
	})
}
