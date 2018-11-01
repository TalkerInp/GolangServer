package mypackages

import (
	"bytes"
	"crypto/tls"
	"net/http"

	"github.com/kolo/xmlrpc"
	"gopkg.in/xmlpath.v2"
)

type HostCollection struct {
	Hosts []Host `json:"Hosts"`
}
type Host struct {
	ID        string `json:"ID"`
	IPAddress string `json:"IPAddress"`
}
type Response struct {
	status  bool
	body    string
	bodyInt int
}

func SetOneClient() HostCollection { // xmlrpc api
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //https
	}

	xmlrpcClient, err := xmlrpc.NewClient("https://192.168.2.108:2634/RPC2", tr)
	if err != nil {
		xmlrpcClient, _ = xmlrpc.NewClient("http://192.168.2.108:2633/RPC2", tr)
	}

	USER := "cms_test" //login
	PASS := "cms_test" //password
	args := USER + ":" + PASS
	xmlArgs := make([]interface{}, len(args)+1)
	xmlArgs[0] = args
	result := []interface{}{}
	xmlrpcClient.Call("one.hostpool.info", xmlArgs[0], &result)

	var response Response
	response.body = result[1].(string)

	path := xmlpath.MustCompile("/HOST_POOL/HOST")
	root, _ := xmlpath.Parse(bytes.NewReader([]byte(response.body)))

	hostpollRoots := path.Iter(root)
	parseResult := HostCollection{}

	for hostpollRoots.Next() {
		hostpollRoot := hostpollRoots.Node()
		host := Host{}
		path = xmlpath.MustCompile("ID")
		if value, ok := path.String(hostpollRoot); ok {
			host.ID = value
		}
		path = xmlpath.MustCompile("NAME")
		if value, ok := path.String(hostpollRoot); ok {
			host.IPAddress = value
		}
		parseResult.Hosts = append(parseResult.Hosts, host)
	}
	return parseResult
}

func GetSoloResponse(hostID string) Host {
	hosts := SetOneClient()
	for i := range hosts.Hosts {
		if hosts.Hosts[i].ID == hostID {
			return hosts.Hosts[i]
		}
	}
	return hosts.Hosts[1]
}
