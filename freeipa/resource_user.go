package freeipa

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ipa "github.com/lukestanbra/go-freeipa/freeipa"
	"github.com/mitchellh/mapstructure"
)

// UserShowResult{
// 	"result":{
// 		"uid":"jdoe4733520369981736867",
// 		"givenname":"John",
// 		"sn":"Doe",
// 		"homedirectory":"/home/jdoe4733520369981736867",
// 		"loginshell":"/bin/sh",
// 		"krbcanonicalname":"jdoe4733520369981736867@EXAMPLE.TEST",
// 		"krbprincipalname":["jdoe4733520369981736867@EXAMPLE.TEST"],
// 		"mail":["jdoe4733520369981736867@example.test"],
// 		"uidnumber":1237800001,
// 		"gidnumber":1237800001,
// 		"nsaccountlock":false,
// 		"has_password":false,
// 		"memberof_group":["ipausers"],
// 		"has_keytab":false
// 	},
// 	"value":"jdoe4733520369981736867"
// }

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
			"fullname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"homedirectory": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"shell": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/bin/sh",
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
	userAddArgsMap["givenname"] = d.Get("firstname")
	userAddArgsMap["sn"] = d.Get("lastname")

	mapstructure.Decode(userAddArgsMap, &userAddArgs)

	log.Printf("%s", userAddArgsMap)
	log.Printf("%s", userAddArgs)

	var userAddOptionalArgs ipa.UserAddOptionalArgs
	userAddOptionalArgsMap := make(map[string]interface{})
	userAddOptionalArgsMap["uid"] = d.Get("username")
	userAddOptionalArgsMap["loginshell"] = d.Get("shell")
	userAddOptionalArgsMap["all"] = true

	if val, ok := d.GetOk("homedirectory"); ok {
		userAddOptionalArgsMap["homedirectory"] = val
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
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("username", r.Result.UID); err != nil {
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

	if err := d.Set("shell", r.Result.Loginshell); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("homedirectory", r.Result.Homedirectory); err != nil {
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
		"homedirectory",
	) {
		var userModOptionalArgs ipa.UserModOptionalArgs
		userModOptionalArgsMap := make(map[string]interface{})
		userModOptionalArgsMap["uid"] = d.Id()
		userModOptionalArgsMap["firstname"] = d.Get("firstname")
		userModOptionalArgsMap["lastname"] = d.Get("lastname")
		userModOptionalArgsMap["shell"] = d.Get("shell")
		userModOptionalArgsMap["all"] = true
		if val, ok := d.GetOk("homedirectory"); ok {
			log.Printf("============= GOT HOMEDIRECTORY OK ===============")
			userModOptionalArgsMap["homedirectory"] = val
		}
		mapstructure.Decode(userModOptionalArgsMap, &userModOptionalArgs)
		r, err := c.UserMod(
			&ipa.UserModArgs{},
			&userModOptionalArgs,
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
