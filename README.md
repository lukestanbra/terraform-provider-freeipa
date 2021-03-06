# To do

- ~~Makefile~~
- Terratest
- ~~Documentation via tfplugindocs~~
- CHANGELOG.md
- goreleaser?
- LICENSE
- Remaining resources:
  - aci
  - Auto Membership Rule
  - automember_task
  - Automount Keys
  - Automount Locations
  - Automount Maps
  - Certificate Authorities
  - CA ACLs
  - certmap
  - Certificate Identity Mapping Global Configuration
  - Certificate Identity Mapping Rules
  - Certificate Profiles
  - certreq
  - class
  - command
  - Configuration
  - Entry
  - Delegations
  - dns_system_records
  - dnsa6record
  - dnsaaaarecord
  - dnsafsdbrecord
  - dnsaplrecord
  - dnsarecord
  - dnscertrecord
  - dnscnamerecord
  - DNS Global Configuration
  - dnsdhcidrecord
  - dnsdlvrecord
  - dnsdnamerecord
  - dnsdsrecord
  - DNS Forward Zones
  - dnshiprecord
  - dnsipseckeyrecord
  - dnskeyrecord
  - dnskxrecord
  - dnslocrecord
  - dnsmxrecord
  - dnsnaptrrecord
  - dnsnsecrecord
  - dnsnsrecord
  - dnsptrrecord
  - DNS Resource Records
  - dnsrprecord
  - dnsrrsigrecord
  - DNS Servers
  - dnssigrecord
  - dnsspfrecord
  - dnssrvrecord
  - dnssshfprecord
  - dnstlsarecord
  - dnstxtrecord
  - dnsurirecord
  - DNS Zones
  - User Groups
  - HBAC Rules
  - HBAC Services
  - HBAC Service Groups
  - Hosts
  - Host Groups
  - Group ID overrides
  - User ID overrides
  - ID Ranges
  - ID Views
  - Kerberos Ticket Policy
  - IPA Locations
  - metaobject
  - Netgroups
  - OTP Configuration
  - OTP Tokens
  - output
  - param
  - Permissions
  - pkinit
  - Privileges
  - Password Policies
  - RADIUS Servers
  - Realm Domains
  - Roles
  - Self Service Permissions
  - SELinux User Maps
  - IPA Servers
  - server_role
  - Services
  - Service delegation rules
  - Service delegation targets
  - servrole
  - Stage Users
  - Sudo Commands
  - Sudo Command Groups
  - Sudo Rules
  - topic
  - Topology Segments
  - Topology suffixes
  - Trusts
  - Global Trust Configuration
  - Trusted domains
  - ~~Users~~
  - Entry
  - Vaults
  - vaultconfig
  - Vault Containers

# Bugs

- ~~Import not pulling in all attributes~~
- Datetime variables do not work properly
  - At some point in the last couple years, the FreeIPA API started returning Datetime
  variables like so:
  ```
  "krbpasswordexpiration": [
    {
      "__datetime__": "20210926115531Z"
    }
  ]
  ```
  The go-freeipa library assumed that it could cast this to a `time.Time`, however this doesn't work. Current workaround is that I've edited the package with a regenerated version of the freeipa client code so that these are just `interface{}`. I'm also thinking I can mostly avoid ever having to use Datetime fields as these won't ever be set declaratively. For this reason I've forked the go-freeipa repo with these changes in.


# Development

Getting a FreeIPA container going is an absolute pain on a Mac.

```
make container
```

This sets up FreeIPA with the user `admin` and password `password`. You can connect to it using the UI at https://ipa.example.test.

Copy the cert off out of the container and add it to the mac keychain.

```
make certificate
```

With this container running locally, you can run the example terraform against it. First however, you must build the provider:

```
make install
```

To run the example:

```
make example
```

If you update the provider code, you'll need to do the following for the core loop:

```
make test
make testacc
make install
make example
```

Regenerate the documentation via

```
go generate
```
