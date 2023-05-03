package main

import (
    "context"
    "os"
    "fmt"

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
         panic("No cluster name is specified")
    }

    awsKeyID := os.Getenv("AWS_KEY_ID")
    if awsKeyID == "" {
         panic("No aws key is specified")
    }

    awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
    if awsSecretAccessKey == "" {
         panic("No aws secret is specified")
    }

    s3uri := os.Getenv("S3_URI")
    if s3uri == "" {
         panic("No S3 URI is specified")
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


    var createImportTaskJSONRequestBody tidbcloud.CreateImportTaskJSONRequestBody
    createImportTaskJSONRequestBody.Name = ptr.String("import-from-api")
    createImportTaskJSONRequestBody.Spec.Source.Type = "S3"
    createImportTaskJSONRequestBody.Spec.Source.Format.Type = "CSV"
    createImportTaskJSONRequestBody.Spec.Source.Uri = s3uri
    createImportTaskJSONRequestBody.Spec.Source.AwsKeyAccess = &struct{
                 AccessKeyId string `json:"access_key_id"`
                 SecretAccessKey string `json:"secret_access_key"`
             } {awsKeyID, awsSecretAccessKey}

    resImport, err := client.CreateImportTaskWithResponse(context.Background(), projectID, clusterID, createImportTaskJSONRequestBody)
    if err != nil {
        panic(err)
    }

     statusCode := resImport.StatusCode()
     switch statusCode {
         case 200:
             fmt.Printf("Started the import job")
         case 400:
             fmt.Printf("Failed to import data: %#v \n", *resImport.JSON400.Message)
     }
}
