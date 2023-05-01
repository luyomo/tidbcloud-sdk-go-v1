package main

import (
    "context"
    "fmt"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
    "github.com/pingcap/tiup/pkg/tui"
)

func main() {
    client, err := tidbcloud.NewDigestClientWithResponses()
    if err != nil {
        panic(err)
    }

    response, err := client.ListProviderRegionsWithResponse(context.Background() )
    if err != nil {
        panic(err)
    }

    var specInfo [][]string
    specInfo = append(specInfo, []string{"Cluster Type", "Cloud Provider", "Region",  "Component Type", "Node Size", "Node Quantity Min", "Node Quantity Step", "Storage MIN", "Storage MAX" })
    for _, item := range *response.JSON200.Items {
        for _, component := range *item.Tidb {
            specInfo = append(specInfo, []string{(*item.ClusterType).(string), (*item.CloudProvider).(string), *item.Region, "TiDB", *component.NodeSize, fmt.Sprintf("%d", *component.NodeQuantityRange.Min), fmt.Sprintf("%d",  *component.NodeQuantityRange.Step), "-", "-"})
        }

        for _, component := range *item.Tikv {
            specInfo = append(specInfo, []string{(*item.ClusterType).(string), (*item.CloudProvider).(string), *item.Region, "TiKV", *component.NodeSize, fmt.Sprintf("%d", *component.NodeQuantityRange.Min), fmt.Sprintf("%d",  *component.NodeQuantityRange.Step), fmt.Sprintf("%d",  *component.StorageSizeGibRange.Min), fmt.Sprintf("%d",  *component.StorageSizeGibRange.Max)  })
        }

        for _, component := range *item.Tiflash {
            specInfo = append(specInfo, []string{(*item.ClusterType).(string), (*item.CloudProvider).(string), *item.Region, "TiFlash", *component.NodeSize, fmt.Sprintf("%d", *component.NodeQuantityRange.Min), fmt.Sprintf("%d",  *component.NodeQuantityRange.Step), fmt.Sprintf("%d",  *component.StorageSizeGibRange.Min), fmt.Sprintf("%d",  *component.StorageSizeGibRange.Max)  })
        }
    }

    tui.PrintTable(specInfo, true) 
}
