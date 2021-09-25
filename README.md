# To do

- Makefile
- Terratest
- Documentation via tfplugindocs
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
  - Users
  - Entry
  - Vaults
  - vaultconfig
  - Vault Containers

# Development

Getting a FreeIPA container going is an absolute pain on a Mac.

```
docker volume create freeipa-data
docker volume create freeipa-run
docker run -it --name freeipa-server-container -h ipa.example.test -v /sys/fs/cgroup:/sys/fs/cgroup:ro --mount 'type=volume,src=freeipa-data,dst=/data' --mount 'type=volume,src=freeipa-run,dst=/run' --tmpfs /tmp --sysctl net.ipv6.conf.all.disable_ipv6=0 -p 127.0.0.1:443:443 -p 127.0.0.1:80:80 freeipa/freeipa-server:centos-8-4.9.2
```

Go through all the gumph, select all the defaults, and choose a simple password, e.g. password.

With this container running locally, you can run the example terraform against it. First however, you must build the provider:

```
go get
go mod tidy
go mod vendor
go build -o ~/.terraform.d/plugins/hashicorp.com/lukestanbra/freeipa/0.0.1/darwin_amd64/terraform-provider-freeipa
```

To run the example:

```
pushd examples
terraform init
terraform apply
popd
```

If you update the provider code, you'll need to do the following for the core loop:

```
rm -rf examples/.terraform*
go build -o ~/.terraform.d/plugins/hashicorp.com/lukestanbra/freeipa/0.0.1/darwin_amd64/terraform-provider-freeipa
pushd examples
terraform init
terraform apply
popd
```

Regenerate the documentation via

```
go generate
```
