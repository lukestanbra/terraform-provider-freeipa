package freeipa

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ipa "github.com/tehwalris/go-freeipa/freeipa"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"firstname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"lastname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"shell": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/bin/bash",
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*ipa.Client)

	_, err := c.UserAdd(
		&ipa.UserAddArgs{
			Givenname: *ipa.String(d.Get("firstname").(string)),
			Sn:        *ipa.String(d.Get("lastname").(string)),
		},
		&ipa.UserAddOptionalArgs{
			UID:        ipa.String(d.Get("username").(string)),
			Loginshell: ipa.String(d.Get("shell").(string)),
		},
	)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create FreeIPA user",
			Detail:   "Unable to create FreeIPA user",
		})
	}

	d.SetId(d.Get("username").(string))

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*ipa.Client)

	user_id := d.Id()

	r, err := c.UserShow(
		&ipa.UserShowArgs{},
		&ipa.UserShowOptionalArgs{
			UID: ipa.String(user_id),
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	// Givenname should be required, however the admin user doesn't have one >:(
	if r.Result.Givenname == nil {
		d.Set("firstname", "")
	} else if err := d.Set("firstname", *r.Result.Givenname); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("lastname", r.Result.Sn); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ipa.Client)
	if d.HasChanges(
		"firstname",
		"lastname",
		"shell",
	) {
		r, err := c.UserMod(
			&ipa.UserModArgs{},
			&ipa.UserModOptionalArgs{
				UID:        ipa.String(d.Id()),
				Givenname:  ipa.String(d.Get("firstname").(string)),
				Sn:         ipa.String(d.Get("lastname").(string)),
				Loginshell: ipa.String(d.Get("shell").(string)),
			},
		)
		log.Printf("%s", r)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*ipa.Client)

	_, err := c.UserDel(
		&ipa.UserDelArgs{},
		&ipa.UserDelOptionalArgs{
			UID: &[]string{d.Id()},
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
