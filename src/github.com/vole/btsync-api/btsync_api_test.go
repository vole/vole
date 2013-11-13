package btsync_api

import (
  "flag"
  "fmt"
  ioutil "io/ioutil"
  "log"
  "os"
  "path"
  "strings"
  "testing"
  "time"
)

// Where temporary test dirs are created.
var Dir = "/tmp"

// Reference to the temp dir used for testing.
var TmpDir = ""

// Reference to a temp file in the temp dir.
var TmpFile *os.File

// Prefix used by the temp dir for easier cleanup.
var Prefix = "btsync_api_test_"

// BTSync API settings.
var login = flag.String("login", "test", "BT Sync API login")
var password = flag.String("password", "test", "BT Sync API password")
var port = flag.Int("port", 8080, "BT Sync API port")

// If tests are failing and you're not sure why, this may help.
var verbose = flag.Bool("verbose", false, "Enable verbose test logging")

// For logging test information and debug stuff.
var logger = log.New(os.Stdout, "[BTSyncAPI] ", log.Ldate|log.Ltime)

// Log a debug message to stdout.
func Debug(msg string, a ...interface{}) {
  if *verbose {
    formatted := fmt.Sprintf(msg, a...)
    logger.Println(fmt.Sprintf("\033[35;1mDEBUG:\033[0m %s", formatted))
  }
}

// Log an info message to stdout.
func Log(msg string, a ...interface{}) {
  formatted := fmt.Sprintf(msg, a...)
  logger.Println(fmt.Sprintf("\033[34;1mINFO:\033[0m %s", formatted))
}

// Create a temp dir to use for tests.
func TestSetup(t *testing.T) {
  Log("Setting up")

  Debug("login: %s", *login)
  Debug("password: %s", *password)
  Debug("port: %d", *port)
  Debug("verbose: %t", *verbose)

  dir, err := ioutil.TempDir(Dir, Prefix)
  if err != nil {
    t.Errorf("Unable to create test directory in %s: %s", Dir, err)
    return
  }

  TmpDir = dir
  Debug("Temp Dir: %s", TmpDir)

  file, err := ioutil.TempFile(TmpDir, Prefix)
  if err != nil {
    t.Errorf("Unable to create temp file in %s: %s", TmpDir, err)
    return
  }

  TmpFile = file
  Debug("Temp File: %s", TmpFile.Name())
}

// Test creating, removing folders.
func TestFolders(t *testing.T) {
  api := New(*login, *password, *port, *verbose)

  Log("Testing AddFolder")

  addFolderResponse, err := api.AddFolder(TmpDir)

  if err != nil {
    t.Errorf("Error making request to add new folder: %s", err)
    return
  }

  if addFolderResponse.Error != 0 {
    t.Errorf("Error adding new folder")
    return
  }

  Log("Testing GetFolders")

  getFoldersResponse, err := api.GetFolders()

  if err != nil {
    t.Errorf("Error making request to get all folders: %s", err)
    return
  }

  if len(*getFoldersResponse) == 0 {
    t.Errorf("Expected at least 1 folder")
    return
  }

  found := false
  var testDir Folder
  for _, v := range *getFoldersResponse {
    if v.Dir == TmpDir {
      found = true
      testDir = v
    }
  }

  if !found {
    t.Errorf("Expected to find %s", TmpDir)
    return
  }

  Log("Testing GetFolder")

  getFolderResponse, err := api.GetFolder(testDir.Secret)
  if err != nil {
    t.Errorf("Error making request to get a single folder: %s", err)
    return
  }

  if len(*getFolderResponse) != 1 {
    t.Errorf("Expected a single folder in response")
    return
  }

  if (*getFolderResponse)[0].Dir != TmpDir {
    t.Errorf("Expected %s to be %s", (*getFolderResponse)[0].Dir, TmpDir)
    return
  }

  Debug("Sleeping for 15 seconds to allow BTSync to pick up new file.")

  time.Sleep(15000 * time.Millisecond)

  Log("Testing GetFiles")

  getFilesResponse, err := api.GetFiles(testDir.Secret)
  if err != nil {
    t.Errorf("Error making request to get files: %s", err)
    return
  }

  if len(*getFilesResponse) != 1 {
    t.Errorf("Expected 1 file")
    return
  }

  if (*getFilesResponse)[0].Name != path.Base(TmpFile.Name()) {
    t.Errorf("Expected %s to be %s", (*getFilesResponse)[0].Name, path.Base(TmpFile.Name()))
    return
  }

  Log("Testing SetFilePrefs")

  setFilePrefsResponse, err := api.SetFilePrefs(testDir.Secret, path.Base(TmpFile.Name()), 1)
  if err != nil {
    t.Errorf("Error making request to set file preferences: %s", err)
    return
  }

  if len(*setFilePrefsResponse) != 1 {
    t.Errorf("Expected response to have a length of 1")
  }

  if (*setFilePrefsResponse)[0].State != "created" {
    t.Errorf("Expected file object in response")
    return
  }

  Log("Testing GetFolderPeers")
  _, err = api.GetFolderPeers(testDir.Secret)
  if err != nil {
    t.Errorf("Error making request to get folder peers: %s", err)
    return
  }

  Log("Testing GetSecrets")

  getSecretsResponse, err := api.GetSecrets(true)
  if err != nil {
    t.Errorf("Error requesting secrets")
    return
  }

  if getSecretsResponse.Encryption == "" {
    t.Errorf("Expected response to have an encrypted key")
    return
  }

  getSecretsResponse, err = api.GetSecretsForSecret(getSecretsResponse.ReadOnly)
  if err != nil {
    t.Errorf("Error requesting secrets for secret: %s", getSecretsResponse.ReadOnly)
    return
  }

  if getSecretsResponse.ReadOnly == "" {
    t.Errorf("Expected response to have a read only key")
    return
  }

  Log("Testing GetFolderPrefs")

  getFolderPrefsResponse, err := api.GetFolderPrefs(testDir.Secret)
  if err != nil {
    t.Errorf("Error requesting prefs for folder")
    return
  }

  if getFolderPrefsResponse.SearchLAN != 1 {
    t.Errorf("Exepected search_lan to be 1")
    return
  }

  Log("Testing SetFolderPrefs")

  prefs := &FolderPreferences{
    SearchLAN:      1,
    SelectiveSync:  1,
    UseDHT:         1,
    UseHosts:       1,
    UseRelayServer: 1,
    UseSyncTrash:   1,
    UseTracker:     1,
  }

  _, err = api.SetFolderPrefs(testDir.Secret, prefs)
  if err != nil {
    t.Errorf("Error making request to set folder preferences")
    return
  }

}

func TestOther(t *testing.T) {
  api := New(*login, *password, *port, *verbose)

  // BUG(aaron): I don't know what the actual list of OS names is.
  osList := []string{"win32", "mac", "linux"}

  Log("Testing GetOS")

  getOSResponse, err := api.GetOS()
  if err != nil {
    t.Errorf("Error getting OS: %s", err)
    return
  }

  found := false
  for _, os := range osList {
    if os == getOSResponse.Name {
      found = true
    }
  }

  if !found {
    t.Errorf("Expected OS to be one of: %s", strings.Join(osList, ","))
    return
  }

  Log("Testing GetVersion")

  getVersionResponse, err := api.GetVersion()
  if err != nil {
    t.Errorf("Error getting version: %s", err)
    return
  }

  if getVersionResponse.Version == "" {
    t.Errorf("Expected version to not be empty")
    return
  }

  Log("Testing GetSpeed")

  _, err = api.GetSpeed()
  if err != nil {
    t.Error("Error getting speed: %s", err)
    return
  }
}

type TestStruct struct {
  Foo string `json:"foo"`
  Bar string `json:"bar"`
}

// Test utility functions.
func TestUtils(t *testing.T) {
  s := &TestStruct{
    Foo: "foo",
    Bar: "bar",
  }

  Log("Testing structToMap")

  m := structToMap(s)

  if m["foo"] != "foo" {
    t.Errorf("Error converting struct to map")
  }
}

func TestCleanup(t *testing.T) {
  api := New(*login, *password, *port, *verbose)

  Log("Tearing down")

  folders, err := api.GetFolders()
  if err != nil {
    t.Errorf("Error getting folders for cleanup.")
    return
  }

  for _, folder := range *folders {
    if strings.HasPrefix(path.Base(folder.Dir), Prefix) {
      Debug("Cleaning up %s", folder.Dir)

      _, err := api.RemoveFolder(folder.Secret)
      if err != nil {
        t.Errorf("Error removing %s", folder.Dir)
        return
      }
    }
  }

  Log("☺")
}
