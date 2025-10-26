package mealie

import "time"

// Action represents the action part of the JSON.
type Action struct {
	ActionType  string `json:"action_type"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	GroupID     string `json:"group_id"`
	HouseholdID string `json:"household_id"`
	ID          string `json:"id"`
}

// Unit represents the unit of measurement for ingredients.
type Unit struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	PluralName         string            `json:"plural_name"`
	Description        string            `json:"description"`
	Extras             map[string]string `json:"extras"`
	Fraction           bool              `json:"fraction"`
	Abbreviation       string            `json:"abbreviation"`
	PluralAbbreviation *string           `json:"plural_abbreviation"`
	UseAbbreviation    bool              `json:"use_abbreviation"`
	Aliases            []string          `json:"aliases"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// Food represents the food item in the recipe.
type Food struct {
	ID                           string            `json:"id"`
	Name                         string            `json:"name"`
	PluralName                   *string           `json:"plural_name"`
	Description                  string            `json:"description"`
	Extras                       map[string]string `json:"extras"`
	LabelID                      *string           `json:"label_id"`
	Aliases                      []string          `json:"aliases"`
	HouseholdsWithIngredientFood []interface{}     `json:"households_with_ingredient_food"`
	Label                        *interface{}      `json:"label"`
	CreatedAt                    time.Time         `json:"created_at"`
	UpdatedAt                    time.Time         `json:"updated_at"`
}

// RecipeIngredient represents an ingredient in the recipe.
type RecipeIngredient struct {
	Quantity     *float64 `json:"quantity"`
	Unit         *Unit    `json:"unit"`
	Food         Food     `json:"food"`
	Note         string   `json:"note"`
	Display      string   `json:"display"`
	Title        *string  `json:"title"`
	OriginalText string   `json:"original_text"`
	ReferenceID  string   `json:"reference_id"`
}

// RecipeInstruction represents a step in the recipe instructions.
type RecipeInstruction struct {
	ID                   string        `json:"id"`
	Title                string        `json:"title"`
	Summary              string        `json:"summary"`
	Text                 string        `json:"text"`
	IngredientReferences []interface{} `json:"ingredient_references"`
}

// Nutrition represents the nutritional information of the recipe.
type Nutrition struct {
	Calories              *string `json:"calories"`
	CarbohydrateContent   *string `json:"carbohydrate_content"`
	CholesterolContent    *string `json:"cholesterol_content"`
	FatContent            *string `json:"fat_content"`
	FiberContent          *string `json:"fiber_content"`
	ProteinContent        *string `json:"protein_content"`
	SaturatedFatContent   *string `json:"saturated_fat_content"`
	SodiumContent         *string `json:"sodium_content"`
	SugarContent          *string `json:"sugar_content"`
	TransFatContent       *string `json:"trans_fat_content"`
	UnsaturatedFatContent *string `json:"unsaturated_fat_content"`
}

// Tag represents a tag
type Tag struct {
	Id      string `json:"id"`
	GroupId string `json:"group_id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
}

// Category represents a category for a recipe
type Category struct {
	Id      string `json:"id"`
	GroupId string `json:"group_id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
}

// Settings represents the settings for the recipe.
type Settings struct {
	Public          bool `json:"public"`
	ShowNutrition   bool `json:"show_nutrition"`
	ShowAssets      bool `json:"show_assets"`
	LandscapeView   bool `json:"landscape_view"`
	DisableComments bool `json:"disable_comments"`
	Locked          bool `json:"locked"`
}

// RecipeContent represents the main content of the JSON.
type RecipeContent struct {
	ID                  string              `json:"id"`
	UserID              string              `json:"user_id"`
	HouseholdID         string              `json:"household_id"`
	GroupID             string              `json:"group_id"`
	Name                string              `json:"name"`
	Slug                string              `json:"slug"`
	Image               string              `json:"image"`
	RecipeServings      float64             `json:"recipe_servings"`
	RecipeYieldQuantity float64             `json:"recipe_yield_quantity"`
	RecipeYield         string              `json:"recipe_yield"`
	TotalTime           *interface{}        `json:"total_time"`
	PrepTime            *interface{}        `json:"prep_time"`
	CookTime            *interface{}        `json:"cook_time"`
	PerformTime         *interface{}        `json:"perform_time"`
	Description         string              `json:"description"`
	RecipeCategory      []Category          `json:"recipe_category"`
	Tags                []Tag               `json:"tags"`
	Tools               []interface{}       `json:"tools"`
	Rating              *interface{}        `json:"rating"`
	OrgURL              string              `json:"org_url"`
	DateAdded           string              `json:"date_added"`
	DateUpdated         time.Time           `json:"date_updated"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
	LastMade            *interface{}        `json:"last_made"`
	RecipeIngredient    []RecipeIngredient  `json:"recipe_ingredient"`
	RecipeInstructions  []RecipeInstruction `json:"recipe_instructions"`
	Nutrition           Nutrition           `json:"nutrition"`
	Settings            Settings            `json:"settings"`
	Assets              []interface{}       `json:"assets"`
	Notes               []interface{}       `json:"notes"`
	Extras              map[string]string   `json:"extras"`
	Comments            []interface{}       `json:"comments"`
}

// Recipe represents the entire recipe JSON structure.
type Recipe struct {
	Action      Action        `json:"action"`
	Content     RecipeContent `json:"content"`
	RecipeScale float64       `json:"recipe_scale"`
}
