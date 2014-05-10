package main

import (
    "os"
    "fmt"
    "net/http"
    "crypto/tls"
    "io/ioutil"
    "strings"
    "crypto/x509"
)

func main() {
  certs := x509.NewCertPool()

  pemData, err := ioutil.ReadFile("myCA.cer")
  if err != nil {
    fmt.Println(err)  // do error
  }
  certs.AppendCertsFromPEM(pemData)
    tr := &http.Transport{
      TLSClientConfig: &tls.Config{ RootCAs: certs }, //,InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    body := strings.NewReader(os.Args[1])
    req, err := http.NewRequest("POST", "https://pepe.imim.es:20000/run", body)
    if err != nil {
        fmt.Println(err)
        return
    }
    req.Header.Add("X-AUTH-TOKEN", "blabla")
    resp, err := client.Do(req)

    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

    body2, err := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body2))

}
