package main

import (
    "context"
    "os"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
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

    response, err := client.ListClustersOfProjectWithResponse(context.Background(), projectID,  &tidbcloud.ListClustersOfProjectParams{})
    if err != nil {
        panic(err)
    }

    var clusterInfo [][]string
    clusterInfo = append(clusterInfo, []string{"Project ID", "id", "Name", "Cluster Type", "Cloud Provider", "Region", "Create Timestamp", "Status", "TiDB Version"})
    for _, item := range response.JSON200.Items {
         clusterInfo = append(clusterInfo, []string{item.ProjectId, item.Id, *item.Name, (*item.ClusterType).(string), (*item.CloudProvider).(string), *item.Region, *item.CreateTimestamp, (*item.Status.ClusterStatus).(string), (*item.Status.TidbVersion)})
    }

    tui.PrintTable(clusterInfo, true) 
}
