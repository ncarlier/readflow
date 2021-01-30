package oembed

import (
	"log"
	"regexp"
)

type entries struct {
	entries []*entry
}

type entry struct {
	Name     string
	Endpoint string
	re       *regexp.Regexp
}

func (db *entries) GetProviderEndpoint(rawurl string) string {
	for _, provider := range db.entries {
		if provider.re.MatchString(rawurl) {
			return provider.Endpoint
		}
	}
	return ""
}

var providers = &entries{}

func init() {
	for _, provider := range Providers {
		for _, endpoint := range provider.Endpoints {
			for _, scheme := range endpoint.Schemes {
				re, err := Scheme2Regexp(scheme)
				if err != nil {
					log.Fatal(err)
				}
				p := &entry{
					Name:     provider.Name,
					Endpoint: endpoint.URL,
					re:       re,
				}
				providers.entries = append(providers.entries, p)
			}
		}
	}
}
