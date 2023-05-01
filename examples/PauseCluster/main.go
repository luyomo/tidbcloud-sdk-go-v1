package main

import (
    "fmt"
    "context"
    "os"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
    "github.com/aws/smithy-go/ptr"
)

func main() {
    projectID := os.Getenv("TIDBCLOUD_PROJECT_ID")
    if projectID == "" {
         panic("No project id is specified")
    }
    clusterName := os.Getenv("TIDBCLOUD_CLUSTER_NAME")
    if projectID == "" {
         panic("No project id is specified")
    }

    client, err := tidbcloud.NewDigestClientWithResponses()
    if err != nil {
        panic(err)
    }

    response, err := client.ListClustersOfProjectWithResponse(context.Background(), projectID,  &tidbcloud.ListClustersOfProjectParams{})
    if err != nil {
        panic(err)
    }

    for _, item := range response.JSON200.Items {
        if *item.Name == clusterName && (*item.Status.ClusterStatus).(string) == "AVAILABLE" {
            var updateClusterJSONRequestBody tidbcloud.UpdateClusterJSONRequestBody 
            updateClusterJSONRequestBody.Config.Paused = ptr.Bool(true)
            updateRes, err := client.UpdateClusterWithResponse(context.Background(), projectID, item.Id, updateClusterJSONRequestBody)
            if err != nil {
                panic(err)
            }
        }
    }
}
