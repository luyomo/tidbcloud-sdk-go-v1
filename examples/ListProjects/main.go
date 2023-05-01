package main

import (
    "fmt"
    "context"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
    "github.com/pingcap/tiup/pkg/tui"
)

func main() {
    client, err := tidbcloud.NewDigestClientWithResponses()
    if err != nil {
        panic(err)
    }

    response, err := client.ListProjectsWithResponse(context.Background(), &tidbcloud.ListProjectsParams{})
    if err != nil {
        panic(err)
    }

    var projectInfo [][]string
    projectInfo = append(projectInfo, []string{"id", "Name", "AwsCmekEnabled", "ClusterCount", "CreateTimestamp", "OrgId"})
    for _, item := range response.JSON200.Items {
        projectInfo = append(projectInfo, []string{*item.Id, *item.Name, fmt.Sprintf("%t", *item.AwsCmekEnabled), fmt.Sprintf("%d", *item.ClusterCount), *item.CreateTimestamp, *item.OrgId})
    }

    tui.PrintTable(projectInfo, true) 
}
