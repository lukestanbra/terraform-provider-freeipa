package freeipa

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ipa "github.com/tehwalris/go-freeipa/freeipa"
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
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("%s", r)

	// Givenname should be required, however the admin user doesn't have one >:(
	if r.Result.Givenname == nil {
		d.Set("firstname", "")
	} else if err := d.Set("firstname", *r.Result.Givenname); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("lastname", r.Result.Sn); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(user_id)

	return diags
}
