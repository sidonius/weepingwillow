package main

type Config struct {
	Log struct {
		ProductionMode bool   `json:"production_mode" toml:"production_mode"`
		OutputPath     string `json:"output_file_path" toml:"output_file_path"`
		Size           int    `json:"size" toml:"size"`
		Backups        int    `json:"backups" toml:"backups"`
		Age            int    `json:"age" toml:"age"`
		ToConsole      bool   `json:"to_console" toml:"to_console"`
	} `json:"log" toml:"log"`

	WebService struct {
		GinMode        string `json:"gin_mode" toml:"gin_mode"`
		Port           string `json:"port" toml:"port"`
		JwtTimeout     int    `json:"jwt_timeout" toml:"jwt_timeout"`
		JwtRefreshTime int    `json:"jwt_refresh_time" toml:"jwt_refresh_time"`
	} `json:"web_service" toml:"web_service"`
}
