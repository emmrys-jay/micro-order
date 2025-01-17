package domain

type UserRole int32

type UserUpdateForQueue struct {
	Id        string   `json:"id,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Email     string   `json:"email,omitempty"`
	Password  string   `json:"password,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Role      UserRole `json:"role,omitempty"`
	IsActive  bool     `json:"is_active,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
	DeletedAt string   `json:"deleted_at,omitempty"`
}
