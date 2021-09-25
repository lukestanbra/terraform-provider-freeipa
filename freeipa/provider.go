package freeipa

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ipa "github.com/tehwalris/go-freeipa/freeipa"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func Provider(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"host": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"username": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"password": &schema.Schema{
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"freeipa_user": dataSourceUser(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"freeipa_user": resourceUser(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

// type apiClient struct {
// 	// Add whatever fields, client or connection info, etc. here
// 	// you would need to setup to communicate with the upstream
// 	// API.
// }

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
		// TODO: myClient.UserAgent = userAgent
		host := d.Get("host").(string)
		user := d.Get("username").(string)
		pw := d.Get("password").(string)
		tspt := &http.Transport{}
		c, err := ipa.Connect(host, tspt, user, pw)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create FreeIPA client",
				Detail:   "Unable to create FreeIPA client",
			})
			return nil, diags
		}

		return c, diags

	}
}
