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
	CreateResourceServer = Permission{
		ID:       "create_resource_server",
		Name:     "create a resource server",
		Category: "resource_server",
	}
	GetAllScopes = Permission{
		ID:       "get_all_scopes",
		Name:     "get all scope",
		Category: "scope",
	}
	GetAllResourceServers = Permission{
		ID:       "get_all_resource_servers",
		Name:     "get all resource servers",
		Category: "resource_server",
	}
	GetClient = Permission{
		ID:       "get_client",
		Name:     "get a client",
		Category: "client",
	}
	UpdateClient = Permission{
		ID:       "update_client",
		Name:     "update a client",
		Category: "client",
	}
	GetAllPermissions = Permission{
		ID:       "get_all_permissions",
		Name:     "Get all permissions",
		Category: "role",
	}
	DeleteScope = Permission{
		ID:       "delete_scope",
		Name:     "Delete a scope",
		Category: "scope",
	}
	CreateRole = Permission{
		ID:       "create_role",
		Name:     "create a role",
		Category: "role",
	}
	UpdateScope = Permission{
		ID:       "update_scope",
		Name:     "Update a scope",
		Category: "scope",
	}
	GetAllRoles = Permission{
		ID:       "get_all_roles",
		Name:     "get all roles",
		Category: "role",
	}
	UpdateUserRole = Permission{
		ID:       "update_user_role",
		Name:     "update user role",
		Category: "user",
	}
	UpdateUserStatus = Permission{
		ID:       "update_user_status",
		Name:     "update user status",
		Category: "user",
	}
	ChangeRoleStatus = Permission{
		ID:       "change_role_status",
		Name:     "change role status",
		Category: "role",
	}
	GetRole = Permission{
		ID:       "get_role",
		Name:     "get a role",
		Category: "role",
	}
	DeleteRole = Permission{
		ID:       "delete_role",
		Name:     "delete a role",
		Category: "role",
	}
	UpdateRole = Permission{
		ID:       "update_role",
		Name:     "update a role",
		Category: "role",
	}
	CreateIdentityProvider = Permission{
		ID:       "create_identity_provider",
		Name:     "create an identity provider",
		Category: "identity_provider",
	}
	UpdateClientStatus = Permission{
		ID:       "update_client_status",
		Name:     "Update client status",
		Category: "client",
	}
	UpdateIdentityProvider = Permission{
		ID:       "update_identity_provider",
		Name:     "update an identity provider",
		Category: "identity_provider",
	}
	GetIdentityProvider = Permission{
		ID:       "get_identity_provider",
		Name:     "get an identity provider",
		Category: "identity_provider",
	}
	DeleteIdentityProvider = Permission{
		ID:       "delete_identity_provider",
		Name:     "delete an identity provider",
		Category: "identity_provider",
	}
	GetAllIdentityProviders = Permission{
		ID:       "get_all_identity_providers",
		Name:     "get all identity providers",
		Category: "identity_provider",
	}
	RevokeUserRole = Permission{
		ID:       "revoke_user_role",
		Name:     "revoke user's role",
		Category: "user",
	}
	ResetUserPassword = Permission{
		ID:       "reset_user_password",
		Name:     "reset user password",
		Category: "user",
	}
)
