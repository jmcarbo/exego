package main

import (
    "os"
    "fmt"
    "net/http"
    "crypto/tls"
    "io/ioutil"
    "strings"
    "crypto/x509"
    "path"
    "github.com/jmcarbo/exego"
    "github.com/spf13/viper"
)

func main() {
  viper.SetConfigName("config") // name of config file (without extension)
  viper.AddConfigPath("/etc/exego/")   // path to look for the config file in
  viper.AddConfigPath("$HOME/.exego")  // call multiple times to add many search paths
  viper.AddConfigPath("./")
  viper.ReadInConfig() // Find and read the config file

  viper.SetDefault("CAcert", "myCA.cer")
  viper.SetDefault("AuthToken", "blabla")
  viper.SetDefault("Addr", "localhost:20000")

  certs := x509.NewCertPool()
  pemData, err := ioutil.ReadFile(viper.GetString("CAcert"))
  if err != nil {
    pemData, err = exego.Asset(path.Join("certs", viper.GetString("CAcert")))
  }
  certs.AppendCertsFromPEM(pemData)
  tr := &http.Transport{
    TLSClientConfig: &tls.Config{ RootCAs: certs }, //,InsecureSkipVerify: true},
  }
  client := &http.Client{Transport: tr}
  body := strings.NewReader(os.Args[len(os.Args)-1])
  req, err := http.NewRequest("POST", "https://"+viper.GetString("Addr")+"/run", body)

  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Add("X-AUTH-TOKEN", viper.GetString("AuthToken"))
  resp, err := client.Do(req)

  if err != nil {
    fmt.Println(err)
    return
  }
  defer resp.Body.Close()

  body2, err := ioutil.ReadAll(resp.Body)
  fmt.Println(string(body2))

}
