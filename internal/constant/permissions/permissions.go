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
	GetAllClients = Permission{
		ID:       "get_all_clients",
		Name:     "get all clients",
		Category: "client",
	}
	GetAllUsers = Permission{
		ID:       "get_all_users",
		Name:     "get all users",
		Category: "user",
	}
	UpdateUser = Permission{
		ID:       "update_user",
		Name:     "update a user",
		Category: "user",
	}
	CreateResourceServer = Permission{
		ID:       "create_resource_server",
		Name:     "create a resource server",
		Category: "resource_server",
	}
	GetAllResourceServers = Permission{
		ID:       "get_all_resource_servers",
		Name:     "get all resource servers",
		Category: "resource_server",
	}
)
