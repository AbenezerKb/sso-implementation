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
	GetUser = Permission{
		ID:       "get_user",
		Name:     "get a user",
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
	DeleteClient = Permission{
		ID:       "delete_client",
		Name:     "delete a client",
		Category: "client",
	}
)
