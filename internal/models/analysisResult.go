package models

type AnalysisResult struct {
	Filename string `json:"filename"`
	Lines    int    `json:"lines"`
	Words    int    `json:"words"`
	Chars    int    `json:"chars"`
}
