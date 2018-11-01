package main

import (
	"bytes"
	"fmt"
	"log"
	"mypackages"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo"
)

func GetHostInfo(hostID string) (string, error) {
	// Поиск IP адреса хоста по ID и выполнение команды по ssh
	hosts := mypackages.SetOneClient()
	sumresult := ""
	for i := range hosts.Hosts {
	savePoint:
		if hosts.Hosts[i].ID == hostID {
			result, err := mypackages.SSHExec(changeIp(hosts.Hosts[i].IPAddress), "lshw -json")
			if err != nil {
				return " ", err
			}

			return result, err
		}
		if hostID == "all" {
			result, err := mypackages.SSHExec(changeIp(hosts.Hosts[i].IPAddress), "lshw -json")
			sumresult = sumresult + result
			if err != nil {
				i++
				goto savePoint
			}
		}
	}
	return sumresult, nil
}

func changeIp(s string) string {
	str := strings.Replace(s, "172.22.22.", "192.168.2.", -1)
	return str
}

//Method for starting EtcdCluster on 3 hosts
func startEtcdCluster(c echo.Context) error {
	id := c.Param("id")
	var tmp string
	switch id {
	case "1":
		tmp = "192.168.2.229"
	case "2":
		tmp = "192.168.2.176"
	case "3":
		tmp = "192.168.2.171"
	}
	cmd := exec.Command("ssh", "root@"+tmp, "source", "/home/egor/makeEtcdCluster")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return c.String(http.StatusOK, out.String())
}

//Usual home page (not important)
func homePage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the home page")
}

func handleRequest() {
	e := echo.New()
	e.GET("/", homePage)
	//Function for post
	e.POST("/ansible/params", mypackages.AddEtcdParams)

	e.GET("/name/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.String(http.StatusOK, id)
	})

	e.GET("/ansible/:command", mypackages.AddEtcdParams2)

	e.GET("/etcd/cluster:id", startEtcdCluster)

	e.GET("/hosts", func(c echo.Context) error {
		return c.JSON(http.StatusOK, mypackages.SetOneClient())
	})

	e.GET("/hosts/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(http.StatusOK, mypackages.GetSoloResponse(id))
	})

	e.GET("/hosts/:id/info", func(c echo.Context) error {
		id := c.Param("id")
		sumresult, err := GetHostInfo(id)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, sumresult)
	})

	e.Logger.Fatal(e.Start(":8181"))
}

func main() {
	fmt.Println("Starting server:...")
	handleRequest()
}
