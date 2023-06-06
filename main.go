package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

var (
	sheet  string
	format string
)

func init() {
	flag.StringVar(&sheet, "sheet", "", "Sheet name (defaults to first sheet)")
	flag.StringVar(&format, "format", "csv", `Output format ("csv", "json", "yaml")`)
}

func main() {
	flag.Parse()
	signal.Ignore(syscall.SIGPIPE)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	f, err := excelize.OpenReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if sheet == "" && len(sheets) > 0 {
		sheet = sheets[0]
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
	}

	if len(rows) == 0 {
		return
	}

	w := os.Stdout

	switch format {
	case "csv":
		err = toCSV(w, rows)
	case "json":
		err = toJSON(w, rows)
	case "yaml":
		err = encode(w, rows, yaml.NewEncoder(w).Encode)
	default:
		err = fmt.Errorf("invalid formt %s", format)
	}
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(w, "\n")
}

func toCSV(w io.Writer, rows [][]string) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	for _, r := range rows {
		for i, l := range r {
			r[i] = normalize(l)
		}
		if len(r) == 0 {
			continue
		}
		if err := cw.Write(r); err != nil {
			if !errors.Is(err, syscall.EPIPE) {
				return err
			}
		}
	}
	return nil
}

func toJSON(w io.Writer, rows [][]string) error {
  enc := json.NewEncoder(w)
  enc.SetIndent("", "  ")
	return encode(w, rows, enc.Encode)
}

func encode(w io.Writer, rows [][]string, encodeFn func(v any) error) error {
	header := rows[0]
	for i, h := range header {
		header[i] = normalize(h)
	}
	list := make([]map[string]any, 0, len(rows)-1)
	for _, row := range rows[1:] {
		item := make(map[string]any, len(row))
		for i, c := range row {
			c = normalize(c)
			var v any
			if n, ok := isInt(c); ok {
				v = n
			} else if n, ok := isFloat(c); ok {
				v = n
			} else if b, ok := isBool(c); ok {
				v = b
			} else {
				v = c
			}
			item[header[i]] = v
		}
		list = append(list, item)
	}
	return encodeFn(list)
}

var rep = strings.NewReplacer("\n", " ", "\r\n", " ")

func normalize(s string) string {
	return strings.TrimSpace(rep.Replace(s))
}

func isInt(s string) (int64, bool) {
	n, err := strconv.ParseInt(s, 10, 64)
	return n, err == nil
}

func isFloat(s string) (float64, bool) {
	n, err := strconv.ParseFloat(s, 64)
	return n, err == nil
}

func isBool(s string) (bool, bool) {
	b, err := strconv.ParseBool(s)
	return b, err == nil
}
