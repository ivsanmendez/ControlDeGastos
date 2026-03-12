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

	PermCategoryCreate Permission = "category:create"
	PermCategoryRead   Permission = "category:read"
	PermCategoryUpdate Permission = "category:update"
	PermCategoryDelete Permission = "category:delete"

	PermReceiptVerify Permission = "receipt:verify"

	PermReportRead Permission = "report:read"

	PermExpenseCategoryCreate Permission = "expense_category:create"
	PermExpenseCategoryRead   Permission = "expense_category:read"
	PermExpenseCategoryUpdate Permission = "expense_category:update"
	PermExpenseCategoryDelete Permission = "expense_category:delete"
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
		PermCategoryCreate,
		PermCategoryRead,
		PermCategoryUpdate,
		PermCategoryDelete,
		PermReceiptVerify,
		PermReportRead,
		PermExpenseCategoryCreate,
		PermExpenseCategoryRead,
		PermExpenseCategoryUpdate,
		PermExpenseCategoryDelete,
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
		PermCategoryCreate,
		PermCategoryRead,
		PermCategoryUpdate,
		PermCategoryDelete,
		PermReceiptVerify,
		PermReportRead,
		PermExpenseCategoryCreate,
		PermExpenseCategoryRead,
		PermExpenseCategoryUpdate,
		PermExpenseCategoryDelete,
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
