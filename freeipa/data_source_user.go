package freeipa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ipa "github.com/lukestanbra/go-freeipa/freeipa"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"firstname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"lastname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fullname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"displayname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"initials": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"homedirectory": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"gecos": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"shell": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"krbcanonicalname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"krbprincipalname": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"uid": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"gid": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"street": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"city": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"postcode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"telephonenumber": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"mobilenumber": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"pager": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ipa.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	user_id := d.Get("username").(string)

	r, err := c.UserShow(
		&ipa.UserShowArgs{},
		&ipa.UserShowOptionalArgs{
			UID: ipa.String(user_id),
			All: ipa.Bool(true),
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	userResultsVarsMap := map[string]interface{}{
		"firstname":        r.Result.Givenname,
		"lastname":         r.Result.Sn,
		"fullname":         r.Result.Cn,
		"displayname":      r.Result.Displayname,
		"initials":         r.Result.Initials,
		"homedirectory":    r.Result.Homedirectory,
		"gecos":            r.Result.Gecos,
		"shell":            r.Result.Loginshell,
		"krbcanonicalname": r.Result.Krbcanonicalname,
		"krbprincipalname": r.Result.Krbprincipalname,
		"uid":              r.Result.Uidnumber,
		"gid":              r.Result.Gidnumber,
		"street":           r.Result.Street,
		"city":             r.Result.L,
		"state":            r.Result.St,
		"postcode":         r.Result.Postalcode,
		"telephonenumber":  r.Result.Telephonenumber,
		"mobilenumber":     r.Result.Mobile,
		"pager":            r.Result.Pager,
	}

	for k, v := range userResultsVarsMap {
		d.Set(k, v)
	}

	// always run
	d.SetId(user_id)

	return diags
}
