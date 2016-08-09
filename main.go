package main

import (
  "os"
  "fmt"
  "time"
  "encoding/csv"

  "github.com/bwinata/stockpile/remote"
)

// Temporary. Will soon get from config file when complete.
const (
  CLIENT_DITTO = "ditto.sensity.com"
  CLIENT_DCC_DEV = "dcc-gina-dev.sensity.com"
  PRIVATE_KEY = "/home/bwinata/.ssh/id_rsa.key"
)

/// -- Entry Point
func main () {
  // Create new .csv file
  fmt.Print ("Creating new csv file...")
  file, err := os.Create ("./test.csv")
  if err != nil {
    fmt.Println ("Error: %v", err)
    return
  }
  defer file.Close ()
  fmt.Println ("OK")

  // Create new CSV writer
  fmt.Print ("Creating new writer...")
  writer := csv.NewWriter (file)
  fmt.Println ("OK")

  keys := []string {
    PRIVATE_KEY,
  }

  // Create new remote access object
  fmt.Print ("Creating new remote access object...")
  access := remote.NewRemoteAccess (writer, 5 * time.Second, true, keys)
  if access == nil {
    fmt.Println ("Error: Cannot create access object!")
    return
  }
  fmt.Println ("OK")

  // Add client: Ditto
  fmt.Printf ("Add new client: %s...", CLIENT_DITTO)
  err = access.AddClient ("ditto", "ubuntu", CLIENT_DITTO, "22")
  if err != nil {
    fmt.Println ("Error: Cannot add client: %s", CLIENT_DITTO)
  }
  fmt.Println ("OK")

  fmt.Printf ("Add new client: %s...", CLIENT_DCC_DEV)
  err = access.AddClient ("dcc", "ubuntu", CLIENT_DCC_DEV, "22")
  if err != nil {
    fmt.Println ("Error: Cannot add client: %s", CLIENT_DCC_DEV)
  }
  fmt.Println ("OK")

  // Connect and create session for each client
  access.Start ()


  // Wait here forever
  forever := make (chan bool)
  <- forever

}
