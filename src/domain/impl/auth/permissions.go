package auth

import (
	"github.com/africarealty/server/src/domain"
	"github.com/africarealty/server/src/kit/auth"
)

// groupRoles specifies links between groups and roles
var groupRoles = map[string][]string{
	domain.AuthGroupSysAdmin: {domain.AuthRoleSysAdmin},
	domain.AuthGroupOwner:    {domain.AuthRoleProfileOwner},
	domain.AuthGroupAgent:    {domain.AuthRoleProfileOwner},
}

// permissions specifies access on resources for session roles
var permissions = map[string][]rolePermissions{
	domain.AuthResUserProfileAll: {rolePermissions{Role: domain.AuthRoleSysAdmin, Permissions: []string{auth.AccessR, auth.AccessW, auth.AccessD}}},
	domain.AuthResUserProfileMy:  {rolePermissions{Role: domain.AuthRoleProfileOwner, Permissions: []string{auth.AccessR, auth.AccessW}}},
}
