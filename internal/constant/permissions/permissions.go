package permissions

const Notneeded = "Not Needed"

type Permission struct {
	ID       string
	Name     string
	Category string
}

// permissions list

var (
	CreateUser = Permission{
		ID:       "create_user",
		Name:     "create a user",
		Category: "user",
	}
	CreateClient = Permission{
		ID:       "create_client",
		Name:     "create a client",
		Category: "client",
	}
	CreateScope = Permission{
		ID:       "create_scope",
		Name:     "create a scope",
		Category: "scope",
	}
	GetScope = Permission{
		ID:       "get_scope",
		Name:     "get a scope",
		Category: "scope",
	}
)
