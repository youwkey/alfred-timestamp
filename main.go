package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/youwkey/alfred-go"
)

const (
	dateLayout     = "2006-01-02"
	datetimeLayout = dateLayout + " 15:04:05"
)

func main() {
	flag.Parse()
	arg1 := flag.Arg(0)
	arg2 := flag.Arg(1)
	sf := alfred.ScriptFilter{}

	if num, err := strconv.ParseInt(arg1, 10, 64); err == nil {
		t := time.Unix(num, 0)
		dateString := t.Format(datetimeLayout)
		item := alfred.NewItem(dateString).Arg(dateString).Text(dateString)
		sf.Items().Append(item)
	} else if t, err := time.ParseInLocation(datetimeLayout, arg1+" "+arg2, time.Local); err == nil {
		tsString := strconv.FormatInt(t.Unix(), 10)
		item := alfred.NewItem(tsString).Arg(tsString).Text(tsString)
		sf.Items().Append(item)
	} else if t, err := time.ParseInLocation(dateLayout, arg1, time.Local); err == nil {
		tsString := strconv.FormatInt(t.Unix(), 10)
		item := alfred.NewItem(tsString).Arg(tsString).Text(tsString)
		sf.Items().Append(item)
	} else {
		sf.Items().Append(alfred.NewInvalidItem("Parse Error"))
	}

	_ = sf.Output()
}
