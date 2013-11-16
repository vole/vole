package btsync_api

type Folder struct {
  Dir      string `json:"dir"`
  Secret   string `json:"secret"`
  Size     int64  `json:"size"`
  Type     string `json:"type"`
  Files    int64  `json:"files"`
  Error    int    `json:"error"`
  Indexing int    `json:"indexing"`
}

type File struct {
  HavePieces  int    `json:"have_pieces"`
  Name        string `json:"name"`
  Size        int64  `json:"size"`
  State       string `json:"state"`
  TotalPieces string `json:"total_pieces"`
  Type        string `json:"type"`
  Download    int    `json:"download"`
}

type Peer struct {
  Id         string `json:"id"`
  Connection string `json:"connection"`
  Name       string `json:"name"`
  Synced     int    `json:"synced"`
  Download   int64  `json:"download"`
  Upload     int64  `json:"upload"`
}

type Preferences struct {
  DeviceName                 string `json:"device_name"`
  DiskLowPriority            bool   `json:"disk_low_priority"`
  DownloadLimit              int    `json:"download_limit"`
  FolderRescanInterval       int    `json:"folder_rescan_interval"`
  LANEncryptData             bool   `json:"lan_encrypt_data"`
  LANUseTCP                  bool   `json:"lan_use_tcp"`
  Lang                       int    `json:"lang"`
  ListeningPort              int    `json:"listening_port"`
  MaxFileSizeDiffForPatching int64  `json:"max_file_size_diff_for_patching"`
  MaxFileSizeForVersioning   int64  `json:"max_file_size_for_versioning"`
  RateLimitLocalPeers        bool   `json:"rate_limit_local_peers"`
  ReadPoolSize               int64  `json:"read_pool_size"`
  SyncMaxTimeDiff            int64  `json:"sync_max_time_diff"`
  SyncTrashTTL               int64  `json:"sync_trash_ttl"`
  UploadLimit                int64  `json:"upload_limit"`
  UseUPnP                    int    `json:"use_upnp"`
  WritePoolSize              int    `json:"write_pool_size"`
}

type FolderPreferences struct {
  SearchLAN      int `json:"search_lan"`
  SelectiveSync  int `json:"selective_sync"`
  UseDHT         int `json:"use_dht"`
  UseHosts       int `json:"use_hosts"`
  UseRelayServer int `json:"use_relay_server"`
  UseSyncTrash   int `json:"use_sync_trash"`
  UseTracker     int `json:"use_tracker"`
}

type Response struct {
  Error   int    `json:"error"`
  Message string `json:"message"`
}

type GetFoldersResponse []Folder

type GetFilesResponse []File

type SetFilePrefsResponse []File

type GetFolderPeersResponse []Peer

type GetSecretsResponse struct {
  ReadOnly   string `json:"read_only"`
  ReadWrite  string `json:"read_write"`
  Encryption string `json:"encryption"`
}

type GetFolderPrefsResponse FolderPreferences
type SetFolderPrefsResponse FolderPreferences

type GetFolderHostsResponse []string

type GetPreferencesResponse Preferences

type GetOSResponse struct {
  Name string `json:"os"`
}

type GetVersionResponse struct {
  Version string `json:"version"`
}

type GetSpeedResponse struct {
  Download int64 `json:"download"`
  Upload   int64 `json:"upload"`
}
