package configs

import (
	telegramApi "github.com/oleksiy-os/porto-events/internal/model/client/telegram"
	"github.com/oleksiy-os/porto-events/internal/store/notion"
)

type (
	Server struct {
		BindAddr string `toml:"bind_addr"`
	}

	Config struct {
		ProductionMode  bool   `toml:"production_mode"`
		LogLevel        uint8  `toml:"log_level"`
		SourcesListPath string `toml:"sources_list_path"`
		Telegram        telegramApi.Telegram
		Notion          notion.Notion
		Server          Server
	}
)
