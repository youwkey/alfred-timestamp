package main

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	t.Parallel()

	type Test struct {
		in  string
		out time.Time
		err error
	}

	date, utc := time.Date, time.UTC
	tests := []Test{
		{in: "1641008096", out: date(2022, 1, 1, 3, 34, 56, 0, utc), err: nil},
		{in: "1641008096123", out: date(2022, 1, 1, 3, 34, 56, 123000000, utc), err: nil},
		{in: "1641008096123456", out: date(2022, 1, 1, 3, 34, 56, 123456000, utc), err: nil},
		{in: "999999999999999999", out: date(33658, 9, 27, 1, 46, 39, 999999000, utc), err: nil},
		{in: "string", out: time.Time{}, err: ErrParseTimestamp},
		{in: "1000000000000000000", out: time.Time{}, err: ErrUnsupportedDigit},
	}

	for i, test := range tests {
		i, test := i, test
		t.Run(fmt.Sprintf("#%d:ParseTimestamp", i), func(t *testing.T) {
			t.Parallel()
			got, err := ParseTimestamp(test.in)
			if !got.Equal(test.out) {
				t.Errorf("#%d: got: %v want: %v", i, got, test.out)
			}
			if !errors.Is(err, test.err) {
				t.Errorf("#%d: got err: %v want err: %v", i, err, test.err)
			}
		})
	}
}

func TestParseDateString(t *testing.T) {
	t.Parallel()

	type Test struct {
		in1 string
		in2 string
		out time.Time
		err error
	}

	utc := time.UTC
	jst, _ := time.LoadLocation("Asia/Tokyo")
	tests := []Test{
		{
			in1: "2022-01-01T12:34:56Z",
			in2: "",
			out: time.Date(2022, 1, 1, 12, 34, 56, 0, utc),
			err: nil,
		},
		{
			in1: "2022-01-01T12:34:56+09:00",
			in2: "",
			out: time.Date(2022, 1, 1, 12, 34, 56, 0, jst),
			err: nil,
		},
		{
			in1: "2022-01-01T12:34:56.123+09:00",
			in2: "",
			out: time.Date(2022, 1, 1, 12, 34, 56, 123000000, jst),
			err: nil,
		},
		{
			in1: "2022-01-01T12:34:56.123456+09:00",
			in2: "",
			out: time.Date(2022, 1, 1, 12, 34, 56, 123456000, jst),
			err: nil,
		},
		{
			in1: "2022T12",
			in2: "",
			out: time.Time{},
			err: ErrParseDateString,
		},
		{
			in1: "2022-12-31",
			in2: "",
			out: time.Date(2022, 12, 31, 0, 0, 0, 0, time.Local),
			err: nil,
		},
		{
			in1: "2022-12-31",
			in2: "01",
			out: time.Date(2022, 12, 31, 1, 0, 0, 0, time.Local),
			err: nil,
		},
		{
			in1: "2022-12-31",
			in2: "01:02",
			out: time.Date(2022, 12, 31, 1, 2, 0, 0, time.Local),
			err: nil,
		},
		{
			in1: "2022/12/31",
			in2: "01:02:03",
			out: time.Date(2022, 12, 31, 1, 2, 3, 0, time.Local),
			err: nil,
		},
		{
			in1: "20221231",
			in2: "01:02:03.456",
			out: time.Date(2022, 12, 31, 1, 2, 3, 456000000, time.Local),
			err: nil,
		},
		{
			in1: "2022-12-31",
			in2: "01:02:03.456789",
			out: time.Date(2022, 12, 31, 1, 2, 3, 456789000, time.Local),
			err: nil,
		},
		{
			in1: "2022",
			in2: "",
			out: time.Time{},
			err: ErrParseDateString,
		},
	}

	for i, test := range tests {
		i, test := i, test
		t.Run(fmt.Sprintf("#%d:ParseDateString", i), func(t *testing.T) {
			t.Parallel()
			got, err := ParseDateString(test.in1, test.in2)
			if !got.Equal(test.out) {
				t.Errorf("#%d: got: %v want: %v", i, got, test.out)
			}
			if !errors.Is(err, test.err) {
				t.Errorf("#%d: got err: %v want err: %v", i, err, test.err)
			}
		})
	}
}
