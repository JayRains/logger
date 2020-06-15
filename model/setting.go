package model

type Logger struct {
	Level       string `json:"level"`
	Identifier  string `json:"identifier"`
	TimeFormat  string `json:"time_format"`
	FileCording bool   `json:"file_cording"`
	OpenColor       bool   `json:"open_color"`
	SavePath    string `json:"save_path"`
}

