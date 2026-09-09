package update_recipe

// Params contains the recipe fields that may be updated. Nil fields are left unchanged.
type Params struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Servings    *int    `json:"servings,omitempty"`
	WorkingTime *int    `json:"working_time,omitempty"`
	WaitingTime *int    `json:"waiting_time,omitempty"`
}

type Response struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Servings    int    `json:"servings"`
	WorkingTime int    `json:"working_time"`
	WaitingTime int    `json:"waiting_time"`
}
