package common

import (
	"crypto/md5"
	"math/rand"
	"strconv"
	"strings"
	"errors"
	"time"
	"fmt"
	"io"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}

		return false
	}

	return true
}


func MD5Calc(input string) string {

	h := md5.New()
	io.WriteString(h, input)
	output := fmt.Sprintf("%X", h.Sum(nil))

	return output
}


func Gen6Code() string {

   return fmt.Sprintf("%06v",
	rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}


func GetClassInfo(level int) (string, string, error) {

	var code, text string
	var err error

	if level == CLASS_LEVEL_MAIN {
		code = "bic_zklad2"
		text = "zklad2txt"

	} else if level == CLASS_LEVEL_MID {
		code = "bic_zklasse_d"
		text = "zklasse_dtxt"

	} else if level == CLASS_LEVEL_SUB {
		code = "bic_zrpa_mtl"
		text = "zrpa_mtltxt"

	} else {
		err = errors.New(ERR_MSG_INVALID_CLASS)
	}

	return code, text, err
}


func ParseRelativeTime(date string) (string, string, error) {

	var bDate, eDate string
	var v int

	last := date[len(date)-1:]

	now := time.Now()

	eTime := now.AddDate(0, 0, -1)
	eDate = eTime.Format("2006-01-02")

	if last == "d" {
		fmt.Sscanf(date, "%dd", &v)

		bTime := now.AddDate(0, 0, 0-v)
		bDate = bTime.Format("2006-01-02")
	} else if last == "m" {
		fmt.Sscanf(date, "%dm", &v)

		bTime := now.AddDate(0, 0-v, 0)
		bDate = bTime.Format("2006-01-02")
	} else {
		return "", "", errors.New("invalid date")
	}

	return bDate, eDate, nil
}


/*
func GetRelativeDate(date string, days int) string {

	const shortForm = "2006-01-02"

        d, _ := time.Parse(shortForm, date)

	t := d.AddDate(0, 0, days*(-1))

        return t.Format(shortForm)
}
*/

func GetSubDay(date string, days int) string {

	const shortForm = "2006-01-02"

        d, _ := time.Parse(shortForm, date)

	t := d.AddDate(0, 0, days*(-1))

        return t.Format(shortForm)
}


func GetAddDay(date string, days int) string {

	const shortForm = "2006-01-02"

        d, _ := time.Parse(shortForm, date)

	t := d.AddDate(0, 0, days)

        return t.Format(shortForm)
}


func GetSubMonth(date string, mons int) string {

	const shortForm = "2006-01-02"

        d, _ := time.Parse(shortForm, date)

	t := d.AddDate(0, mons*(-1), 0)

        return t.Format(shortForm)
}


func GetSubYear(date string, years int) string {

	const shortForm = "2006-01-02"

        d, _ := time.Parse(shortForm, date)

	t := d.AddDate(years*(-1), 0, 0)

        return t.Format(shortForm)
}


func GetSubDayCount(begin string, end string) float64 {

	b, _ := time.Parse("2006-01-02", begin)
	e, _ := time.Parse("2006-01-02", end)
	d := e.Sub(b)

	return (d.Hours() / 24)
}


func GetShortVendorCode(code string) string {

	if len(code) == 10 {
		return code[3:len(code)]
	}

	return code
}


func GetLongVendorCode(code string) string {

	if len(code) == 7 {
		code = fmt.Sprintf("000%s", code)
		return code
	}

	return code
}


func GetShortMatlCode(code string) string {

	if len(code) == 18 {
		return code[10:len(code)]
	}

	return code
}


func GetLongMatlCode(code string) string {

	if len(code) == 8 {
		code = fmt.Sprintf("0000000000%s", code)
		return code
	} else if len(code) == 6 {
		code = fmt.Sprintf("000000000000%s", code)
		return code
	}

	return code
}


func GetLastMonth(month string) (string, error) {

	if len(month) != 6 {
		return "", errors.New("invalid month")
	}

	imon, err := strconv.Atoi(month)
	if err != nil {
		return "", err
	}

	m := month[4:6]
	if m == "01" {
		imon = imon - 89
	} else {
		imon = imon - 1
	}

	return strconv.Itoa(imon), nil
}


func GetFirstMonth(month string) (string, error) {

	if len(month) != 6 {
		return "", errors.New("invalid month")
	}

	yyyy  := month[0:4]
	first := fmt.Sprintf("%s01", yyyy)

	return first, nil
}


func GetClassCodeLevel(code string) (int, string) {

	var level int
	size := len(code)

	if size == 4 {
		level = CLASS_LEVEL_MAIN

	} else if size == 6 {
		//XXX TODO: 40xxxx style material code
		level = CLASS_LEVEL_MID

	} else if size == 8 {
		h2 := code[0:2]
		if h2 == "70" {
			//XXX TODO: exclude 70010101, 70020101
			level = CLASS_LEVEL_MATL
			code = fmt.Sprintf("0000000000%s", code)
			return level, code
		}

		h4 := code[0:4]
		if h4 == "5500" || h4 == "7500" || h4 == "7940" ||
			h4 == "8000" || h4 == "9000" {
			level = CLASS_LEVEL_MATL
			code = fmt.Sprintf("0000000000%s", code)
			return level, code
		}

		level = CLASS_LEVEL_SUB

	} else {
		level = CLASS_LEVEL_NONE
	}

	return level, code
}


func TranslateDate(date string) string {

	if len(date) != 10 {
		return date
	}

	return strings.Replace(date, ".", "-", -1)
}


