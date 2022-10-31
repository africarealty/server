package auth

import (
	"context"
	"net/http"
	"strings"
)

type ResourcePolicyBuilder struct {
	resource    string
	permissions []string
	conditions  []ConditionFn
}

func Resource(resource string, permissions string) *ResourcePolicyBuilder {
	b := &ResourcePolicyBuilder{
		resource:   resource,
		conditions: []ConditionFn{},
	}
	b.permissions = b.convertPermissions(permissions)
	return b
}

// convertPermissions converts permissions from "rwxd" string to []string{"r", w", "x", "d"}
func (a *ResourcePolicyBuilder) convertPermissions(permissions string) []string {
	var res []string
	s := strings.ToLower(permissions)
	if strings.Contains(s, AccessR) {
		res = append(res, AccessR)
	}
	if strings.Contains(s, AccessW) {
		res = append(res, AccessW)
	}
	if strings.Contains(s, AccessX) {
		res = append(res, AccessX)
	}
	if strings.Contains(s, AccessD) {
		res = append(res, AccessD)
	}
	return res
}

func (a *ResourcePolicyBuilder) When(f ConditionFn) *ResourcePolicyBuilder {
	a.conditions = append(a.conditions, f)
	return a
}

func (a *ResourcePolicyBuilder) WhenNot(f ConditionFn) *ResourcePolicyBuilder {
	a.conditions = append(a.conditions, func(c context.Context, r *http.Request) (bool, error) { res, err := f(c, r); return !res, err })
	return a
}

func (a *ResourcePolicyBuilder) Resolve(ctx context.Context, r *http.Request) (*AuthorizationResource, error) {
	// check conditions
	for _, cond := range a.conditions {
		if condRes, err := cond(ctx, r); err != nil {
			return nil, err
		} else {
			if !condRes {
				return nil, nil
			}
		}
	}
	return &AuthorizationResource{
		Resource:    a.resource,
		Permissions: a.permissions,
	}, nil
}

func (a *ResourcePolicyBuilder) B() ResourcePolicy {
	return a
}
