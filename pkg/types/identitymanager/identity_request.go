package identitymanager

type CreateIdentityRequest struct {
	Alias      string            `json:"alias" validate:"required" example:"personal-account"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
