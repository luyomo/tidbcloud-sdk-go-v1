# tidbcloud-sdk-go-v1

Translations:

* [简体中文](README_zh.md)
* [日本語](README_ja.md)

## Overview
tidbcloud-sdk-go-v1 is the TiDB Cloud SDK for the Go programming language. It implements [TiDB Cloud API](https://docs.pingcap.com/tidbcloud/api/v1beta).

### Hello world
  + Once you register one TiDB Cloud account, please follow the [TiDB Cloud document](https://docs.pingcap.com/tidbcloud/api-overview) to get the API key.
  + Set the credential config file as below. 
      - Config the credentials into default file
          ```
          workstation$ more ~/.tidbcloud/credentials
          [credential]
          public_key=xxxxxxxx
          private_key=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx 
          ```
      - Config the credentials into customized file
          ```
          workstation$ export TIDBCLOUD_CREDENTIAL_FILE=/tmp/credentials
          workstation$ more ~/credentials
          [credential]
          public_key=xxxxxxxx
          private_key=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx 
          ```
      - Set env variables
          ```
          workstation$ export TIDBCLOUD_PUBLIC_KEY=xxxxxxxx
          workstation$ export TIDBCLOUD_PRIVATE_KEY=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
          ```
  + Run first example
      ```
      workstation$ cd examples/ListProjects
      workstation$ go run main.go
      id                   Name           AwsCmekEnabled  ClusterCount  CreateTimestamp  OrgId
      --                   ----           --------------  ------------  ---------------  -----
      1111111111111111111  ProjectName    false           1             1678671443       2222222222222222222
      ```
