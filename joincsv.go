package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	skip := true
	n := 1
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	switch os.Args[1] {
	case "-h", "--help":
		fmt.Println(`joincsv
Sample usage: joincsv labels.csv content1.csv content2.csv

joincsv accepts two or more CSV files as inputs.

The first CSV file must contain two or more rows of labels.
The first row in that CSV file should be your desired new headings in order.
The second (and subsequent) rows should have those headings in the
columns where those values live in your content files.

The second, and remaining, CSV files should all contain contents.

This program will join those together and apply label headings from
your label CSV.

There is one optional flag "-k", or "--keep". Use this if your content
CSVs don't have header rows (otherwise you'll lose your first row of data)!`)
		os.Exit(0)
	case "-k", "--keep":
		if len(os.Args) < 4 {
			fmt.Println("Error: you need to provide at least two CSV files (labels and contents)")
			os.Exit(1)
		}
		skip = false
		n = 2
	default:
		if len(os.Args) < 3 {
			if len(os.Args) < 4 {
				fmt.Println("Error: you need to provide at least two CSV files (labels and contents)")
				os.Exit(1)
			}
		}
	}
	var idxs [][]int
	outcsv := csv.NewWriter(os.Stdout)
	for i, path := range os.Args[n:] {
		c, err := readCSV(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if i == 0 {
			if len(c) < 2 {
				fmt.Println("Error: your labels csv must contain at least two rows (see help -h)")
				os.Exit(1)
			}
			idxs = flatten(shrink(c[0]), labels(c[1:]))
			outcsv.Write(shrink(c[0]))
			continue
		}
		if skip {
			c = skipHeader(c)
		}
		for _, r := range c {
			outcsv.Write(row(idxs, r))
		}
	}
	outcsv.Flush()
	os.Exit(0)
}

// takes a file path, checks if there is a byte-order mark (BOM),
// which it will discard, then interprets the file as a csv,
// and returns the csv contents as a [][]string
func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	ru, _, err := buf.ReadRune()
	if err != nil {
		return nil, err
	}
	if ru != '\uFEFF' {
		buf.UnreadRune() // ignore error because we just read the rune
	}
	return csv.NewReader(buf).ReadAll()
}

// if the contents csv has a header we want to skip, trim the first row from the slice
func skipHeader(c [][]string) [][]string {
	if len(c) < 2 {
		return c
	}
	return c[1:]
}

// join fields: use the CSV package rather than just a simple strings.Join(s, ",") in case the
// fields include characters that need csv escaping
func join(fields []string) string {
	var b strings.Builder
	c := csv.NewWriter(&b)
	err := c.Write(fields)
	if err != nil {
		panic(err) // this ideally shouldn't happen!
	}
	c.Flush()
	return strings.TrimSpace(b.String())
}

// row takes your indexes and a row from the content and returns a row for the output CSV
func row(idxs [][]int, vals []string) []string {
	ret := make([]string, len(idxs))
	for i, v := range idxs {
		if len(v) == 1 {
			ret[i] = vals[v[0]]
		} else if len(v) > 1 {
			j := make([]string, len(v))
			for ii := range j {
				j[ii] = vals[v[ii]]
			}
			ret[i] = join(j)
		}
	}
	return ret
}

// takes values from the labels csv and returns a map that links
// the labels to one or more column indexes in the content csv
func labels(c [][]string) map[string][]int {
	l := make(map[string][]int)
	for _, row := range c {
		for idx, col := range row {
			if col == "" { // can have blank values to ignore a column
				continue
			}
			l[col] = append(l[col], idx) // add the column index to the label
		}
	}
	return l
}

// flatten takes the labels then returns the indexes where those labels
// live in the content
func flatten(h []string, l map[string][]int) [][]int {
	idxs := make([][]int, len(h))
	for i, v := range h {
		idxs[i] = l[v]
	}
	return idxs
}

// shrink removes any trailing empty headers
func shrink(h []string) []string {
	for i := len(h) - 1; i > -1; i-- {
		if h[i] != "" {
			return h[:i+1]
		}
	}
	return h
}
