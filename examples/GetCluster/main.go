package main

import (
    "fmt"
    "context"
    "os"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
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
            getClusterRes, err := client.GetClusterWithResponse(context.Background(), projectID, item.Id )
            if err != nil {
                panic(err)
            }
            fmt.Printf("Response: <%#v>", getClusterRes.JSON200)
        }
    }
}
