package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "text/template"
    "github.com/docopt/docopt-go"
)

func main() {
  usage := `
      
Usage:
  create-basic-configdrive -H HOSTNAME -S SSH_FILE [-p PATH] [-t TOKEN | -d URL] [-n NAME] [-e URL] [-i URL] [-l URLS] [-u URL] [-h] [-v]

Options:
  -H HOSTNAME  Machine hostname.
  -S FILE      SSH keys file.
  -p DEST      Create config-drive ISO image to the given path.
  -t TOKEN     Token ID from https://discovery.etcd.io.
               
  -d URL       Full URL path to discovery endpoint.
               [default: https//discovery.etcd.io/TOKEN]
                
  -n NAME      etcd node name (defaults to HOSTNAME).
               
  -e URL       Advertise URL for client communication.
               [default: http://\$public_ipv4:2379]
               
  -i URL       URL for server communication.
               [default: http://\$private_ipv4:2380]
               
  -l URLS      Listen URLS for client communication.
               [default: http://0.0.0.0:2379,http://0.0.0.0:4001]
               
  -u URL       Listen URL for server communication.
               [default: http://0.0.0.0:2380]
               
  -v           Show version.
  -h           This help.
`

  arguments, _ := docopt.Parse(usage, nil, true, "Coreos go tools 0.1", false)

  tmpl_text := `#cloud-config
coreos:
  etcd2:
    name: {{ .ETCD_NAME }}
    advertise-client-urls: {{ .ETCD_ADDR>
    initial-advertise-peer-urls: {{ .ETCD_PEER_URLS>
    discovery: <ETCD_DISCOVERY>
    listen-peer-urls: <ETCD_LISTEN_PEER_URLS>
    listen-client-urls: <ETCD_LISTEN_CLIENT_URLS>
  units:
    - name: etcd2.service
      command: start
    - name: fleet.service
      command: start
ssh_authorized_keys:
  - <SSH_KEY>
hostname: <HOSTNAME>
`
  var tmpl_map map[string]string
  
  tmpl_map["ETCD_NAME"], ok := arguments["-n"].(string) 
  if ok == false {
    tmpl_map["ETCD_NAME"], _ := arguments["-H"].(string) 
  }
  tmpl, _ := template.New("test").Parse(tmpl_text)
  _ = tmpl.Execute(os.Stdout, tmpl_map)

  
      dest, ok := arguments["-p"].(string)
      if ok == false {
        dest, _ = os.Getwd()
      }
      
      fmt.Println(arguments)
      workdir, _ := ioutil.TempDir(dest, "coreos")
      _ = os.MkdirAll(workdir + "/openstack/latest", 0777)
      f, _ := os.Create(workdir + "/openstack/latest/user_data")
      defer f.Close()
      // defer RemoveAll(workdir)

}