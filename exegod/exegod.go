package main

import (
  "fmt"
  "os"
  "os/exec"
  "syscall"
  "time"
  "bytes"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "log"
  "github.com/jmcarbo/exego"
  "crypto/tls"
  "github.com/spf13/viper"
)

const (
  EXEGOD_VERSION = "0.2.0"

  ERR_UNAUTHORIZED = 401
)

var authToken string

type Command struct {
  Command    string         `json:"command"`
  ExitStatus int            `json:"exit_status"`
  Output     string         `json:"output"`
  TimeStart  time.Time      `json:"time_start"`
  TimeFinish time.Time      `json:"time_finish"`
  Duration   time.Duration  `json:"duration"`
}

func HandleShellRequest(w http.ResponseWriter, req *http.Request) {
  token := req.Header.Get("X-AUTH-TOKEN")

  if token != authToken {
    w.WriteHeader(ERR_UNAUTHORIZED)
    w.Write([]byte("Invalid authentication token\r\n"))
    return
  }

  body, err := ioutil.ReadAll(req.Body)

  if err != nil {
    fmt.Println("Read error: ", err)
     w.WriteHeader(400)
    w.Write([]byte("Read error"))
    return
  }

  cmd := Exec(string(body))
  cmd.Print()

  w.Header().Set("X-EXEGOD-DURATION", cmd.Duration.String())
  w.Header().Set("X-EXEGOD-STATUS", string(cmd.ExitStatus))

  w.WriteHeader(200)
  w.Write([]byte(cmd.ToJson()))
  w.(http.Flusher).Flush()
}

func (cmd *Command) Run(command string) {
  var output bytes.Buffer
  shell := exec.Command("bash", "-c", command)

  shell.Stdout  = &output
  shell.Stderr  = &output
  cmd.Command   = command
  cmd.TimeStart = time.Now()

  err := shell.Start()
  if (err != nil) {
    fmt.Println("Execution failed:", err)
  }

  err = shell.Wait()

  cmd.TimeFinish = time.Now()
  cmd.Duration   = cmd.TimeFinish.Sub(cmd.TimeStart)

  if msg, ok := err.(*exec.ExitError); ok {
    cmd.ExitStatus = msg.Sys().(syscall.WaitStatus).ExitStatus()
  } else {
    cmd.ExitStatus = 0
  }

  cmd.Output = string(output.Bytes())
}

func (cmd *Command) ToJson() (s string) { 
  buff, err := json.Marshal(cmd)

  if (err != nil) {
    s = ""
    return
  }

  return string(buff)
}

func (cmd *Command) Print() {
  fmt.Println("Command:   ", cmd.Command)
  fmt.Println("ExitStatus:", cmd.ExitStatus)
  fmt.Println("Duration:  ", cmd.Duration)
  fmt.Println("Output:    ", cmd.Output)
}

func (cmd *Command) Success() bool {
  return cmd.ExitStatus == 0
}

func Exec(str string) *Command {
  var command *Command
  command = new(Command)
  command.Run(str)

  return command
}

func EnvVarDefined(name string) bool {
  result := os.Getenv(name)
  return len(result) > 0
}

func main() {
  viper.SetConfigName("config") // name of config file (without extension)
  viper.AddConfigPath("/etc/exego/")   // path to look for the config file in
  viper.AddConfigPath("$HOME/.exego")  // call multiple times to add many search paths
  viper.AddConfigPath("./")
  viper.ReadInConfig() // Find and read the config file

  viper.SetDefault("CAcert", "myCA.cer")
  viper.SetDefault("AuthToken", "blabla")
  viper.SetDefault("Addr", "localhost:20000")

  authToken = viper.GetString("AuthToken")
  bindAddr := viper.GetString("Addr")


  log.Printf("Exegod v%s\n", EXEGOD_VERSION)
  fmt.Printf("Starting server on %s\n", bindAddr)

  http.HandleFunc("/run", HandleShellRequest)

  //_, cerr := os.Open()
  //_, kerr := os.Open(exego.A"mycert1.key")

  //if os.IsNotExist(cerr) || os.IsNotExist(kerr) {
  //  log.Fatalln(cerr, kerr)
  //  return
  //}

  //http.ListenAndServeTLS(bindAddr, exego.Asset("certs/mycert1.cer"), exego.Asset("certs/mycert1.key"), nil)
  cer,_ := exego.Asset("certs/mycert1.cer")
  key,_ := exego.Asset("certs/mycert1.key")
  cert, _ := tls.X509KeyPair(cer, key)
  config := tls.Config{Certificates: []tls.Certificate{cert}}
  listener, _ := tls.Listen("tcp", bindAddr, &config)
  server := &http.Server{Addr: bindAddr, Handler: nil}
  server.Serve(listener)
}
