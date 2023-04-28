package main

import (
    "context"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
    "github.com/pingcap/tiup/pkg/tui"
)

func main() {
    client, err := tidbcloud.NewClientWithResponses("https://api.tidbcloud.com")
    if err != nil {
        panic(err)
    }

    response, err := client.ListClustersOfProjectWithResponse(context.Background(), "1111111111111111111",  &tidbcloud.ListClustersOfProjectParams{})
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
