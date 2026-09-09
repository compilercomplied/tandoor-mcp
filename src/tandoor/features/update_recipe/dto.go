package update_recipe

// Params contains the recipe fields that may be updated. Nil fields are left unchanged.
type Params struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Servings    *int     `json:"servings,omitempty"`
	WorkingTime *int     `json:"working_time,omitempty"`
	WaitingTime *int     `json:"waiting_time,omitempty"`
	SourceURL   *string  `json:"source_url,omitempty"`
	Private     *bool    `json:"private,omitempty"`
	Keywords    *[]IDRef `json:"keywords,omitempty"`
	Properties  *[]IDRef `json:"properties,omitempty"`
}

type IDRef struct {
	ID int `json:"id"`
}

type Response struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Servings    int     `json:"servings"`
	WorkingTime int     `json:"working_time"`
	WaitingTime int     `json:"waiting_time"`
	SourceURL   *string `json:"source_url"`
	Private     bool    `json:"private"`
}
