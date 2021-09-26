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
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"full_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"initials": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"home_directory": &schema.Schema{
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
			"krb_canonical_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"krb_principal_name": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"email": &schema.Schema{
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
			"telephone_number": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"mobile_number": &schema.Schema{
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
			"fax_number": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"organisational_unit": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"job_title": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"manager": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"car_license": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ssh_public_key": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ssh_public_key_fingerprint": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_auth_type": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_class": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"token_radius_config_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"token_radius_username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"department_number": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"employee_number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"employee_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"preferred_language": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"has_password": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
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
		"first_name":                 r.Result.Givenname,
		"last_name":                  r.Result.Sn,
		"full_name":                  r.Result.Cn,
		"display_name":               r.Result.Displayname,
		"initials":                   r.Result.Initials,
		"home_directory":             r.Result.Homedirectory,
		"gecos":                      r.Result.Gecos,
		"shell":                      r.Result.Loginshell,
		"krb_canonical_name":         r.Result.Krbcanonicalname,
		"krb_principal_name":         r.Result.Krbprincipalname,
		"email":                      r.Result.Mail,
		"uid":                        r.Result.Uidnumber,
		"gid":                        r.Result.Gidnumber,
		"street":                     r.Result.Street,
		"city":                       r.Result.L,
		"state":                      r.Result.St,
		"postcode":                   r.Result.Postalcode,
		"telephone_number":           r.Result.Telephonenumber,
		"mobile_number":              r.Result.Mobile,
		"pager":                      r.Result.Pager,
		"fax_number":                 r.Result.Facsimiletelephonenumber,
		"organisational_unit":        r.Result.Ou,
		"job_title":                  r.Result.Title,
		"manager":                    r.Result.Manager,
		"car_license":                r.Result.Carlicense,
		"ssh_public_key":             r.Result.Ipasshpubkey,
		"ssh_public_key_fingerprint": r.Result.Sshpubkeyfp,
		"user_auth_type":             r.Result.Ipauserauthtype,
		"user_class":                 r.Result.Userclass,
		"token_radius_config_link":   r.Result.Ipatokenradiusconfiglink,
		"token_radius_username":      r.Result.Ipatokenradiususername,
		"department_number":          r.Result.Departmentnumber,
		"employee_number":            r.Result.Employeenumber,
		"employee_type":              r.Result.Employeetype,
		"preferredl_anguage":         r.Result.Preferredlanguage,
		"has_password":               r.Result.HasPassword,
	}

	for k, v := range userResultsVarsMap {
		d.Set(k, v)
	}

	// always run
	d.SetId(user_id)

	return diags
}
