package unit

// UnitResponse represents a unit in Tandoor
type UnitResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	PluralName   *string `json:"plural_name"`
	Description  *string `json:"description"`
	BaseUnit     *string `json:"base_unit"`
	OpenDataSlug *string `json:"open_data_slug"`
}

// UnitListResponse represents a paginated list of units
type UnitListResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []UnitResponse `json:"results"`
}

// UnitParam represents parameters to create a unit
type UnitParam struct {
	Name         string  `json:"name"`
	PluralName   *string `json:"plural_name,omitempty"`
	Description  *string `json:"description,omitempty"`
	BaseUnit     *string `json:"base_unit,omitempty"`
	OpenDataSlug *string `json:"open_data_slug,omitempty"`
}

type UpdateParams struct {
	Name         *string `json:"name,omitempty"`
	PluralName   *string `json:"plural_name,omitempty"`
	Description  *string `json:"description,omitempty"`
	BaseUnit     *string `json:"base_unit,omitempty"`
	OpenDataSlug *string `json:"open_data_slug,omitempty"`
}

// UnitRef is a reference to a unit for nested structures
type UnitRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// FoodRef is a reference to a food item for nested structures
type FoodRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UnitConversionResponse represents a unit conversion in Tandoor
type UnitConversionResponse struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	BaseAmount      float64      `json:"base_amount"`
	BaseUnit        UnitResponse `json:"base_unit"`
	ConvertedAmount float64      `json:"converted_amount"`
	ConvertedUnit   UnitResponse `json:"converted_unit"`
	Food            *FoodRef     `json:"food"`
	OpenDataSlug    *string      `json:"open_data_slug"`
}

// UnitConversionListResponse represents a paginated list of unit conversions
type UnitConversionListResponse struct {
	Count    int                      `json:"count"`
	Next     *string                  `json:"next"`
	Previous *string                  `json:"previous"`
	Results  []UnitConversionResponse `json:"results"`
}

// UnitConversionParam represents parameters to create a unit conversion
type UnitConversionParam struct {
	BaseAmount      float64  `json:"base_amount"`
	BaseUnit        UnitRef  `json:"base_unit"`
	ConvertedAmount float64  `json:"converted_amount"`
	ConvertedUnit   UnitRef  `json:"converted_unit"`
	Food            *FoodRef `json:"food,omitempty"`
}

// ListUnitsParams represents query parameters for listing units
type ListUnitsParams struct {
	Query    *string `json:"query,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"page_size,omitempty"`
}

// ListConversionsParams represents query parameters for listing unit conversions
type ListConversionsParams struct {
	FoodID   *int    `json:"food_id,omitempty"`
	Query    *string `json:"query,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"page_size,omitempty"`
}
