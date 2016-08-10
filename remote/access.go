package remote

import (
  "fmt"
  "sync"
  "time"
  "errors"
  "encoding/csv"
)

/// --- Private Types
type callback func () []string

/// -- Public Data Types
type Access struct {
  Clients         map[string]*RemoteClient
  ConnectedEnum   []string
  Writer          *csv.Writer
  ConsoleEnabled  bool
  PrivateKeys     []string
  TimeStamp       time.Time
  Period          time.Duration
  Sync            sync.WaitGroup
  cb              callback
}

/// -- Private Functions
// -----------------------------------------------------------------------------
func (a * Access) login (name string) (error) {
  if a != nil {
    if client := a.Clients[name]; client != nil {
      if err := client.Connect (); err != nil {
        return err
      }
    }
  } else {
    return errors.New ("Invalid object. Cannot login")
  }
  return nil
}

// -----------------------------------------------------------------------------
func (a * Access) tick () {
  for {
    time.Sleep (a.Period)
    now := time.Now ()

    a.TimeStamp = now
    for _, client := range (a.Clients) {
      client.timeSync <- now
    }
  }
}

// -----------------------------------------------------------------------------
func (a * Access) spawn (client *RemoteClient) {
  // Let's run this forever, or at least until the program ends
  for {
    client.newSession ()
    select {
    case timeStamp := <-client.timeSync:
      client.results = client.GetResources ()
      if a.ConsoleEnabled {
        fmt.Printf ("Time: %v, Client: %s, Resources: %v\n", timeStamp, client.host, client.results)
      }
      a.Sync.Done ()
    }
    client.session.Close ()
  }
}

// -----------------------------------------------------------------------------
func (a * Access) csvSync () {
  // Loop forever: Will wait until spawn routines are synced with resource values
  // at which point they will get written to the .csv file.
  for {
    a.Sync.Add (len (a.Clients))
    a.Sync.Wait ()

    // Write to csv file
    var record []string
    record = append (record, a.TimeStamp.Format ("01/02/2006 15:04:05"))

    if a.cb != nil {
      slice := a.cb ()
      for _, elems := range (slice) {
        record = append (record, elems)
      }
    }

    for _, element := range (a.ConnectedEnum) {
      fmt.Printf ("Client: %s\n", element)
      client := a.Clients[element]
      for _, val := range (client.results) {
        record = append (record, val)
      }
    }

    a.Writer.Write (record)
    a.Writer.Flush ()
  }
}

/// -- Public Functions
// -----------------------------------------------------------------------------
func NewRemoteAccess (writer *csv.Writer, cb callback, period time.Duration, enableConsole bool, privateKeys []string) (*Access) {
  newAccess := new (Access)
  if newAccess == nil {
    return nil
  } else {
    newAccess.Clients = make (map[string]*RemoteClient)
    newAccess.Writer = writer
    newAccess.ConsoleEnabled = enableConsole
    newAccess.PrivateKeys = privateKeys
    newAccess.Period = period
    newAccess.cb = cb
  }
  return newAccess
}

// -----------------------------------------------------------------------------
func (a * Access) AddClient (id, username, host, port string) (error) {
  if a != nil {
    // Create new client
    client := newClient (username, host, port, a.PrivateKeys, 5 * time.Second)
    if client == nil {
      return errors.New ("Invalid object. Cannot create new client")
    }
    a.Clients[id] = client
  } else {
    return errors.New ("Invalid object. Cannot add new client!")
  }
  return nil
}

// -----------------------------------------------------------------------------
func (a * Access) Start () {
 if a != nil && a.Writer != nil {
   for key, client := range (a.Clients) {
     fmt.Printf ("Logging into client: %s...", key)
     err := a.login (key)
     if err != nil {
       fmt.Printf ("Ouch! Cannot start client %s. Err = %v\n", key, err)
     } else {
       fmt.Println ("OK")
       a.ConnectedEnum = append (a.ConnectedEnum, key)
       go a.spawn (client)
     }
   }
   go a.tick ()

   a.csvSync () // Stay here

 } else {
   panic ("Access object is invalid!")
 }
}
