// Copyright 2022 youwkey. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"strconv"
	"strings"
	"time"

	"github.com/youwkey/alfred-go"
)

const (
	dateLayout            = "2006-01-02"
	datetimeLayout        = dateLayout + " 15:04:05"
	unixTimeMaxDigit      = 12
	unixMilliTimeMaxDigit = 15
	unixMicroTimeMaxDigit = 18
)

// parse errors.
var (
	ErrParseTimestamp   = errors.New("parse timestamp error")
	ErrUnsupportedDigit = errors.New("unsupported unix timestamp digit")
	ErrParseDateString  = errors.New("parse date string error")
)

// ParseTimestamp returns the time.Time value represented by the timestamp string.
func ParseTimestamp(str string) (time.Time, error) {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}, ErrParseTimestamp
	}

	digit := len(str)

	switch {
	case digit <= unixTimeMaxDigit:
		return time.Unix(num, 0), nil
	case digit <= unixMilliTimeMaxDigit:
		return time.UnixMilli(num), nil
	case digit <= unixMicroTimeMaxDigit:
		return time.UnixMicro(num), nil
	default:
		return time.Time{}, ErrUnsupportedDigit
	}
}

// ParseDateString returns the time.Time value represented by the date string.
func ParseDateString(dateString, timeString string) (time.Time, error) {
	dateTimeLayouts := [...]string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.000Z07:00",
		"2006-01-02T15:04:05.000000Z07:00",
	}
	dateLayouts := [...]string{"2006-01-02", "2006/01/02", "20060102"}
	timeLayouts := [...]string{"15", "15:04", "15:04:05", "15:04:05.000", "15:04:05.000000"}

	local := time.Local

	if strings.Contains(dateString, "T") {
		for _, layout := range dateTimeLayouts {
			if t, err := time.ParseInLocation(layout, dateString, local); err == nil {
				return t, nil
			}
		}

		return time.Time{}, ErrParseDateString
	}

	var baseTime time.Time

	for _, layout := range dateLayouts {
		if t, err := time.ParseInLocation(layout, dateString, local); err == nil {
			baseTime = t

			break
		}
	}

	if baseTime.IsZero() {
		return time.Time{}, ErrParseDateString
	} else if timeString == "" {
		return baseTime, nil
	}

	for _, layout := range timeLayouts {
		if t, err := time.ParseInLocation(layout, timeString, local); err == nil {
			year, month, day := baseTime.Date()
			hour, min, sec := t.Clock()
			nsec := t.Nanosecond()

			return time.Date(year, month, day, hour, min, sec, nsec, local), nil
		}
	}

	return time.Time{}, ErrParseDateString
}

func main() {
	flag.Parse()
	arg1 := flag.Arg(0)
	arg2 := flag.Arg(1)
	sf := alfred.ScriptFilter{}

	if t, err := ParseTimestamp(arg1); err == nil {
		dateString := t.Format(datetimeLayout)
		item := alfred.NewItem(dateString).Arg(dateString).Text(dateString)
		sf.Items().Append(item)
	} else if t, err := ParseDateString(arg1, arg2); err == nil {
		tsString := strconv.FormatInt(t.Unix(), 10)
		item := alfred.NewItem(tsString).Arg(tsString).Text(tsString)
		sf.Items().Append(item)
	} else {
		sf.Items().Append(alfred.NewInvalidItem("Parse Error"))
	}

	_ = sf.Output()
}
