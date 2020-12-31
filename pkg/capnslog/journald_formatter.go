// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// +build !windows

package capnslog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func NewJournaldFormatter() (Formatter, error) {
	if !Enabled() {
		return nil, errors.New("No systemd detected")
	}
	return &journaldFormatter{}, nil
}

type journaldFormatter struct{}

func (j *journaldFormatter) Format(pkg string, l LogLevel, _ int, entries ...interface{}) {
	var pri Priority
	switch l {
	case CRITICAL:
		pri = PriCrit
	case ERROR:
		pri = PriErr
	case WARNING:
		pri = PriWarning
	case NOTICE:
		pri = PriNotice
	case INFO:
		pri = PriInfo
	case DEBUG:
		pri = PriDebug
	case TRACE:
		pri = PriDebug
	default:
		panic("Unhandled loglevel")
	}
	msg := fmt.Sprint(entries...)
	tags := map[string]string{
		"PACKAGE":           pkg,
		"SYSLOG_IDENTIFIER": filepath.Base(os.Args[0]),
	}
	err := Send(msg, pri, tags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (j *journaldFormatter) Flush() {}
