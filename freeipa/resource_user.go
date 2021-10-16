package freeipa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ipa "github.com/lukestanbra/go-freeipa/freeipa"
	"github.com/mitchellh/mapstructure"
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
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"full_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"initials": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"home_directory": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gecos": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"shell": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"krb_principal_name": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"email": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"uid": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"gid": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"street": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"city": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"postcode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"telephone_number": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"mobile_number": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"pager": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"fax_number": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"organisational_unit": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_title": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"manager": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"car_license": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*ipa.Client)

	var userAddArgs ipa.UserAddArgs
	userAddArgsMap := make(map[string]interface{})
	userAddArgsMap["givenname"] = d.Get("first_name")
	userAddArgsMap["sn"] = d.Get("last_name")

	mapstructure.Decode(userAddArgsMap, &userAddArgs)

	var userAddOptionalArgs ipa.UserAddOptionalArgs
	userAddOptionalArgsMap := make(map[string]interface{})
	userAddOptionalArgsMap["uid"] = d.Get("username")
	userAddOptionalArgsMap["all"] = true

	// Maps terraform resource variable names to FreeIPA Client names
	userAddOptionalVarsMap := map[string]string{
		"full_name":           "cn",
		"display_name":        "displayname",
		"initials":            "initials",
		"home_directory":      "homedirectory",
		"gecos":               "gecos",
		"shell":               "loginshell",
		"krb_principal_name":  "krbprincipalname",
		"email":               "mail",
		"uid":                 "uidnumber",
		"gid":                 "gidnumber",
		"street":              "street",
		"city":                "l",
		"state":               "st",
		"postcode":            "postalcode",
		"telephone_number":    "telephonenumber",
		"mobile_number":       "mobile",
		"pager":               "pager",
		"fax_number":          "facsimiletelephonenumber",
		"organisational_unit": "ou",
		"job_title":           "title",
		"manager":             "manager",
		"car_license":         "carlicense",
	}

	for k, v := range userAddOptionalVarsMap {
		if val, ok := d.GetOk(k); ok {
			userAddOptionalArgsMap[v] = val
		}
	}

	mapstructure.Decode(userAddOptionalArgsMap, &userAddOptionalArgs)

	_, err := c.UserAdd(
		&userAddArgs,
		&userAddOptionalArgs,
	)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create FreeIPA user",
			Detail:   err.Error(),
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
			All: ipa.Bool(true),
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	userShowResultsVarsMap := map[string]interface{}{
		"first_name":          r.Result.Givenname,
		"last_name":           r.Result.Sn,
		"full_name":           r.Result.Cn,
		"display_name":        r.Result.Displayname,
		"initials":            r.Result.Initials,
		"home_directory":      r.Result.Homedirectory,
		"gecos":               r.Result.Gecos,
		"shell":               r.Result.Loginshell,
		"krb_principal_name":  r.Result.Krbprincipalname,
		"email":               r.Result.Mail,
		"uid":                 r.Result.Uidnumber,
		"gid":                 r.Result.Gidnumber,
		"street":              r.Result.Street,
		"city":                r.Result.L,
		"state":               r.Result.St,
		"postcode":            r.Result.Postalcode,
		"telephone_number":    r.Result.Telephonenumber,
		"mobile_number":       r.Result.Mobile,
		"pager":               r.Result.Pager,
		"fax_number":          r.Result.Facsimiletelephonenumber,
		"organisational_unit": r.Result.Ou,
		"job_title":           r.Result.Title,
		"manager":             r.Result.Manager,
		"car_license":         r.Result.Carlicense,
	}

	for k, v := range userShowResultsVarsMap {
		d.Set(k, v)
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ipa.Client)

	// Maps terraform resource variable names to FreeIPA Client names
	userModOptionalVarsMap := map[string]string{
		"first_name":          "givenname",
		"last_name":           "sn",
		"full_name":           "cn",
		"display_name":        "displayname",
		"initials":            "initials",
		"home_directory":      "homedirectory",
		"gecos":               "gecos",
		"shell":               "loginshell",
		"krb_principal_name":  "krbprincipalname",
		"email":               "mail",
		"uid":                 "uidnumber",
		"gid":                 "gidnumber",
		"street":              "street",
		"city":                "l",
		"state":               "st",
		"postcode":            "postalcode",
		"telephone_number":    "telephonenumber",
		"mobile_number":       "mobile",
		"pager":               "pager",
		"fax_number":          "facsimiletelephonenumber",
		"organisational_unit": "ou",
		"job_title":           "title",
		"manager":             "manager",
		"car_license":         "carlicense",
	}

	// Extract the keys from the map above
	keys := make([]string, len(userModOptionalVarsMap))
	i := 0
	for k := range userModOptionalVarsMap {
		keys[i] = k
		i++
	}

	if d.HasChanges(keys...) {
		var userModOptionalArgs ipa.UserModOptionalArgs
		userModOptionalArgsMap := make(map[string]interface{})
		userModOptionalArgsMap["uid"] = d.Id()
		userModOptionalArgsMap["all"] = true

		for k, v := range userModOptionalVarsMap {
			if val, ok := d.GetOk(k); ok {
				userModOptionalArgsMap[v] = val
			}
		}

		mapstructure.Decode(userModOptionalArgsMap, &userModOptionalArgs)

		_, err := c.UserMod(
			&ipa.UserModArgs{},
			&userModOptionalArgs,
		)

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
