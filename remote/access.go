package remote

import (
  "fmt"
  "time"
  "errors"
  "encoding/csv"
)

/// -- Public Data Types
type Access struct {
  Clients         map[string]*RemoteClient
  Writer          *csv.Writer
  ConsoleEnabled  bool
  PrivateKeys     []string
  Period          time.Duration
  TimeSync        chan time.Time
}

/// -- Private Functions
// -----------------------------------------------------------------------------
func (a * Access) login (name string) (error) {
  if a != nil {
    if client := a.Clients[name]; client != nil {
      if err := client.Connect (); err != nil {
        return err
      }
      if err := client.newSession (); err != nil {
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
    a.TimeSync <- time.Now ()
  }
}

// -----------------------------------------------------------------------------
func (a * Access) spawn (client *RemoteClient) {
  // Let's run this forever, or at least until the program ends
  for {
    select {
    case timeStamp := <-a.TimeSync:
      fmt.Printf ("Time: %v, Client: %s, Resources:%v\n", timeStamp, client.host, client.GetResources ())
    }
  }
}

/// -- Public Functions
// -----------------------------------------------------------------------------
func NewRemoteAccess (writer *csv.Writer, period time.Duration, enableConsole bool, privateKeys []string) (*Access) {
  newAccess := new (Access)
  if newAccess == nil {
    return nil
  } else {
    newAccess.Clients = make (map[string]*RemoteClient)
    newAccess.Writer = writer
    newAccess.ConsoleEnabled = enableConsole
    newAccess.PrivateKeys = privateKeys
    newAccess.Period = period
    newAccess.TimeSync = make (chan time.Time)
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
       go a.spawn (client)
     }
   }
   fmt.Printf ("Spawning clients...\n")
   go a.tick ()
 } else {
   panic ("Access object is invalid!")
 }
}
