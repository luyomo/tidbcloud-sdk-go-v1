package main

import (
    "fmt"
    "context"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
    "github.com/pingcap/tiup/pkg/tui"
)

func main() {
    client, err := tidbcloud.NewClientWithResponses("https://api.tidbcloud.com")
    if err != nil {
        panic(err)
    }

    response, err := client.ListProjectsWithResponse(context.Background(), &tidbcloud.ListProjectsParams{})
    if err != nil {
        panic(err)
    }

    var clusterInfo [][]string
    clusterInfo = append(clusterInfo, []string{"id", "Name", "AwsCmekEnabled", "ClusterCount", "CreateTimestamp", "OrgId"})
    for _, item := range response.JSON200.Items {
        clusterInfo = append(clusterInfo, []string{*item.Id, *item.Name, fmt.Sprintf("%t", *item.AwsCmekEnabled), fmt.Sprintf("%d", *item.ClusterCount), *item.CreateTimestamp, *item.OrgId})
    }

    tui.PrintTable(clusterInfo, true) 
}
