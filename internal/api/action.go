package api

type Update struct {
	Filename string `json:"filename" yaml:"filename"`
	URL      string `json:"url" yaml:"url"`
}