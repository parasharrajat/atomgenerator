package atomgenerator

import (
	"bytes"
	"regexp"
	"testing"
	"time"
)

func TestFullXmlFeed(t *testing.T) {
	pubDate, _ := time.Parse("2006-01-02 15:04", "2008-09-10 11:12")

	f := Feed{
		Title:   "title",
		PubDate: pubDate,
		Link:    "http://www.myblog.bogus",
		Authors: []Author{
			Author{
				Name:  "author name",
				Email: "author email",
				Uri:   "author uri",
			},
		},
	}

	entryPubDate, _ := time.Parse("2006-01-02 15:04", "2009-10-11 12:13")
	f.AddEntry(&Entry{
		Title:       "entry title",
		PubDate:     entryPubDate,
		Link:        "http://www.myblog.bogus/entry",
		Description: "entry description",
		Content:     "<p>entry content</p>",
	})

	atom, err := f.GenXml()
	if err != nil {
		t.Error(err)
	}

	expected := []byte(`<?xml version="1.0" encoding="UTF-8"?> <feed xmlns="http://www.w3.org/2005/Atom">
		  <title>title</title>
		  <link href="http://www.myblog.bogus" rel="alternate"></link>
		  <id>http://www.myblog.bogus</id>
		  <updated>2008-09-10T11:12:00Z</updated>
		  <author>
		   <name>author name</name>
		   <email>author email</email>
		   <uri>author uri</uri>
		  </author>
		  <entry>
		   <title>entry title</title>
		   <link href="http://www.myblog.bogus/entry" rel="alternate"></link>
		   <updated>2009-10-11T12:13:00Z</updated>
		   <id>tag:www.myblog.bogus,2009-10-11:/entry</id>
		   <summary type="html">entry description</summary>
		   <content type="html">&lt;p&gt;entry content&lt;/p&gt;</content>
		  </entry>
		 </feed>`)

	whitespace := regexp.MustCompile(`\s+`)
	noWs := func(b []byte) []byte {
		return whitespace.ReplaceAll(atom, []byte(" "))
	}

	if !bytes.Equal(noWs(atom), noWs(expected)) {
		t.Errorf("XML differs: expected %s, got %s.\n", expected, atom)
	}
}

func TestValidation(t *testing.T) {
	now := time.Now()

	f := Feed{Title: "title"}
	if errs := f.Validate(); len(errs) != 1 {
		t.Error("Expected an error for a feed without PubDate.")
	}

	f = Feed{PubDate: now}
	if errs := f.Validate(); len(errs) != 1 {
		t.Error("Expected an error for a feed without Title.")
	}

	f = Feed{
		Title:   "title",
		PubDate: now,
	}
	if errs := f.Validate(); len(errs) != 0 {
		t.Error("Expected no errors for a feed with Title&PubDate and no entries.")
	}

	e := &Entry{Title: "entry title"}
	f.AddEntry(e)
	if errs := f.Validate(); len(errs) != 2 {
		t.Error("Expected two errors for lack of author and no entry PubDate.")
	}

	e.PubDate = now
	if errs := f.Validate(); len(errs) != 1 {
		t.Error("Expected an error for lack of author.")
	}

	f.Authors = append(f.Authors, Author{Name: "name"})
	if errs := f.Validate(); len(errs) != 0 {
		t.Error("Expected no errors for complete feed.")
	}

	f.Authors = append(f.Authors, Author{Email: "email"})
	if errs := f.Validate(); len(errs) != 1 {
		t.Error("Expected an error for entry Author without Name.")
	}
}
