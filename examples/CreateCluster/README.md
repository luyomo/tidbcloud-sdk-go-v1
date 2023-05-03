# How to run the example
```
export TIDB_PROJECT_ID=1111111111111111111
workstation$ cd tidbcloud-sdk-go-v1/examples/CreateCluster
workstation$ go run main.go
id                  Message
--                  --
ddddddddddddddddddd Succeeded in creating the cluster
```

# Default value
| Column Name   | Value       |
| CloudProvider | AWS         |
| ClusterType   | DEDICATED   |
| Name          | APICluster  |
| Region        | us-east-1   |
| IpAccessList  | 0.0.0.0/0   |
| TiDB          | 1/2C8G      |
| TiKV          | 3/2C8G/200G |
