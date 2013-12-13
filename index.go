package main

import (
	"encoding/csv"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Town struct {
	ID        int
	Name      string
	AltNames  []string
	State     string
	Latitude  float64
	Longitude float64
}

type TownRef struct {
	Name string
	Town *Town
}

var townIndex townRefList

type townRefList []TownRef

func (t townRefList) Len() int           { return len(t) }
func (t townRefList) Less(i, j int) bool { return t[i].Name < t[j].Name }
func (t townRefList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func indexTown(t *Town) {
	names := []string{strings.ToLower(t.Name)}
	for _, n := range t.AltNames {
		names = appendIfMissing(names, n)
	}
	if s := stripPrefixes(t.Name, []string{"Township of ", "Town of ", "City of ", "Village of ", "Borough of "}); s != "" {
		names = appendIfMissing(names, s)
	}
	for _, n := range names {
		townIndex = append(townIndex, TownRef{n, t})
	}
}

func appendIfMissing(b []string, s string) []string {
	s = strings.ToLower(s)
	for _, v := range b {
		if s == v {
			return b
		}
	}
	return append(b, s)
}

func stripPrefixes(s string, prefixes []string) string {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return s[len(p):]
		}
	}
	return ""
}

func indexTowns(f io.Reader) error {
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.Comma = '\t'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		var t Town
		t.ID, _ = strconv.Atoi(record[0])
		t.Name = record[1]
		t.Latitude, _ = strconv.ParseFloat(record[4], 64)
		t.Longitude, _ = strconv.ParseFloat(record[5], 64)
		t.State = record[10]

		for _, n := range strings.Split(record[3], ",") {
			if n = strings.TrimSpace(n); n != "" {
				t.AltNames = append(t.AltNames, n)
			}
		}

		indexTown(&t)
	}

	sort.Sort(townIndex)
	return nil
}

func searchTowns(p string) []*Town {
	p = strings.ToLower(p)
	n := sort.Search(len(townIndex), func(i int) bool { return strings.HasPrefix(townIndex[i].Name, p) || townIndex[i].Name >= p })
	if n >= len(townIndex) || !strings.HasPrefix(townIndex[n].Name, p) {
		return nil
	}
	var res []*Town
outer:
	for _, t := range townIndex[n:] {
		if !strings.HasPrefix(t.Name, p) {
			break
		}
		for _, r := range res {
			if r == t.Town {
				continue outer
			}
		}
		res = append(res, t.Town)
	}
	return res
}
