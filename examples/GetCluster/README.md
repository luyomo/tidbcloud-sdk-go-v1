# How to run the example
```
export TIDBCLOUD_PROJECT_ID=1111111111111111111
export TIDBCLOUD_CLUSTER_NAME=APICluster 
workstation$ cd tidbcloud-sdk-go-v1/examples/ListProjects
workstation$ go run main.go
id                   Name           AwsCmekEnabled  ClusterCount  CreateTimestamp  OrgId
--                   ----           --------------  ------------  ---------------  -----
1111111111111111111  ProjectName    false           1             1678671443       2222222222222222222
```
