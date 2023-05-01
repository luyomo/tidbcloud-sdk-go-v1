# How to run the example
```
workstation$ cd tidbcloud-sdk-go-v1/examples/ListProviderRegions
workstation$ go run main.go
Cluster Type  Cloud Provider  Region           Component Type  Node Size  Node Quantity Min  Node Quantity Step  Storage MIN  Storage MAX
------------  --------------  ------           --------------  ---------  -----------------  ------------------  -----------  -----------
DEDICATED     AWS             ap-northeast-2   TiDB            2C8G       1                  1                   -            -   
DEDICATED     AWS             ap-northeast-2   TiDB            4C16G      1                  1                   -            -   
DEDICATED     AWS             ap-northeast-2   TiDB            8C16G      1                  1                   -            -   
DEDICATED     AWS             ap-northeast-2   TiDB            16C32G     1                  1                   -            -   
DEDICATED     AWS             ap-northeast-2   TiKV            2C8G       3                  3                   200          500 
DEDICATED     AWS             ap-northeast-2   TiKV            4C16G      3                  3                   200          2048
DEDICATED     AWS             ap-northeast-2   TiKV            8C32G      3                  3                   200          4096
DEDICATED     AWS             ap-northeast-2   TiKV            8C64G      3                  3                   200          4096
DEDICATED     AWS             ap-northeast-2   TiKV            16C64G     3                  3                   200          4096
DEDICATED     AWS             ap-northeast-2   TiFlash         8C64G      0                  1                   200          2048
DEDICATED     AWS             ap-northeast-2   TiFlash         16C128G    0                  1                   200          2048
```
