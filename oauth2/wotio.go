package oauth2

import (
	"net/http"
	"crypto/tls"
	"crypto/x509"
	"os"
	"bytes"
	"fmt"
	"io/ioutil"

	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"time"
)

func GetWotioRootCA() ([]byte) {
	var cert []byte
	certPath := os.Getenv("WOTIO_CA_CERT_PATH")
	if certPath != "" {
		var err error
		cert, err = ioutil.ReadFile(certPath)
		if err != nil {
			logrus.Warnln("Could not read wotio certificate: ", err)
		}
	}
	return cert
}

func GetWotioCertPool(cert []byte) (*x509.CertPool) {
	pool := x509.NewCertPool()
	if cert != nil {
		pool.AppendCertsFromPEM(cert)
	}
	return pool
}


func WotioCreateToken(token string, expires_in int64) (error) {
	WotioToken := os.Getenv("WOTIO_TOKEN")
	WotioTokenUrl := os.Getenv("WOTIO_TOKEN_URL")
	if WotioTokenUrl == "" {
		return errors.New("WOTIO_TOKEN_URL is not set.")
	}
	start:= time.Now()
	end := start.Add(time.Duration(expires_in)*time.Second)
	json := fmt.Sprintf("{ \"token\": \"%s\", \"start\": \"%s\", \"end\": \"%s\", \"type\":\"hydra-oauth2-v1\" }", token, start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339))
	logrus.Debugln("WotioCreateToken:",json)

	req, err := http.NewRequest("POST", WotioTokenUrl, bytes.NewBufferString(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + WotioToken)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: GetWotioCertPool(GetWotioRootCA())}}}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Debugln("WotioCreateToken response:", resp.Status, string(body))
	return nil
}

func WotioSetTokenScopes(token string, scopes []string) (error) {
	WotioToken := os.Getenv("WOTIO_TOKEN")
	WotioTokenUrl := os.Getenv("WOTIO_TOKEN_URL")
	if WotioTokenUrl == "" {
		return errors.New("WOTIO_TOKEN_URL is not set.")
	}
	
	j, _ := json.Marshal(scopes)
	logrus.Debugln("WotioSetTokenScopes:",string(j))

	url := fmt.Sprintf("%s/%s/scopes", WotioTokenUrl, token)
	logrus.Debugln("WotioSetTokenScopes:",url)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(string(j)))
	req.Header.Set("Authorization", "Bearer " + WotioToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: GetWotioCertPool(GetWotioRootCA())}}}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Debugln("WotioSetTokenScopes response:", resp.Status, string(body))
	return nil
}
