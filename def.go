package main

import "modules/transform_html"

type OutputJson = struct {
	Begin           int
	End             int
	Etl_config_path string
	Source_name     string
	OutputData			any
}


type SourceConfig struct {
	Name     string `yaml:"name"`
	Root_url string
	Menu     struct {
		Page_limit     int
		Cards_per_page int
		Default_url    string
		Page_url_sub   string
		First_page_url string

		Rules []*transform_html.ParserTransfromRule
	}

	Card struct {
		Rules []*transform_html.ParserTransfromRule
	}
}

type HttpConfigRetries struct {
	Max_retries uint
	Backoff_factor uint
	Timeout uint
	status_forcelist []uint
}

type HttpConfig struct {
	Retries HttpConfigRetries
	Headers map[string]string	
}

type EtlConfig struct {
	Http HttpConfig
	Sources []SourceConfig `yaml:"sources"`
}