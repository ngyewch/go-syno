package filestation

import (
	"encoding/json"
	"github.com/ngyewch/go-syno/api"
	"io"
	"strconv"
	"strings"
)

var (
	fileStationErrorCodes = map[int]string{
		400: "Invalid parameter of file operation",
		401: "Unknown error of file operation",
		402: "System is too busy",
		403: "Invalid user does this file operation",
		404: "Invalid group does this file operation",
		405: "Invalid user and group does this file operation",
		406: "Can't get user/group information from the account server",
		407: "Operation not permitted",
		408: "No such file or directory",
		409: "Non-supported file system",
		410: "Failed to connect internet-based file system (e.g., CIFS)",
		411: "Read-only file system",
		412: "Filename too long in the non-encrypted file system",
		413: "Filename too long in the encrypted file system",
		414: "File already exists",
		415: "Disk quota exceeded",
		416: "No space left on device",
		417: "Input/output error",
		418: "Illegal name or path",
		419: "Illegal file name",
		420: "Illegal file name on FAT file system",
		421: "Device or resource busy",
		599: "No such task of the file operation",
	}
)

type Api struct {
	client *api.Client
}

type ListShareRequest struct {
	Offset        int
	Limit         int
	SortBy        string
	SortDirection string
	OnlyWritable  bool
	Additional    []string
}

type ListShareResponse struct {
	Total  int            `json:"total"`
	Offset int            `json:"offset"`
	Shares []SharedFolder `json:"shares,omitempty"`
}

type ListRequest struct {
	FolderPath    string
	Offset        int
	Limit         int
	SortBy        string
	SortDirection string
	Pattern       []string
	FileType      string
	GotoPath      string
	Additional    []string
}

type ListResponse Folder

type GetInfoRequest struct {
	Path       []string
	Additional []string
}

type GetInfoResponse struct {
	Files []File `json:"files,omitempty"`
}

type DownloadRequest struct {
	Path []string
	Mode string
}

type Folder struct {
	Total  int    `json:"total"`
	Offset int    `json:"offset"`
	Files  []File `json:"files,omitempty"`
}

type File struct {
	Path       string          `json:"path,omitempty"`
	Name       string          `json:"name,omitempty"`
	IsDir      bool            `json:"isdir"`
	Children   *Folder         `json:"children,omitempty"`
	Additional *FileAdditional `json:"additional,omitempty"`
}

type FileAdditional struct {
	RealPath       string       `json:"real_path,omitempty"`
	Size           int64        `json:"size"`
	Owner          *Owner       `json:"owner,omitempty"`
	Time           *Time        `json:"time,omitempty"`
	Perm           *Permissions `json:"perm,omitempty"`
	MountPointType string       `json:"mount_point_type,omitempty"`
	Type           string       `json:"type,omitempty"`
}

type Permissions struct {
	Posix     int  `json:"posix"`
	IsAclMode bool `json:"is_acl_mode"`
	Acl       *Acl `json:"acl,omitempty"`
}

type SharedFolder struct {
	Path       string                  `json:"path,omitempty"`
	Name       string                  `json:"name,omitempty"`
	Additional *SharedFolderAdditional `json:"additional,omitempty"`
}

type SharedFolderAdditional struct {
	RealPath       string                   `json:"real_path,omitempty"`
	Owner          *Owner                   `json:"owner,omitempty"`
	Time           *Time                    `json:"time,omitempty"`
	Perm           *SharedFolderPermissions `json:"perm,omitempty"`
	MountPointType string                   `json:"mount_point_type,omitempty"`
	VolumeStatus   *VolumeStatus            `json:"volume_status,omitempty"`
}

type Owner struct {
	User  string `json:"user,omitempty"`
	Group string `json:"group,omitempty"`
	Uid   int    `json:"uid,omitempty"`
	Gid   int    `json:"gid,omitempty"`
}

type Time struct {
	Atime  int64 `json:"atime,omitempty"`
	Mtime  int64 `json:"mtime,omitempty"`
	Ctime  int64 `json:"ctime,omitempty"`
	Crtime int64 `json:"crtime,omitempty"`
}

type SharedFolderPermissions struct {
	ShareRight string    `json:"share_right,omitempty"`
	Posix      int       `json:"posix"`
	AdvRight   *AdvRight `json:"adv_right,omitempty"`
	AclEnable  bool      `json:"acl_enable"`
	IsAclMode  bool      `json:"is_acl_mode"`
	Acl        *Acl      `json:"acl,omitempty"`
}

type AdvRight struct {
	DisableDownload bool `json:"disable_download"`
	DisableList     bool `json:"disable_list"`
	DisableModify   bool `json:"disable_modify"`
}

type Acl struct {
	Append bool `json:"append"`
	Del    bool `json:"del"`
	Exec   bool `json:"exec"`
	Read   bool `json:"read"`
	Write  bool `json:"write"`
}

type VolumeStatus struct {
	FreeSpace  int64 `json:"freespace"`
	TotalSpace int64 `json:"totalspace"`
	ReadOnly   bool  `json:"readonly"`
}

func New(client *api.Client) *Api {
	return &Api{
		client: client,
	}
}

func (a *Api) ListShare(req ListShareRequest) (*api.Response[ListShareResponse], error) {
	paramMap := make(map[string]string)
	paramMap["offset"] = strconv.Itoa(req.Offset)
	if req.Limit != 0 {
		paramMap["limit"] = strconv.Itoa(req.Limit)
	}
	if req.SortBy != "" {
		paramMap["sort_by"] = req.SortBy
	}
	if req.SortDirection != "" {
		paramMap["sort_direction"] = req.SortDirection
	}
	if req.OnlyWritable {
		paramMap["onlywritable"] = strconv.FormatBool(req.OnlyWritable)
	}
	if len(req.Additional) > 0 {
		jsonBytes, err := json.Marshal(req.Additional)
		if err != nil {
			return nil, err
		}
		paramMap["additional"] = string(jsonBytes)
	}

	var res api.Response[ListShareResponse]
	err := a.client.Request("SYNO.FileStation.List", 2, "list_share", paramMap, &res)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &res, nil
}

func (a *Api) List(req ListRequest) (*api.Response[ListResponse], error) {
	paramMap := make(map[string]string)
	if req.FolderPath != "" {
		paramMap["folder_path"] = req.FolderPath
	}
	paramMap["offset"] = strconv.Itoa(req.Offset)
	if req.Limit != 0 {
		paramMap["limit"] = strconv.Itoa(req.Limit)
	}
	if req.SortBy != "" {
		paramMap["sort_by"] = req.SortBy
	}
	if req.SortDirection != "" {
		paramMap["sort_direction"] = req.SortDirection
	}
	if len(req.Pattern) > 0 {
		paramMap["pattern"] = strings.Join(req.Pattern, ",")
	}
	if req.FileType != "" {
		paramMap["filetype"] = req.FileType
	}
	if req.GotoPath != "" {
		paramMap["goto_path"] = req.GotoPath
	}
	if len(req.Additional) > 0 {
		jsonBytes, err := json.Marshal(req.Additional)
		if err != nil {
			return nil, err
		}
		paramMap["additional"] = string(jsonBytes)
	}

	var res api.Response[ListResponse]
	err := a.client.Request("SYNO.FileStation.List", 2, "list", paramMap, &res)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &res, nil
}

func (a *Api) GetInfo(req GetInfoRequest) (*api.Response[GetInfoResponse], error) {
	paramMap := make(map[string]string)
	if len(req.Path) > 0 {
		jsonBytes, err := json.Marshal(req.Path)
		if err != nil {
			return nil, err
		}
		paramMap["path"] = string(jsonBytes)
	}
	if len(req.Additional) > 0 {
		jsonBytes, err := json.Marshal(req.Additional)
		if err != nil {
			return nil, err
		}
		paramMap["additional"] = string(jsonBytes)
	}

	var res api.Response[GetInfoResponse]
	err := a.client.Request("SYNO.FileStation.List", 2, "getinfo", paramMap, &res)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &res, nil
}

func (a *Api) Download(req DownloadRequest) (io.ReadCloser, error) {
	paramMap := make(map[string]string)
	if len(req.Path) > 0 {
		jsonBytes, err := json.Marshal(req.Path)
		if err != nil {
			return nil, err
		}
		paramMap["path"] = string(jsonBytes)
	}
	if req.Mode != "" {
		paramMap["mode"] = req.Mode
	}
	return a.client.RawRequest("SYNO.FileStation.Download", 2, "download", paramMap)
}
