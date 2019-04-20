package utils

import (
	"github.com/prometheus/log"
	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
	"io/ioutil"
	"net/http"
)

func GetData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func GetKerberizedData(url string, krbConfig string, principal string, keytabPath string, realm string, spn string) ([]byte, error) {
	cfg, err := config.Load(krbConfig)
	if err != nil {
		log.Error(err)
	}
	ktFromFile, err := keytab.Load(keytabPath)
	if err != nil {
		log.Fatalf("could not load keytab: %v", err)
	}

	cl := client.NewClientWithKeytab(principal, realm, ktFromFile, cfg)
	err = cl.Login()
	if err != nil {
		log.Fatalf("could not login client: %v", err)
	}

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("could create request: %v", err)
	}

	spnegoCl := spnego.NewClient(cl, nil, spn)

	// Make the request
	resp, err := spnegoCl.Do(r)
	if err != nil {
		log.Fatalf("error making request: %v", err)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
