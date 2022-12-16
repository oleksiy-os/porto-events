package notion

// TODO: Notion integration not ready to use. Under development

// Notion internal config data
type Notion struct {
	Token        string `toml:"token"`
	PageEventsId string `toml:"page_id_events"`
	PageConfigId string `toml:"page_id_config"`
	Timer        uint8  `toml:"timer_check"` // how often check events in source, in hours
}
