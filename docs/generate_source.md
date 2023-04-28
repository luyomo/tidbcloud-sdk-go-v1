# Background
[TiDB Cloud](https://www.pingcap.com/tidb-cloud/) provides API to manage TiDB Cluser like cluster creation, decommission, data import etc. The curl is generally used to send the http request to perform those operations. If you want to manage the clusters with golang, you need one SDK for easy programming.
This topic will introduce how to prepare the tidbcloud golang SDK using [oapi-codegen](https://github.com/deepmap/oapi-codegen) within 10 minutes.

# swagger API definition
TiDB Cloud use swagger 2.0 to define the API specification. Please find the [link](https://download.pingcap.org/tidbcloud-oas.json) for the json file. Next go to [swagger edit](https://editor.swagger.io/#/) and import the API specification definition using File/Import URL. Once the specification is imported into the editor, click (Edit/Convert to openapi 3.0) to convert swagger 2.0 to openapi 3.0 which is required by [oapi-codegen](https://github.com/deepmap/oapi-codegen). It takes several minutes for conversion and download it by clicking [Save as YAML] after it is available. Here name this downloaded file as openapi

## Manual conversion
The enum has to been adjusted in the download file as below, otherwise the compilation will fail:
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
          <span style="color:blue">schema</span>:
            format: enum
            type: string
            enum:
            - DEDICATED
            - DEVELOPER
```

# Source generation
## Download oapi-codegen
Use the below command to install the oapi-codegen and set the search path to download directory. The version is v1.12.4 when this test is performed.
```
workstation$ go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
workstation$  /home/pi/go/bin/oapi-codegen --version
v1.12.4
```
## Generate golang source code for TiDb Cloud API
In this case, the generated file is put into [pkg/tidbcloud](https://github.com/luyomo/tidbcloud-sdk-go-v1/tree/master/pkg/tidbcloud).
```
workstation$ oapi-codegen -package tidbcloud openapi.yaml > tidbcloud/cli.go
```

## Add http digest authentication
The generated source code does not have any http authentication. Since TiDB Cloud provides the http digest authentication, it has to been implemented into this SDK. Import the [auth.go](https://github.com/luyomo/tidbcloud-sdk-go-v1/blob/master/pkg/tidbcloud/auth.go) and replace all the [return c.Client.Do(req)] with [return c.AddDigestHeader(req)] in the [cli.go](https://github.com/luyomo/tidbcloud-sdk-go-v1/blob/master/pkg/tidbcloud/cli.go) file.
```
//return c.Client.Do(req)
return c.AddDigestHeader(req)
```

# Hello world - List all the projects
Please find the first example of [projects list](https://github.com/luyomo/tidbcloud-sdk-go-v1/tree/master/examples/ListProjects).
