package mypackages

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/labstack/echo"
	"go.etcd.io/etcd/client"
)

func SetEtcdKeyValue(key string, value string) {
	cfg := client.Config{
		Endpoints: []string{"http://192.168.2.229:2379", "http://192.168.2.176:2379", "http://192.168.2.171:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)
	// set key and value
	log.Print("Setting " + key + " key with " + value + " value")
	resp, err := kapi.Set(context.Background(), key, value, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}
	// get key's value
	log.Print("Getting " + key + " key value")
	resp, err = kapi.Get(context.Background(), key, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
}

// struct for method "addEtcdParams"
type EtcdParams struct {
	Vardevice string `json:"vardevice"`
	Varnumber string `json:"varnumber"`
	Varlable  string `json:"varlable"`
	Varstate  string `json:"varstate"`
}

//POST method for adding etcd key and value
func AddEtcdParams(c echo.Context) error {
	etcdParams := EtcdParams{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	var s bytes.Buffer
	s.Write(b)
	fmt.Printf("%s\n", s.String())
	fmt.Println("Hello")
	if err != nil {
		log.Printf("failed reading the request body for addEtcdParams: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &etcdParams)
	if err != nil {
		log.Println("Failed unmarshaling in addEtcdParams  : ", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is our params: %#v", etcdParams)
	return c.String(http.StatusOK, "Done")
}

// GET method for adding etcd key and value
func AddEtcdParams2(c echo.Context) error {
	command := c.Param("command")
	vardevice := c.QueryParam("vardevice")
	varlable := c.QueryParam("varlable")
	varnumber := c.QueryParam("varnumber")
	varstate := c.QueryParam("varstate")
	varpart_start := c.QueryParam("varpart_start")
	varpart_end := c.QueryParam("varpart_end")
	fmt.Println("varpartEnd = ", varpart_end)

	SetEtcdKeyValue("vardevice", vardevice)
	SetEtcdKeyValue("varlable", varlable)
	SetEtcdKeyValue("varnumber", varnumber)
	SetEtcdKeyValue("varstate", varstate)
	SetEtcdKeyValue("varpart_start", varpart_start)
	SetEtcdKeyValue("varpart_end", varpart_end)

	cmd := exec.Command("ansible-playbook", "/home/egor/golang/projects/mainproject/src/ansible/"+command+".yml")
	var b bytes.Buffer
	cmd.Stdout = &b
	// cmd.Stdin = os.Stdin
	// cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// log.Fatal("Ошибка тут:"+err.Error())
		return c.String(http.StatusOK, "error with one of value in '"+command+".yml'")
	}
	return c.String(http.StatusOK, b.String())

}
