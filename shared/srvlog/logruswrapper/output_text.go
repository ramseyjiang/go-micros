package logruswrapper

import (
	"bytes"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	grey   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 36
	white  = 37
)

type IQTextFormatter struct {
	UseColour         bool
	IncludeTimePrefix bool
	TimePrefixFormat  string
}

func (iqtf *IQTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	levelText := "????"
	traceColour := green

	var levelColor int
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = blue
		levelText = "INFO"
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = white
		traceColour = grey
		levelText = "DBUG"
	case logrus.WarnLevel:
		levelColor = yellow
		levelText = "WARN"
	case logrus.ErrorLevel:
		levelColor = red
		levelText = "ERRR"
	case logrus.FatalLevel:
		levelColor = red
		levelText = "FATL"
	case logrus.PanicLevel:
		levelColor = red
		levelText = "PANC"
	default:
		levelColor = blue
	}

	// Log timestamp
	if iqtf.IncludeTimePrefix {
		if iqtf.TimePrefixFormat == "" {
			fmt.Fprintf(b, "%s ", entry.Time.Format(time.StampMicro))
		} else {
			fmt.Fprintf(b, "%s ", entry.Time.Format(iqtf.TimePrefixFormat))
		}
	}

	if iqtf.UseColour {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m ", levelColor, levelText)
	} else {
		fmt.Fprintf(b, "%s ", levelText)
	}

	// Write [package.FunctionName:Line]
	originText := ""
	if valFunction, functionExists := entry.Data["source_func"]; functionExists {
		if valFile, fileExists := entry.Data["source_file"]; fileExists {
			if valLine, lineExists := entry.Data["source_line"]; lineExists {
				originText = fmt.Sprintf("[%v %v:%v]", valFunction, valFile, valLine)
				delete(entry.Data, "source_line")
			} else {
				originText = fmt.Sprintf("[%v %v]", valFunction, valFile)
			}
			delete(entry.Data, "source_file")
		} else {
			originText = fmt.Sprintf("[%v]", valFunction)
		}
		delete(entry.Data, "source_func")
	} else if valFunction, functionExists := entry.Data["origin_func"]; functionExists {
		if valFile, fileExists := entry.Data["origin_file"]; fileExists {
			if valLine, lineExists := entry.Data["origin_line"]; lineExists {
				originText = fmt.Sprintf("[%v %v:%v]", valFunction, valFile, valLine)
				delete(entry.Data, "origin_line")
			} else {
				originText = fmt.Sprintf("[%v %v]", valFunction, valFile)
			}
			delete(entry.Data, "origin_file")
		} else {
			originText = fmt.Sprintf("[%v]", valFunction)
		}
		delete(entry.Data, "origin_func")
	} else if valApp, appExists := entry.Data["app"]; appExists {
		// We don't have any calling information
		originText = fmt.Sprintf("%v", valApp)
	}

	if originText != "" {
		if iqtf.UseColour {
			fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m ", traceColour, originText)
		} else {
			fmt.Fprintf(b, "%s ", originText)
		}
	}

	fmt.Fprintf(b, "%s ", entry.Message)
	// b.Write([]byte("test"))

	b.WriteByte('\n')
	return b.Bytes(), nil
}
