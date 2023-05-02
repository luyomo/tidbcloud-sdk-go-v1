package main

import (
    "context"
    "os"
    "fmt"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
)

func main() {

    projectID := os.Getenv("TIDBCLOUD_PROJECT_ID")
    if projectID == "" {
        panic("No project id is specified")
    }

    clusterName := os.Getenv("TIDBCLOUD_CLUSTER_NAME")
    if projectID == "" {
         panic("No cluster name is specified")
    }

    client, err := tidbcloud.NewDigestClientWithResponses()
    if err != nil {
        panic(err)
    }

    response, err := client.ListClustersOfProjectWithResponse(context.Background(), projectID,  &tidbcloud.ListClustersOfProjectParams{})
    if err != nil {
        panic(err)
    }

    clusterID := ""
    for _, item := range response.JSON200.Items {
        if *item.Name == clusterName {
            clusterID = item.Id 
        }
    }

    if clusterID == "" {
        panic("No valid cluster id found")
    }


    resDelete, err := client.DeleteClusterWithResponse(context.Background(), projectID, clusterID)
    if err != nil {
        panic(err)
    }

    statusCode := resDelete.StatusCode()
    switch statusCode {
        case 200:
            fmt.Printf("Succeeded in deleting the cluster")
        case 400:
            fmt.Printf("Failed to delete the cluster: %#v \n", *resDelete.JSON400.Message)
    }
}
