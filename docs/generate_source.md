# Background
TiDB Cloud provides API to manage TiDB Cluser like cluster creation, decommission, data import etc. The curl is used to send the http request. If you want to manage the cluster with golang, you need one SDK for easy programming.
This topic will introduce how to prepare the golang SDK using [oapi-codegen](https://github.com/deepmap/oapi-codegen) within 10 minutes.

# swagger API definition
TiDB Cloud use swagger 2.0 to define the API service. Please find the [link](https://download.pingcap.org/tidbcloud-oas.json) for the data source. Next go to [swagger edit](https://editor.swagger.io/#/) and import the API definition using File/Import URL. Once import the TiDB Cloud API definition, click Edit/Convert to openapi 3.0 to convert swagger 2.0 to openapi 3.0 which is required by [oapi-codegen](https://github.com/deepmap/oapi-codegen). It takes several minutes to convert the specification and download it by
clicking [Save as YAML]. Name this download file as openapi

## Manual conversion
The enum has to been adjusted in the download file as below:
### original definition
```
... 
  items:
    title: ListProviderRegionsItem
      type: object
      properties:
        cluster_type:
          description: |-
            The cluster type.
            - `"DEVELOPER"`: a [Serverless Tier](https://docs.pingcap.com/tidbcloud/select-cluster-tier#serverless-tier) cluster
            - `"DEDICATED"`: a [Dedicated Tier](https://docs.pingcap.com/tidbcloud/select-cluster-tier#dedicated-tier) cluster

            **Warning:** `"DEVELOPER"` will soon be changed to `"SERVERLESS"` to represent Serverless Tier clusters.
          example: DEDICATED
          format: enum
          type: string
          enum:
          - DEDICATED
          - DEVELOPER
```
### adjusted definition
```
... 
  items:
    title: ListProviderRegionsItem
      type: object
      properties:
        cluster_type:
          description: |-
            The cluster type.
            - `"DEVELOPER"`: a [Serverless Tier](https://docs.pingcap.com/tidbcloud/select-cluster-tier#serverless-tier) cluster
            - `"DEDICATED"`: a [Dedicated Tier](https://docs.pingcap.com/tidbcloud/select-cluster-tier#dedicated-tier) cluster

            **Warning:** `"DEVELOPER"` will soon be changed to `"SERVERLESS"` to represent Serverless Tier clusters.
          example: DEDICATED
          schema:
            format: enum
            type: string
            enum:
            - DEDICATED
            - DEVELOPER
```

# Source generation
## Download oapi-codegen
Use the below command to install the oapi-codegen and set the search path to download directory. The version is v1.12.4 when I tested. 
```
workstation$ go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
workstation$  /home/pi/go/bin/oapi-codegen --version
v1.12.4
```
## Generate golang source code for TiDb Cloud API
```
workstation$ oapi-codegen -package client openapi.yaml > client/tidbcloud.go
```

## Add http digest authentication
The generated source code does not have any http authentication. Since TiDB Cloud provides the http digest authentication, it has to been implemented. Import the auth.go and replace [ return c.Client.Do(req) ] with [return c.AddDigestHeader(req)]
```
//return c.Client.Do(req)
return c.AddDigestHeader(req)
```

# Hello world - List all the projects
```
package main

import (
    "fmt"
    "context"
    cli "tidbcloud/client"
)

func main() {}
    client, err := cli.NewClientWithResponses("https://api.tidbcloud.com")
    if err != nil {
        panic(err)
    }
 
    response, err := client.ListProjectsWithResponse(context.Background(), &cli.ListProjectsParams{})
    if err != nil {
        panic(err)
    }
 
    for _, item := range response.JSON200.Items {
         fmt.Printf("id: <%#v>, Name: <%#v>, AwsCmekEnabled: <%#v>, ClusterCount: <%#v>, CreateTimestamp: <%#v>, OrgId: <%#v> \n", *item.Id, *item.Name, *item.AwsCmekEnabled, *     item.ClusterCount, *item.CreateTimestamp, *item.OrgId)
    }
}
```
