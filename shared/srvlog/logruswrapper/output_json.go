package logruswrapper

import "github.com/sirupsen/logrus"

type IQJSONFormatter struct {
	lf logrus.JSONFormatter
}

func (iqjf *IQJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Custom formatter here
	return iqjf.lf.Format(entry)
}
