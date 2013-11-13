// Provides a client for interacting with the Bittorrent Sync API.
package btsync_api

import (
  //  "fmt"
  "log"
  "os"
  "strconv"
  "strings"
)

const endpoint = "http://127.0.0.1:%d/api/?"

// Client used to connect to the Bittorrent Sync API.
type BTSyncAPI struct {
  Username string
  Password string
  Port     int
  Endpoint string
  Debug    bool
  Logger   *log.Logger
}

// Instantiates a new BTSyncAPI.
func New(login string, password string, port int, debug bool) *BTSyncAPI {
  logger := log.New(os.Stdout, "[BTSyncAPI] ", log.Ldate|log.Ltime)
  return &BTSyncAPI{login, password, port, endpoint, debug, logger}
}

// Make a request to the API. Requires the method name and a map
// of query string parameters.
func (api *BTSyncAPI) Request(method string, args map[string]string) *Request {
  return &Request{
    API:    api,
    Method: method,
    Args:   args,
  }
}

// Add a folder to Sync with the given secret.
func (api *BTSyncAPI) AddFolderWithSecret(folder string, secret string) (response *Response, err error) {
  args := map[string]string{
    "dir": folder,
  }

  if secret != "" {
    args["secret"] = secret
  }

  request := api.Request("add_folder", args)

  response = &Response{}
  err = request.GetResponse(response)

  return
}

// Add a new folder to Sync.
func (api *BTSyncAPI) AddFolder(folder string) (*Response, error) {
  return api.AddFolderWithSecret(folder, "")
}

// Remove a folder from Sync using its secret.
func (api *BTSyncAPI) RemoveFolder(secret string) (response *Response, err error) {
  request := api.Request("remove_folder", map[string]string{
    "secret": secret,
  })

  response = &Response{}
  err = request.GetResponse(response)

  return
}

// Get information about a folder.
func (api *BTSyncAPI) GetFolder(secret string) (response *GetFoldersResponse, err error) {
  args := map[string]string{}

  if secret != "" {
    args["secret"] = secret
  }

  request := api.Request("get_folders", args)

  response = &GetFoldersResponse{}
  err = request.GetResponse(response)

  return
}

// Get a list of all folders.
func (api *BTSyncAPI) GetFolders() (*GetFoldersResponse, error) {
  return api.GetFolder("")
}

// Get a list of all files for the given secret under the given path.
func (api *BTSyncAPI) GetFilesForPath(secret string, path string) (response *GetFilesResponse, err error) {
  args := map[string]string{
    "secret": secret,
  }

  if path != "" {
    args["path"] = path
  }

  request := api.Request("get_files", args)

  response = &GetFilesResponse{}
  err = request.GetResponse(response)

  return
}

// Get a list of all files for a folder.
func (api *BTSyncAPI) GetFiles(secret string) (*GetFilesResponse, error) {
  return api.GetFilesForPath(secret, "")
}

// Set preferences for a file.
func (api *BTSyncAPI) SetFilePrefs(secret string, path string, download int) (response *SetFilePrefsResponse, err error) {
  request := api.Request("set_file_prefs", map[string]string{
    "secret":   secret,
    "path":     path,
    "download": strconv.Itoa(download),
  })

  response = &SetFilePrefsResponse{}
  err = request.GetResponse(response)

  return
}

// Get a list of all peers for a folder.
func (api *BTSyncAPI) GetFolderPeers(secret string) (response *GetFolderPeersResponse, err error) {
  request := api.Request("get_folder_peers", map[string]string{
    "secret": secret,
  })

  response = &GetFolderPeersResponse{}
  err = request.GetResponse(response)

  return
}

// Get the secrets for a folder.
func (api *BTSyncAPI) GetSecretsForSecret(secret string) (response *GetSecretsResponse, err error) {
  request := api.Request("get_secrets", map[string]string{
    "secret": secret,
  })

  response = &GetSecretsResponse{}
  err = request.GetResponse(response)

  return
}

// Generate new secrets.
func (api *BTSyncAPI) GetSecrets(encryption bool) (response *GetSecretsResponse, err error) {
  args := map[string]string{}

  if encryption {
    args["type"] = "encryption"
  }

  request := api.Request("get_secrets", args)

  response = &GetSecretsResponse{}
  err = request.GetResponse(response)

  return
}

// Get preferences for a folder.
func (api *BTSyncAPI) GetFolderPrefs(secret string) (response *GetFolderPrefsResponse, err error) {
  request := api.Request("get_folder_prefs", map[string]string{
    "secret": secret,
  })

  response = &GetFolderPrefsResponse{}
  err = request.GetResponse(response)

  return
}

// Set preferences for a folder.
func (api *BTSyncAPI) SetFolderPrefs(secret string, prefs *FolderPreferences) (response *SetFolderPrefsResponse, err error) {
  args := structToMap(prefs)
  args["secret"] = secret

  request := api.Request("set_folder_prefs", args)

  response = &SetFolderPrefsResponse{}
  err = request.GetResponse(response)

  return
}

// Get a list of hosts for a folder.
func (api *BTSyncAPI) GetFolderHosts(secret string) (response *GetFolderHostsResponse, err error) {
  request := api.Request("get_folder_hosts", map[string]string{
    "secret": secret,
  })

  response = &GetFolderHostsResponse{}
  err = request.GetResponse(response)

  return
}

// Set the list of hosts for a folder.
func (api *BTSyncAPI) SetFolderHosts(secret string, hosts []string) (response *Response, err error) {
  request := api.Request("set_folder_hosts", map[string]string{
    "secret": secret,
    "hosts":  strings.Join(hosts, ","),
  })

  response = &Response{}
  err = request.GetResponse(response)

  return
}

// Get Sync preferences.
func (api *BTSyncAPI) GetPreferences() (response *GetPreferencesResponse, err error) {
  request := api.Request("get_prefs", map[string]string{})

  response = &GetPreferencesResponse{}
  err = request.GetResponse(response)

  return
}

// Set Sync preferences.
func (api *BTSyncAPI) SetPreferences(prefs Preferences) (response *Response, err error) {
  request := api.Request("set_prefs", map[string]string{})

  prefsMap := structToMap(prefs)

  for key, value := range prefsMap {
    request.Args[key] = string(value)
  }

  response = &Response{}
  err = request.GetResponse(response)

  return
}

// Get name of OS.
func (api *BTSyncAPI) GetOS() (response *GetOSResponse, err error) {
  request := api.Request("get_os", map[string]string{})

  response = &GetOSResponse{}
  err = request.GetResponse(response)

  return
}

// Get Bittorrent Sync version.
func (api *BTSyncAPI) GetVersion() (response *GetVersionResponse, err error) {
  request := api.Request("get_version", map[string]string{})

  response = &GetVersionResponse{}
  err = request.GetResponse(response)

  return
}

// Get current upload and download speed.
func (api *BTSyncAPI) GetSpeed() (response *GetSpeedResponse, err error) {
  request := api.Request("get_speed", map[string]string{})

  response = &GetSpeedResponse{}
  err = request.GetResponse(response)

  return
}
