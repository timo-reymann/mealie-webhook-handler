package mealie

import (
	"os"
	"testing"
	"time"
)

func TestParseWebhook(t *testing.T) {
	// Read the JSON file
	content, err := os.ReadFile("testdata/webhook_payload.json")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal the JSON into the Recipe struct
	recipe, err := ParseWebhook(content)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Assertions for the Action field
	if recipe.Action.ActionType != "post" {
		t.Errorf("Expected ActionType 'post', got %s", recipe.Action.ActionType)
	}
	if recipe.Action.Title != "foo" {
		t.Errorf("Expected Title 'foo', got %s", recipe.Action.Title)
	}
	if recipe.Action.URL != "https://webhook.site/2b55af7c-1a7a-4cfa-b6b6-3f5ed65b8e3d" {
		t.Errorf("Expected URL 'https://webhook.site/2b55af7c-1a7a-4cfa-b6b6-3f5ed65b8e3d', got %s", recipe.Action.URL)
	}
	if recipe.Action.GroupID != "5e582dd1-3eef-4389-a031-53e077e696b8" {
		t.Errorf("Expected GroupID '5e582dd1-3eef-4389-a031-53e077e696b8', got %s", recipe.Action.GroupID)
	}
	if recipe.Action.HouseholdID != "49fce3b4-ebb0-49ee-b976-e8a0bc17c81e" {
		t.Errorf("Expected HouseholdID '49fce3b4-ebb0-49ee-b976-e8a0bc17c81e', got %s", recipe.Action.HouseholdID)
	}
	if recipe.Action.ID != "94078c56-4c4d-4baf-a88b-9118f153cc61" {
		t.Errorf("Expected ID '94078c56-4c4d-4baf-a88b-9118f153cc61', got %s", recipe.Action.ID)
	}

	// Assertions for the RecipeContent field
	if recipe.Content.ID != "40c91397-17cb-4239-b582-bd0be7b82e48" {
		t.Errorf("Expected RecipeContent.ID '40c91397-17cb-4239-b582-bd0be7b82e48', got %s", recipe.Content.ID)
	}
	if recipe.Content.UserID != "7c2f1e7c-8384-45aa-b3cc-849317f13993" {
		t.Errorf("Expected UserID '7c2f1e7c-8384-45aa-b3cc-849317f13993', got %s", recipe.Content.UserID)
	}
	if recipe.Content.HouseholdID != "49fce3b4-ebb0-49ee-b976-e8a0bc17c81e" {
		t.Errorf("Expected HouseholdID '49fce3b4-ebb0-49ee-b976-e8a0bc17c81e', got %s", recipe.Content.HouseholdID)
	}
	if recipe.Content.GroupID != "5e582dd1-3eef-4389-a031-53e077e696b8" {
		t.Errorf("Expected GroupID '5e582dd1-3eef-4389-a031-53e077e696b8', got %s", recipe.Content.GroupID)
	}
	if recipe.Content.Name != "Hühnchen-Wok mit Teriyaki-Erdnussbuttersoße" {
		t.Errorf("Expected Name 'Hühnchen-Wok mit Teriyaki-Erdnussbuttersoße', got %s", recipe.Content.Name)
	}
	if recipe.Content.Slug != "huhnchen-wok-mit-teriyaki-erdnussbuttersosse" {
		t.Errorf("Expected Slug 'huhnchen-wok-mit-teriyaki-erdnussbuttersosse', got %s", recipe.Content.Slug)
	}
	if recipe.Content.Image != "38" {
		t.Errorf("Expected Image '38', got %s", recipe.Content.Image)
	}
	if recipe.Content.RecipeServings != 0.0 {
		t.Errorf("Expected RecipeServings 0.0, got %f", recipe.Content.RecipeServings)
	}
	if recipe.Content.RecipeYieldQuantity != 0.0 {
		t.Errorf("Expected RecipeYieldQuantity 0.0, got %f", recipe.Content.RecipeYieldQuantity)
	}
	if recipe.Content.RecipeYield != "" {
		t.Errorf("Expected RecipeYield '', got %s", recipe.Content.RecipeYield)
	}
	if recipe.Content.Description != "Schnell gemacht, man braucht nur einen Wok, die Zutaten und ca. 30 Minuten. Dabei schmeckt es auch noch lecker und ist halbwegs gesund ;)" {
		t.Errorf("Expected Description 'Schnell gemacht, man braucht nur einen Wok, die Zutaten und ca. 30 Minuten. Dabei schmeckt es auch noch lecker und ist halbwegs gesund ;)', got %s", recipe.Content.Description)
	}
	if recipe.Content.OrgURL != "https://recipes.timo-reymann.de/recipes/Huehnchen-Wok-mit-Teriyaki-Erdnussbuttersosse.html" {
		t.Errorf("Expected OrgURL 'https://recipes.timo-reymann.de/recipes/Huehnchen-Wok-mit-Teriyaki-Erdnussbuttersosse.html', got %s", recipe.Content.OrgURL)
	}
	if recipe.Content.DateAdded != "2025-09-15" {
		t.Errorf("Expected DateAdded '2025-09-15', got %s", recipe.Content.DateAdded)
	}

	// Parse and compare timestamps
	expectedDateUpdated, _ := time.Parse(time.RFC3339, "2025-09-15T10:45:30.670639+00:00")
	if !recipe.Content.DateUpdated.Equal(expectedDateUpdated) {
		t.Errorf("Expected DateUpdated '2025-09-15T10:45:30.670639+00:00', got %v", recipe.Content.DateUpdated)
	}

	expectedCreatedAt, _ := time.Parse(time.RFC3339, "2025-09-15T10:39:28.051067+00:00")
	if !recipe.Content.CreatedAt.Equal(expectedCreatedAt) {
		t.Errorf("Expected CreatedAt '2025-09-15T10:39:28.051067+00:00', got %v", recipe.Content.CreatedAt)
	}

	expectedUpdatedAt, _ := time.Parse(time.RFC3339, "2025-09-15T10:45:30.674365+00:00")
	if !recipe.Content.UpdatedAt.Equal(expectedUpdatedAt) {
		t.Errorf("Expected UpdatedAt '2025-09-15T10:45:30.674365+00:00', got %v", recipe.Content.UpdatedAt)
	}

	// Assertions for RecipeIngredients
	if len(recipe.Content.RecipeIngredient) != 15 {
		t.Fatalf("Expected 15 RecipeIngredients, got %d", len(recipe.Content.RecipeIngredient))
	}

	// Example assertion for the first ingredient
	firstIngredient := recipe.Content.RecipeIngredient[0]
	if *firstIngredient.Quantity != 500.0 {
		t.Errorf("Expected Quantity 500.0, got %f", *firstIngredient.Quantity)
	}
	if firstIngredient.Unit.Name != "Gramm" {
		t.Errorf("Expected Unit.Name 'Gramm', got %s", firstIngredient.Unit.Name)
	}
	if firstIngredient.Food.Name != "Hühnchen" {
		t.Errorf("Expected Food.Name 'Hühnchen', got %s", firstIngredient.Food.Name)
	}
	if firstIngredient.Display != "500 Gramm Hühnchen" {
		t.Errorf("Expected Display '500 Gramm Hühnchen', got %s", firstIngredient.Display)
	}

	// Assertions for RecipeInstructions
	if len(recipe.Content.RecipeInstructions) != 10 {
		t.Fatalf("Expected 10 RecipeInstructions, got %d", len(recipe.Content.RecipeInstructions))
	}

	// Example assertion for the first instruction
	firstInstruction := recipe.Content.RecipeInstructions[0]
	if firstInstruction.Text != "Reis parallel kochen lassen" {
		t.Errorf("Expected Text 'Reis parallel kochen lassen', got %s", firstInstruction.Text)
	}

	// Assertions for Settings
	if !recipe.Content.Settings.Public {
		t.Errorf("Expected Public true, got %v", recipe.Content.Settings.Public)
	}
	if recipe.Content.Settings.ShowNutrition {
		t.Errorf("Expected ShowNutrition false, got %v", recipe.Content.Settings.ShowNutrition)
	}
	if recipe.Content.Settings.ShowAssets {
		t.Errorf("Expected ShowAssets false, got %v", recipe.Content.Settings.ShowAssets)
	}
	if recipe.Content.Settings.LandscapeView {
		t.Errorf("Expected LandscapeView false, got %v", recipe.Content.Settings.LandscapeView)
	}
	if recipe.Content.Settings.DisableComments {
		t.Errorf("Expected DisableComments false, got %v", recipe.Content.Settings.DisableComments)
	}
	if recipe.Content.Settings.Locked {
		t.Errorf("Expected Locked false, got %v", recipe.Content.Settings.Locked)
	}

	// Assertions for Nutrition (all nil)
	if recipe.Content.Nutrition.Calories != nil {
		t.Errorf("Expected Calories nil, got %v", recipe.Content.Nutrition.Calories)
	}
	if recipe.Content.Nutrition.CarbohydrateContent != nil {
		t.Errorf("Expected CarbohydrateContent nil, got %v", recipe.Content.Nutrition.CarbohydrateContent)
	}
	if recipe.Content.Nutrition.CholesterolContent != nil {
		t.Errorf("Expected CholesterolContent nil, got %v", recipe.Content.Nutrition.CholesterolContent)
	}
	if recipe.Content.Nutrition.FatContent != nil {
		t.Errorf("Expected FatContent nil, got %v", recipe.Content.Nutrition.FatContent)
	}
	if recipe.Content.Nutrition.FiberContent != nil {
		t.Errorf("Expected FiberContent nil, got %v", recipe.Content.Nutrition.FiberContent)
	}
	if recipe.Content.Nutrition.ProteinContent != nil {
		t.Errorf("Expected ProteinContent nil, got %v", recipe.Content.Nutrition.ProteinContent)
	}
	if recipe.Content.Nutrition.SaturatedFatContent != nil {
		t.Errorf("Expected SaturatedFatContent nil, got %v", recipe.Content.Nutrition.SaturatedFatContent)
	}
	if recipe.Content.Nutrition.SodiumContent != nil {
		t.Errorf("Expected SodiumContent nil, got %v", recipe.Content.Nutrition.SodiumContent)
	}
	if recipe.Content.Nutrition.SugarContent != nil {
		t.Errorf("Expected SugarContent nil, got %v", recipe.Content.Nutrition.SugarContent)
	}
	if recipe.Content.Nutrition.TransFatContent != nil {
		t.Errorf("Expected TransFatContent nil, got %v", recipe.Content.Nutrition.TransFatContent)
	}
	if recipe.Content.Nutrition.UnsaturatedFatContent != nil {
		t.Errorf("Expected UnsaturatedFatContent nil, got %v", recipe.Content.Nutrition.UnsaturatedFatContent)
	}

	// Assertions for RecipeScale
	if recipe.RecipeScale != 1 {
		t.Errorf("Expected RecipeScale 1, got %f", recipe.RecipeScale)
	}
}
