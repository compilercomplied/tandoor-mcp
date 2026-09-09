package food

// FoodResponse represents a food item in Tandoor
type FoodResponse struct {
	ID                  int     `json:"id"`
	Name                string  `json:"name"`
	PluralName          *string `json:"plural_name"`
	Description         string  `json:"description"`
	FullName            string  `json:"full_name"`
	IgnoreShopping      bool    `json:"ignore_shopping"`
	SupermarketCategory *struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"supermarket_category"`
	Parent   *int `json:"parent"`
	NumChild int  `json:"numchild"`
}

// FoodListResponse represents a paginated list of foods
type FoodListResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []FoodResponse `json:"results"`
}

// FoodParam represents request parameters to create a food item
type FoodParam struct {
	Name           string  `json:"name"`
	PluralName     *string `json:"plural_name,omitempty"`
	Description    string  `json:"description,omitempty"`
	IgnoreShopping bool    `json:"ignore_shopping,omitempty"`
	Parent         *int    `json:"parent,omitempty"`
}

type UpdateParams struct {
	Name           *string `json:"name,omitempty"`
	PluralName     *string `json:"plural_name,omitempty"`
	Description    *string `json:"description,omitempty"`
	IgnoreShopping *bool   `json:"ignore_shopping,omitempty"`
	Parent         *int    `json:"parent,omitempty"`
}

// FoodInheritFieldResponse represents a field inheritance configuration
type FoodInheritFieldResponse struct {
	ID    int     `json:"id"`
	Name  *string `json:"name"`
	Field *string `json:"field"`
}

// FoodInheritFieldListResponse represents a list of inheritance fields
type FoodInheritFieldListResponse struct {
	Count   int                        `json:"count"`
	Results []FoodInheritFieldResponse `json:"results"`
}

// ListParams represents query parameters for listing foods
type ListParams struct {
	Query    *string `json:"query,omitempty"`
	Root     *int    `json:"root,omitempty"`
	Tree     *int    `json:"tree,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"page_size,omitempty"`
}
