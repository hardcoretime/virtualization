package rewriter

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestRewriter() *RuleBasedRewriter {
	apiGroupRules := map[string]APIGroupRule{
		"original.group.io": {
			GroupRule: GroupRule{
				Group:            "original.group.io",
				Versions:         []string{"v1", "v1alpha1"},
				PreferredVersion: "v1",
			},
			ResourceRules: map[string]ResourceRule{
				"someresources": {
					Kind:             "SomeResource",
					ListKind:         "SomeResourceList",
					Plural:           "someresources",
					Singular:         "someresource",
					Versions:         []string{"v1", "v1alpha1"},
					PreferredVersion: "v1",
					Categories:       []string{"all"},
					ShortNames:       []string{"sr", "srs"},
				},
				"anotherresources": {
					Kind:             "AnotherResource",
					ListKind:         "AnotherResourceList",
					Plural:           "anotherresources",
					Singular:         "anotherresource",
					Versions:         []string{"v1", "v1alpha1"},
					PreferredVersion: "v1",
					ShortNames:       []string{"ar"},
				},
			},
		},
		"other.group.io": {
			GroupRule: GroupRule{
				Group:            "original.group.io",
				Versions:         []string{"v2alpha3"},
				PreferredVersion: "v2alpha3",
			},
			ResourceRules: map[string]ResourceRule{
				"otherresources": {
					Kind:             "OtherResource",
					ListKind:         "OtherResourceList",
					Plural:           "otherresources",
					Singular:         "otherresource",
					Versions:         []string{"v1", "v1alpha1"},
					PreferredVersion: "v1",
					ShortNames:       []string{"or"},
				},
			},
		},
	}

	rules := &RewriteRules{
		KindPrefix:         "Prefixed", // KV
		ResourceTypePrefix: "prefixed", // kv
		ShortNamePrefix:    "p",
		Categories:         []string{"prefixed"},
		RenamedGroup:       "prefixed.resources.group.io",
		Rules:              apiGroupRules,
	}

	return &RuleBasedRewriter{
		Rules: rules,
	}
}

func TestRewriteAPIEndpoint(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		expect string
	}{
		{
			"rewritable group",
			"/apis/original.group.io",
			"/apis/prefixed.resources.group.io",
		},
		{
			"rewritable group and version",
			"/apis/original.group.io/v1",
			"/apis/prefixed.resources.group.io/v1",
		},
		{
			"rewritable resource list",
			"/apis/original.group.io/v1/someresources",
			"/apis/prefixed.resources.group.io/v1/prefixedsomeresources",
		},
		{
			"rewritable resource by name",
			"/apis/original.group.io/v1/someresources/srname",
			"/apis/prefixed.resources.group.io/v1/prefixedsomeresources/srname",
		},
		{
			"rewritable resource status",
			"/apis/original.group.io/v1/someresources/srname/status",
			"/apis/prefixed.resources.group.io/v1/prefixedsomeresources/srname/status",
		},
		{
			"rewritable CRD",
			"/apis/apiextensions.k8s.io/v1/customresourcedefinitions/someresources.original.group.io",
			"/apis/apiextensions.k8s.io/v1/customresourcedefinitions/prefixedsomeresources.prefixed.resources.group.io",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.path)
			require.NoError(t, err, "should parse path '%s'", tt.path)

			ep := ParseAPIEndpoint(u)

			rwr := createTestRewriter()

			newEp := rwr.RewriteAPIEndpoint(ep)

			if tt.expect == "" {
				require.Nil(t, newEp, "should not rewrite path '%s', got %+v", tt.path, newEp)
			}

			require.NotNil(t, newEp, "should rewrite path '%s', got nil originEndpoint")

			require.Equal(t, tt.expect, newEp.Path(), "expect rewrite for path '%s' to be '%s', got '%s'", tt.path, tt.expect, ep.Path())
		})
	}

}