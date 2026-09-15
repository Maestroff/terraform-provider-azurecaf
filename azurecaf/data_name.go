package azurecaf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// dataName creates and returns the schema for the azurecaf_name data source.
//
// This data source provides the same naming functionality as the azurecaf_name resource
// but is evaluated during the plan phase, making the generated name visible before
// resource creation. This is the recommended approach for most use cases.
//
// Key benefits of using the data source over the resource:
//   - Names are generated during terraform plan, providing early visibility
//   - No state management required (data sources are read-only)
//   - Better for single resource name generation
//   - Integrates naturally with Terraform's data flow
//
// Use the resource version when you need to generate multiple related names
// using the resource_types parameter.
func dataName() *schema.Resource {
	resourceMapsKeys := make([]string, 0, len(ResourceDefinitions))
	for k := range ResourceDefinitions {
		resourceMapsKeys = append(resourceMapsKeys, k)
	}

	return &schema.Resource{
		ReadContext: dataNameRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"prefixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			"suffixes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
				ForceNew: true,
			},
			"component_order": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MinItems:    5,
				MaxItems:    5,
				Description: "Left-to-right order for prefixes, slug, name, random, and suffixes. Omit to use the existing provider order.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(defaultComponentOrder, false),
				},
			},
			"random_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      0,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"separator": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "-",
			},
			"clean_input": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"passthrough": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(resourceMapsKeys, false),
				ForceNew:     true,
			},
			"random_seed": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"use_slug": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"use_legacy_slug": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Use legacy slug for backward compatibility (default: false in v4.0.0+, set true to maintain v3.x behavior)",
			},
			"error_when_exceeding_max_length": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Return an error instead of omitting name components when the composed name exceeds the resource type maximum length.",
			},
		},
	}
}

func dataNameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := getNameReadResult(d, meta); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func getNameReadResult(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	prefixes := convertInterfaceToString(d.Get("prefixes").([]interface{}))
	suffixes := convertInterfaceToString(d.Get("suffixes").([]interface{}))
	separator := d.Get("separator").(string)
	resourceType := d.Get("resource_type").(string)
	cleanInput := d.Get("clean_input").(bool)
	passthrough := d.Get("passthrough").(bool)
	useSlug := d.Get("use_slug").(bool)
	useLegacySlug := d.Get("use_legacy_slug").(bool)
	errorWhenExceedingMaxLength := d.Get("error_when_exceeding_max_length").(bool)
	componentOrder, err := componentOrderFromResourceData(d)
	if err != nil {
		return err
	}
	randomLength := d.Get("random_length").(int)
	randomSeed := int64(d.Get("random_seed").(int))

	convention := ConventionCafClassic

	randomSuffix := randSeq(int(randomLength), &randomSeed)

	namePrecedence := nameComponentPrecedence

	resourceName, err := getResourceNameWithComponentOrder(resourceType, separator, prefixes, name, suffixes, randomSuffix, convention, cleanInput, passthrough, useSlug, useLegacySlug, namePrecedence, componentOrder, errorWhenExceedingMaxLength)
	if err != nil {
		return err
	}
	d.Set("result", resourceName)

	d.SetId(resourceName)
	return nil
}
