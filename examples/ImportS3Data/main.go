package main

import (
    "context"
    "os"
    "fmt"

    "github.com/luyomo/tidbcloud-sdk-go-v1/pkg/tidbcloud"
    "github.com/aws/smithy-go/ptr"
//    "github.com/pingcap/tiup/pkg/tui"
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
    fmt.Printf("The cluster id: <%s> \n", clusterID)


    var createImportTaskJSONRequestBody tidbcloud.CreateImportTaskJSONRequestBody
    createImportTaskJSONRequestBody.Name = ptr.String("import-from-api")
    createImportTaskJSONRequestBody.Spec.Source.Type = "S3"
    createImportTaskJSONRequestBody.Spec.Source.Format.Type = "CSV"
    createImportTaskJSONRequestBody.Spec.Source.Uri = "s3://jay-ticdc/csvdata/"
    createImportTaskJSONRequestBody.Spec.Source.AwsKeyAccess = &struct{
                 AccessKeyId string `json:"access_key_id"`
                 SecretAccessKey string `json:"secret_access_key"`
             } {"AKIA2TXTRGT4RX4HUT7W", "IK472KMfu5ND96fI4swpBFgtNh1sp5HYpUBy0QVr"}

    resImport, err := client.CreateImportTaskWithResponse(context.Background(), projectID, clusterID, createImportTaskJSONRequestBody)
    if err != nil {
        panic(err)
    }

    fmt.Printf("response: <%#v> \n", *resImport.JSON400.Message)
    statusCode := resImport.StatusCode()
    fmt.Printf("Status code: %#v \n", statusCode)

//     statusCode := response.StatusCode()
//     switch statusCode {
//         case 200:
//             clusterInfo = append(clusterInfo, []string{response.JSON200.Id, "Succeeded in creating the cluster"})
//         case 400:
//             clusterInfo = append(clusterInfo, []string{"-", "Exsited resource"})
//     }
//     tui.PrintTable(clusterInfo, true)
}
