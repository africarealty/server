package impl

import (
	"context"
	"github.com/africarealty/server/src/kit/auth"
	"github.com/africarealty/server/src/kit/log"
	"net/http"
)

type resourcePolicyManager struct {
	routePoliciesMap map[string][]auth.ResourcePolicy
	logger           log.CLoggerFunc
}

func (s *resourcePolicyManager) l() log.CLogger {
	return s.logger().Cmp("resource-policy-manager")
}

func NewResourcePolicyManager(logger log.CLoggerFunc) auth.ResourcePolicyManager {
	return &resourcePolicyManager{
		routePoliciesMap: map[string][]auth.ResourcePolicy{},
		logger:           logger,
	}
}

func (s *resourcePolicyManager) RegisterResourceMapping(routeId string, policies ...auth.ResourcePolicy) {
	s.l().Mth("register").Trc(routeId, " registered")
	s.routePoliciesMap[routeId] = policies
}

func (s *resourcePolicyManager) GetRequestedResources(ctx context.Context, routeId string, r *http.Request) ([]*auth.AuthorizationResource, error) {
	l := s.l().Mth("get-resources")

	var resources []*auth.AuthorizationResource
	var codes []string

	if policies, ok := s.routePoliciesMap[routeId]; ok {
		for _, policy := range policies {
			resource, err := policy.Resolve(ctx, r)
			if err != nil {
				return nil, err
			}
			if resource == nil {
				continue
			}
			resources = append(resources, resource)
			codes = append(codes, resource.Resource)
		}
	}
	l.F(log.FF{"routeId": routeId, "resources": codes}).Trc()

	return resources, nil
}
