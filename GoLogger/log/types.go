/**
 */
package log

import (
	"encoding/json"
)

const (
	LEVEL_DEBUG		= 0
	LEVEL_INFO		= 10
	LEVEL_NOTICE	= 20
	LEVEL_WARN		= 30
	LEVEL_ERROR		= 40
	LEVEL_FETAL		= 50
)

type LogText struct {
	Datetime	string	`json:"datetime,omitempty"`
	Timestamp	float64	`json:"timestamp,omitempty"`
	LevelInt	int32	`json:"levelInt,omitempty"`
	LevelStr	string	`json:"levelStr,omitempty"`
	File		string	`json:"file,omitempty"`
	Line		int32	`json:"line,omitempty"`
	Function	string	`json:"function,omitempty"`
	Text		string	`json:"text"`
}

var levelString = [...]string {
	"DEBUG", "DEBUG", "DEBUG", "DEBUG", "DEBUG", "DEBUG", "DEBUG", "DEBUG", "DEBUG", "DEBUG",
	"INFO ", "INFO ", "INFO ", "INFO ", "INFO ", "INFO ", "INFO ", "INFO ", "INFO ", "INFO ",
	"NOTCE", "NOTCE", "NOTCE", "NOTCE", "NOTCE", "NOTCE", "NOTCE", "NOTCE", "NOTCE", "NOTCE",
	"WARN ", "WARN ", "WARN ", "WARN ", "WARN ", "WARN ", "WARN ", "WARN ", "WARN ", "WARN ",
	"ERROR", "ERROR", "ERROR", "ERROR", "ERROR", "ERROR", "ERROR", "ERROR", "ERROR", "ERROR",
	"FETAL", "FETAL", "FETAL", "FETAL", "FETAL", "FETAL", "FETAL", "FETAL", "FETAL", "FETAL",
}

func GetLevelString(level int32) string {
	if (level < int32(len(levelString))) {
		return levelString[level]
	} else {
		return "FETAL"
	}
}

func ParseLogText(text string) (ret LogText, err error) {
	err = json.Unmarshal([]byte(text), &ret)
	return
}
