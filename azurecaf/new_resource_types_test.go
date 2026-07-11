package azurecaf

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Covers the resource types added for azurerm v4.80 parity.
func TestNewResourceTypes_Azurerm480(t *testing.T) {
	provider := Provider()
	nameResource := provider.ResourcesMap["azurecaf_name"]
	if nameResource == nil {
		t.Fatal("azurecaf_name resource not found")
	}

	cases := []struct {
		resourceType string
		slug         string
	}{
		{"azurerm_windows_function_app", "winfa"},
		{"azurerm_windows_function_app_slot", "winfas"},
		{"azurerm_windows_web_app_slot", "wwapps"},
		{"azurerm_managed_redis", "amr"},
		{"azurerm_managed_devops_pool", "mdp"},
		{"azurerm_trusted_signing_account", "tsa"},
		{"azurerm_eventgrid_partner_namespace", "egpn"},
		{"azurerm_eventgrid_partner_registration", "egpr"},
		{"azurerm_eventgrid_namespace_topic", "egnt"},
		{"azurerm_mongo_cluster_firewall_rule", "mongofw"},
		{"azurerm_mongo_cluster_user", "mongou"},
		{"azurerm_network_security_perimeter_profile", "nspp"},
		{"azurerm_network_security_perimeter_access_rule", "nspar"},
		{"azurerm_cognitive_account_project", "proj"},
	}

	for _, tc := range cases {
		t.Run(tc.resourceType, func(t *testing.T) {
			def, ok := ResourceDefinitions[tc.resourceType]
			if !ok {
				t.Fatalf("resource type %s not found in ResourceDefinitions", tc.resourceType)
			}
			if def.CafPrefix != tc.slug {
				t.Errorf("expected slug %q for %s, got %q", tc.slug, tc.resourceType, def.CafPrefix)
			}

			resourceData := schema.TestResourceDataRaw(t, nameResource.Schema, map[string]interface{}{
				"name":          "demo",
				"resource_type": tc.resourceType,
				"prefixes":      []interface{}{"dev"},
				"random_seed":   1,
				"random_length": 5,
				"clean_input":   true,
			})

			if err := nameResource.Create(resourceData, nil); err != nil {
				t.Fatalf("failed to generate name for %s: %v", tc.resourceType, err)
			}

			result := resourceData.Get("result").(string)
			if !strings.Contains(result, tc.slug) {
				t.Errorf("expected result to contain slug %q, got %q", tc.slug, result)
			}
			if !regexp.MustCompile(def.ValidationRegExp).MatchString(result) {
				t.Errorf("result %q does not match validation regex %s", result, def.ValidationRegExp)
			}
			if len(result) < def.MinLength || len(result) > def.MaxLength {
				t.Errorf("result %q length %d outside [%d, %d]", result, len(result), def.MinLength, def.MaxLength)
			}
		})
	}
}
