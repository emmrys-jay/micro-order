package domain

type ProductStatus int32

type ProductUpdateFromQueue struct {
	Id            string        `json:"id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Description   string        `json:"description,omitempty"`
	Price         float64       `json:"price,omitempty"`
	Quantity      int32         `json:"quantity,omitempty"`
	Status        ProductStatus `json:"status,omitempty"`
	OwnerId       string        `json:"owner_id,omitempty"`
	OwnerName     string        `json:"owner_name,omitempty"`
	OwnerPhone    string        `json:"owner_phone,omitempty"`
	OwnerEmail    string        `json:"owner_email,omitempty"`
	CreatedAt     string        `json:"created_at,omitempty"`
	UpdatedAt     string        `json:"updated_at,omitempty"`
	NameIsUpdated bool          `json:"name_is_updated"`
}
