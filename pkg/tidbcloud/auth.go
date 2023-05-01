package tidbcloud

import (
	"fmt"
    "os"
	"net/http"
	"strings"
    "errors"

    "gopkg.in/ini.v1"
    "github.com/icholy/digest"
)

const (
    API_URL string = "https://api.tidbcloud.com"
)

func NewDigestClientWithResponses(opts ...ClientOption) (*ClientWithResponses, error) {
    client, err := NewDigestClient(API_URL, opts...)
    if err != nil {
        return nil, err
    }
    return &ClientWithResponses{client}, nil
}

// Clone from NewClient from cli
func NewDigestClient(server string, opts ...ClientOption) (*Client, error) {
    // create a client with sane default values
    client := Client{
        Server: server,
    }
    // mutate client and add all optional params
    for _, o := range opts {
        if err := o(&client); err != nil {
            return nil, err
        }
    }
    // ensure the server URL always has a trailing slash
    if !strings.HasSuffix(client.Server, "/") {
        client.Server += "/"
    }

    name, password, err := readCredentials()
    if err != nil {
        return nil, err
    }

    // create httpClient, if not already present
    if client.Client == nil {
        client.Client = &http.Client{
            Transport: &digest.Transport{
                Username: name,
                Password: password,
            },
        }
    }
    return &client, nil
}

func readCredentials() (name, password string, err error) {
    // ~/.tidbcloud/credentials

    // 01. Read env variables
    name = os.Getenv("TIDBCLOUD_PUBLIC_KEY")
    password = os.Getenv("TIDBCLOUD_PRIVATE_KEY")

    if name != "" && password != "" {
        return
    }

    // 02. Read credentials from customized file
    credentialsFile := os.Getenv("TIDBCLOUD_CREDENTIAL_FILE")

    // 03. Read credentials from default file: ~/.tidbcloud/credentials
    var dirname string
    if credentialsFile == "" {
        dirname, err = os.UserHomeDir()
        if err != nil {
            return
        }
        credentialsFile = fmt.Sprintf("%s/.tidbcloud/credentials", dirname)
    }

    if _, err = os.Stat(credentialsFile); err == nil {
        cfg, inierr := ini.Load(fmt.Sprintf("%s/.tidbcloud/credentials", dirname))
	    if inierr != nil {
            err = inierr
            return
	    }

	    name = cfg.Section("credential").Key("public_key").String()
	    password = cfg.Section("credential").Key("private_key").String()
    } else {
        errors.New("Credentials not found")
    }

    return
}

