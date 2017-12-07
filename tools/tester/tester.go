package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

func main() {
	var v interface{}
	err := json.NewDecoder(os.Stdin).Decode(&v)
	if err != nil {
		log.Fatal(err)
	}
	print(v, "")
}

func print(v interface{}, indent string) {
	switch t := v.(type) {
	case map[string]interface{}:
		fmt.Printf("map[string]interface{}{\n")
		for _, k := range keys(t) {
			nindent := indent + " "
			fmt.Printf("%s%q: ", nindent, k)
			print(t[k], nindent)
			fmt.Printf(",\n")
		}
		fmt.Printf("%s}", indent)
	case []interface{}:
		fmt.Printf("[]interface{}{\n")
		for _, v := range t {
			nindent := indent + " "
			print(v, nindent)
			fmt.Printf(",\n")
		}
		fmt.Printf("%s}", indent)
	case string:
		fmt.Printf("%q", t)
	case float64: // positions
		fmt.Printf("\"%v\"", t)
	default:
		fmt.Printf("%v", t)
	}
}

func keys(m map[string]interface{}) []string {
	var ks []string
	for k := range m {
		ks = append(ks, k)
	}

	if found := isNode(ks); found != nil {
		return found
	}

	sort.Strings(ks)
	return ks
}

var nodeKeys = []string{
	"InternalName",
	"InternalType",
	"Properties",
	"Children",
	"StartOffset",
	"EndOffset",
}

func isNode(a []string) []string {
	var seen []string

	for _, k := range nodeKeys {
		for _, l := range a {
			if k == l {
				seen = append(seen, k)
			}
		}
	}
	if len(seen) != len(a) {
		return nil
	}
	return seen
}
