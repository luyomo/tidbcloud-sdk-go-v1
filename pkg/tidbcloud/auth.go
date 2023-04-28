package tidbcloud

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
    "os"
	"log"
	"net/http"
	"net/url"
	"strings"
    "errors"

    "gopkg.in/ini.v1"
)

type myjar struct {
	jar map[string][]*http.Cookie
}

// DigestHeaders tracks the state of authentication
type DigestHeaders struct {
	Realm     string
	Qop       string
	Method    string
	Nonce     string
	Opaque    string
	Algorithm string
	HA1       string
	HA2       string
	Cnonce    string
	Path      string
	Nc        int16
	Username  string
	Password  string
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

func (d *DigestHeaders) digestChecksum() {
	switch d.Algorithm {
	case "MD5":
		// A1
		h := md5.New()
		A1 := fmt.Sprintf("%s:%s:%s", d.Username, d.Realm, d.Password)
		io.WriteString(h, A1)
		d.HA1 = fmt.Sprintf("%x", h.Sum(nil))

		// A2
		h = md5.New()
		A2 := fmt.Sprintf("%s:%s", d.Method, d.Path)
		io.WriteString(h, A2)
		d.HA2 = fmt.Sprintf("%x", h.Sum(nil))
	case "MD5-sess":
		// A1
		h := md5.New()
		A1 := fmt.Sprintf("%s:%s:%s", d.Username, d.Realm, d.Password)
		io.WriteString(h, A1)
		haPre := fmt.Sprintf("%x", h.Sum(nil))
		h = md5.New()
		A1 = fmt.Sprintf("%s:%s:%s", haPre, d.Nonce, d.Cnonce)
		io.WriteString(h, A1)
		d.HA1 = fmt.Sprintf("%x", h.Sum(nil))

		// A2
		h = md5.New()
		A2 := fmt.Sprintf("%s:%s", d.Method, d.Path)
		io.WriteString(h, A2)
		d.HA2 = fmt.Sprintf("%x", h.Sum(nil))
	default:
		//token
	}
}

func (d *Client) ReadCredentials() (name, password string, err error) {
  // ~/.tidbcloud/credentials

    // 01. Read env variables
    name = os.Getenv("TIDBCLOUD_PUBLIC_KEY")
    password = os.Getenv("TIDBCLOUD_PRIVATE_KEY")

    if name != "" && password != "" {
        return
    }

    credentialsFile := os.Getenv("TIDBCLOUD_CREDENTIAL_FILE")

    var dirname string
    if credentialsFile == "" {
        dirname, err = os.UserHomeDir()
        if err != nil {
            return
        }
        credentialsFile = fmt.Sprintf("%s/.tidbcloud/credentials", dirname)
    }

    if _, err := os.Stat(credentialsFile); err == nil {
        cfg, err := ini.Load(fmt.Sprintf("%s/.tidbcloud/credentials", dirname))
	    if err != nil {
            return name, password, err
	    }

	    name = cfg.Section("credential").Key("public_key").String()
	    password = cfg.Section("credential").Key("private_key").String()
    } else {
        errors.New("Credentials not found")
    }

    return
}

// ApplyAuth adds proper auth header to the passed request
func (d *DigestHeaders) ApplyAuth(req *http.Request) {
	d.Nc += 0x1
	d.Cnonce = randomKey()
	d.Method = req.Method
	d.Path = req.URL.RequestURI()
	d.digestChecksum()
	response := h(strings.Join([]string{d.HA1, d.Nonce, fmt.Sprintf("%08x", d.Nc),
		d.Cnonce, d.Qop, d.HA2}, ":"))
	AuthHeader := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc=%08x, qop=%s, response="%s", algorithm=%s`,
		d.Username, d.Realm, d.Nonce, d.Path, d.Cnonce, d.Nc, d.Qop, response, d.Algorithm)
	if d.Opaque != "" {
		AuthHeader = fmt.Sprintf(`%s, opaque="%s"`, AuthHeader, d.Opaque)
	}
	req.Header.Set("Authorization", AuthHeader)
}

func(c *Client)AddDigestHeader(req *http.Request) (*http.Response, error) {
    name, password, err := c.ReadCredentials()
    if err != nil {
        return nil, err
    }

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
    
    // 02. If failed as unauthorized request
	if resp.StatusCode == 401 {
		authn := digestAuthParams(resp)
		algorithm := authn["algorithm"]
		d := &DigestHeaders{}
		d.Path = req.URL.Path
		d.Realm = authn["realm"]
		d.Qop = authn["qop"]
		d.Nonce = authn["nonce"]
		d.Opaque = authn["opaque"]
		if algorithm == "" {
			d.Algorithm = "MD5"
		} else {
			d.Algorithm = authn["algorithm"]
		}
		d.Nc = 0x0
		d.Username = name
		d.Password = password

		d.ApplyAuth(req)

        // Send second request
		resp, err = c.Client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != 200 {
			d = &DigestHeaders{}
			err = fmt.Errorf("response status code was %v", resp.StatusCode)
		}

		return resp, err
    }
    return nil, nil
}

// Auth authenticates against a given URI
func (d *DigestHeaders) Auth(username string, password string, uri string) (*DigestHeaders, error) {

	client := &http.Client{}
	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	client.Jar = jar

    // 01. Create the first request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatal(err)
	}
    fmt.Printf("Request: <%#v> \n", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

    // 02. If failed as unauthorized request
	if resp.StatusCode == 401 {

        // 03. Parse result from request and prepare digest header
		authn := digestAuthParams(resp)
		algorithm := authn["algorithm"]
		d := &DigestHeaders{}
		u, _ := url.Parse(uri)
		d.Path = u.RequestURI()
		d.Realm = authn["realm"]
		d.Qop = authn["qop"]
		d.Nonce = authn["nonce"]
		d.Opaque = authn["opaque"]
		if algorithm == "" {
			d.Algorithm = "MD5"
		} else {
			d.Algorithm = authn["algorithm"]
		}
		d.Nc = 0x0
		d.Username = username
		d.Password = password

        // 04. Set the auth into request
		// req, err = http.NewRequest("GET", uri, nil)
		d.ApplyAuth(req)

        // 05. Send second request
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != 200 {
			d = &DigestHeaders{}
			err = fmt.Errorf("response status code was %v", resp.StatusCode)
		}
        fmt.Printf("RESP: <%#v> \n", resp.Body)

        b, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Fatalln(err)
        }
        fmt.Println(string(b))

		return d, err
	}
	return nil, fmt.Errorf("response status code should have been 401, it was %v", resp.StatusCode)
}

/*
Parse Authorization header from the http.Request. Returns a map of
auth parameters or nil if the header is not a valid parsable Digest
auth header.
*/
func digestAuthParams(r *http.Response) map[string]string {
	s := strings.SplitN(r.Header.Get("Www-Authenticate"), " ", 2)
	if len(s) != 2 || s[0] != "Digest" {
		return nil
	}

	result := map[string]string{}
	for _, kv := range strings.Split(s[1], ",") {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[strings.Trim(parts[0], "\" ")] = strings.Trim(parts[1], "\" ")
	}
	return result
}

func randomKey() string {
	k := make([]byte, 12)
	for bytes := 0; bytes < len(k); {
		n, err := rand.Read(k[bytes:])
		if err != nil {
			panic("rand.Read() failed")
		}
		bytes += n
	}
	return base64.StdEncoding.EncodeToString(k)
}

/*
H function for MD5 algorithm (returns a lower-case hex MD5 digest)
*/
func h(data string) string {
	digest := md5.New()
	digest.Write([]byte(data))
	return fmt.Sprintf("%x", digest.Sum(nil))
}
