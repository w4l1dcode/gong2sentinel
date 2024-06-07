# gong2sentinel

A Go program that exports Gong audit logs and user permissions on calls to Microsoft Sentinel SIEM.
Two tables are used; `GongAuditLogs` for audit logs and `GongCallUserAccess` for  user permissions on calls.

## Running

First create a yaml file, such as `dev.yml`:
```yaml
log:
  level: DEBUG

microsoft:
  app_id: ""
  secret_key: ""
  tenant_id: ""
  subscription_id: ""
  resource_group: ""
  workspace_name: ""
  expires_months: 
  update_table: 
  dcr:
    endpoint: ""
    rule_id: ""
    stream_name_auditing: ""
    stream_name_user_access: ""

gong:
  access_key: ""
  access_secret: ""
  lookup_hours: """
```

And now run the program from source code:
```shell
% make
go run ./cmd/... -config=dev.yml
INFO[0000] shipping logs                                 module=sentinel_logs table_name=GongAuditLogs total=82
INFO[0002] shipped logs                                  module=sentinel_logs table_name=GongAuditLogs
INFO[0002] successfully sent logs to sentinel            total=82
```

## Building

```shell
% make build
```
