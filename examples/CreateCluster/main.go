package main

import (
    "context"
    "os"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
    "github.com/aws/smithy-go/ptr"
    "github.com/pingcap/tiup/pkg/tui"
)

func main() {

    projectID := os.Getenv("TIDBCLOUD_PROJECT_ID")
    if projectID == "" {
        panic("No project id is specified")
    }

    client, err := tidbcloud.NewDigestClientWithResponses()
    if err != nil {
        panic(err)
    }

    createClusterJSONRequestBody := tidbcloud.CreateClusterJSONRequestBody{
        CloudProvider: "AWS",
        ClusterType: "DEDICATED",
        Name: "APICluster",
        Region: "us-east-1",
    }

    createClusterJSONRequestBody.Config.Components = & struct{
        Tidb struct {
            NodeQuantity int32 `json:"node_quantity"`
            NodeSize string `json:"node_size"`
        } `json:"tidb"`
         Tiflash *struct {
             NodeQuantity int32 `json:"node_quantity"`
             NodeSize string `json:"node_size"`
             StorageSizeGib int32 `json:"storage_size_gib"`
         } `json:"tiflash,omitempty"`
        Tikv struct {
            NodeQuantity int32 `json:"node_quantity"`
            NodeSize string `json:"node_size"`
            StorageSizeGib int32 `json:"storage_size_gib"`
        } `json:"tikv"`
    }{
        struct {
            NodeQuantity int32 `json:"node_quantity"`
            NodeSize string `json:"node_size"`
        }{ 1, "2C8G"} ,
         nil,
        struct {
            NodeQuantity int32 `json:"node_quantity"`
            NodeSize string `json:"node_size"`
            StorageSizeGib int32 `json:"storage_size_gib"`
        }{ 3, "2C8G" ,200 },
    }

    createClusterJSONRequestBody.Config.IpAccessList =  &[]struct{
             Cidr string `json:"cidr"`
             Description *string `json:"description,omitempty"`
         }{{"0.0.0.0/0", ptr.String("All all destinations to access") }}

    createClusterJSONRequestBody.Config.RootPassword = "1234Abcd"
    createClusterJSONRequestBody.Config.Port = ptr.Int32(4000)

    var clusterInfo [][]string
    clusterInfo = append(clusterInfo, []string{"id", "Message"})
    response, err := client.CreateClusterWithResponse(context.Background(), projectID,  createClusterJSONRequestBody)
    if err != nil {
        panic(err)
    }

    statusCode := response.StatusCode()
    switch statusCode {
        case 200:
            clusterInfo = append(clusterInfo, []string{response.JSON200.Id, "Succeeded in creating the cluster"})
        case 400:
            clusterInfo = append(clusterInfo, []string{"-", "Exsited resource"})
    }
    tui.PrintTable(clusterInfo, true)
}
