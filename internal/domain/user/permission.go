package user

type Permission string

const (
	PermExpenseCreate    Permission = "expense:create"
	PermExpenseReadOwn   Permission = "expense:read:own"
	PermExpenseReadAll   Permission = "expense:read:all"
	PermExpenseDeleteOwn Permission = "expense:delete:own"
	PermExpenseDeleteAll Permission = "expense:delete:all"

	PermContributionCreate Permission = "contribution:create"
	PermContributionRead   Permission = "contribution:read"
	PermContributionDelete Permission = "contribution:delete"

	PermContributorCreate Permission = "contributor:create"
	PermContributorRead   Permission = "contributor:read"
	PermContributorUpdate Permission = "contributor:update"
	PermContributorDelete Permission = "contributor:delete"
)

var rolePermissions = map[Role][]Permission{
	RoleUser: {
		PermExpenseCreate,
		PermExpenseReadOwn,
		PermExpenseDeleteOwn,
		PermContributionCreate,
		PermContributionRead,
		PermContributionDelete,
		PermContributorCreate,
		PermContributorRead,
		PermContributorUpdate,
		PermContributorDelete,
	},
	RoleAdmin: {
		PermExpenseCreate,
		PermExpenseReadOwn,
		PermExpenseReadAll,
		PermExpenseDeleteOwn,
		PermExpenseDeleteAll,
		PermContributionCreate,
		PermContributionRead,
		PermContributionDelete,
		PermContributorCreate,
		PermContributorRead,
		PermContributorUpdate,
		PermContributorDelete,
	},
}

// RoleHasPermission reports whether the given role includes the permission.
func RoleHasPermission(role Role, perm Permission) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}
