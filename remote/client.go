package remote

import (
  "io"
  "os"
  "fmt"
  "time"
  "bytes"
  "errors"
  "strings"
  "io/ioutil"

  "golang.org/x/crypto/ssh"

  "github.com/bwinata/stockpile/resources"
)

/// -- Private Data Types
type RemoteClient struct {
  username   string
  host       string
  protoCfg   *ssh.ClientConfig
  handler    *ssh.Client
  session    *ssh.Session
  Output     chan string
  stdout     io.Reader
  stdin      io.Writer
  timeSync   chan time.Time
  results    []string
}

/// -- Private
// -----------------------------------------------------------------------------
func readKey (keyPath string) (ssh.Signer, error) {
  fileHandle, err := os.Open (keyPath)
  if err != nil {
    return nil, err
  }

  defer fileHandle.Close ()

  buf, err := ioutil.ReadAll (fileHandle)
  if err != nil {
    return nil, err
  }

  signer, err := ssh.ParsePrivateKey (buf)
  if err != nil {
    return nil, err
  }

  return signer, nil
}

// -----------------------------------------------------------------------------
func makeKeyrings (keys []string) ([]ssh.AuthMethod) {
  auths := []ssh.AuthMethod{}

  for _, key := range (keys) {
    signer, err := readKey (key)
    if err == nil {
      auths = append (auths, ssh.PublicKeys (signer))
    } else {
      fmt.Printf ("Cannot create keyring for key: %s. Err=%v\n", key, err)
      return nil
    }
  }
  return auths
}

// -----------------------------------------------------------------------------
func (client * RemoteClient) listenOutput () {
    var buf bytes.Buffer
    io.Copy (&buf, client.stdout)
    client.Output <- buf.String ()
}

// -----------------------------------------------------------------------------
func (client * RemoteClient) Connect () (error) {
  if client != nil {
    connection, err := ssh.Dial ("tcp", client.host, client.protoCfg)
    if err != nil {
      return err
    }
    client.handler = connection
  } else {
    return errors.New ("Invalid object")
  }
  return nil
}

// -----------------------------------------------------------------------------
func (client * RemoteClient) newSession () (error) {
  if client != nil {
    session, err := client.handler.NewSession ()
    if err != nil {
      return err
    }
    client.session = session

    stdout, err := client.session.StdoutPipe ()
    if err != nil {
      return err
    }
    client.stdout = stdout

    go client.listenOutput ()

  } else {
    return errors.New ("Invalid object")
  }

  return nil
}

// -----------------------------------------------------------------------------
func (client * RemoteClient) GetResources () ([]string) {
  if client != nil && client.Output != nil  {
    if err := client.Exec (fmt.Sprintf ("%s;%s;%s",
                                        resources.CPU_USAGE_PERCENTAGE,
                                        resources.MEMORY_USAGE_MB,
                                        resources.MEMORY_FREE_MB)); err != nil {
      fmt.Printf ("Cannot run command. Err=%v\n", err)
    }

    // Wait for response
    select {
    case response := <- client.Output:
      return strings.Split (response, "\n")
    }
  } else {
    panic (errors.New ("Invalid object"))
  }
}

// -----------------------------------------------------------------------------
func (client * RemoteClient) Exec (command string) (error) {
  if client != nil {
    if err := client.session.Run (command); err != nil {
      return err
    }
  } else {
    return errors.New ("Invalid object")
  }

  return nil
}

// -----------------------------------------------------------------------------
func newClient (username, hostname, port string, keys []string, timeout time.Duration) (*RemoteClient) {
  newClient := new (RemoteClient)
  if newClient == nil {
    return nil
  }

  // Attempt to create keyrings first
  auths := makeKeyrings (keys)

  newClient.username = username
  newClient.host = fmt.Sprintf ("%s:%s", hostname, port)
  newClient.Output = make (chan string)
  newClient.timeSync = make (chan time.Time)
  newClient.protoCfg = &ssh.ClientConfig {
    User    : username,
    Auth    : auths,
    Timeout : timeout * time.Second,
  }

  return newClient
}
