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
  MACHINE = "ditto.sensity.com" // IP addresses are allowed too
  PRIVATE_KEY = "/home/bwinata/.ssh/id_rsa.key"
)

func testCallback () ([]string) {
  var record []string
  record = append (record, "10")
  return record
}


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
  access := remote.NewRemoteAccess (writer, testCallback, 5 * time.Second, true, keys)
  if access == nil {
    fmt.Println ("Error: Cannot create access object!")
    return
  }
  fmt.Println ("OK")

  // Add clients:
  fmt.Printf ("Add new client: %s...", MACHINE)
  err = access.AddClient ("ditto", "ubuntu", MACHINE, "22")
  if err != nil {
    fmt.Println ("Error: Cannot add client: %s", MACHINE)
  }
  fmt.Println ("OK")

  // Connect and create session for each client
  access.Start ()

  // Wait here forever
  forever := make (chan bool)
  <- forever

}
