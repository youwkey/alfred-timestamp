package main

import (
	"flag"
	"github.com/youwkey/alfred-go"
	"strconv"
	"time"
)

const dateLayout = "2006-01-02"
const datetimeLayout = dateLayout + " 15:04:05"

func main() {
	flag.Parse()
	arg1 := flag.Arg(0)
	arg2 := flag.Arg(1)
	sf := alfred.ScriptFilter{}
	sf.SetEmptyTitle("Parse Error", "")

	if num, err := strconv.ParseInt(arg1, 10, 64); err == nil {
		t := time.Unix(num, 0)
		dateString := t.Format(datetimeLayout)
		sf.Append(alfred.Item{
			Title: dateString,
			Arg:   dateString,
			Text: &alfred.Text{
				Copy:      dateString,
				LargeType: dateString,
			},
		})
	} else if t, err := time.ParseInLocation(datetimeLayout, arg1+" "+arg2, time.Local); err == nil {
		tsString := strconv.FormatInt(t.Unix(), 10)
		sf.Append(alfred.Item{
			Title: tsString,
			Arg:   tsString,
			Text: &alfred.Text{
				Copy:      tsString,
				LargeType: tsString,
			},
		})
	} else if t, err := time.ParseInLocation(dateLayout, arg1, time.Local); err == nil {
		tsString := strconv.FormatInt(t.Unix(), 10)
		sf.Append(alfred.Item{
			Title: tsString,
			Arg:   tsString,
			Text: &alfred.Text{
				Copy:      tsString,
				LargeType: tsString,
			},
		})
	}

	sf.Output()
}
