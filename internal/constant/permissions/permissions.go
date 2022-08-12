package permissions

const Notneeded = "Not Needed"

// permissions list
const (
	CreateSystemUser = "create a user"
)

var PermissionCategory = map[string]string{
	CreateSystemUser: "user",
}
