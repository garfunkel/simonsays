package main

import (
	"log"
	"encoding/xml"
	"fmt"
	"regexp"
	"html"
	"strings"
	"math/rand"
	"time"
	"flag"
)

type Jira struct {
	XmlName xml.Name `xml:"res"`
	Version float32 `xml:"version,attr"`
	Channels []struct {
		XmlName xml.Name `xml:"channel"`
		Title string `xml:"title"`
		Link string `xml:"link"`
		Description string `xml:"description"`
		Language string `xml:"language"`
		Issue struct {
			XmlName xml.Name `xml:"issue"`
			Start int `xml:"start,attr"`
			End int `xml:"end,attr"`
			Total int `xml:"total,attr"`
		} `xml:"issue"`
		BuildInfo struct {
			XmlName xml.Name `xml:"build-info"`
			Version string `xml:"version"`
			BuildNumber int `xml:"build-number"`
			BuildDate string `xml:"build-date"`
		} `xml:"build-info"`
		Items []struct {
			XmlName xml.Name `xml:"item"`
			Title string `xml:"title"`
			Link string `xml:"link"`
			Project struct {
				XmlName xml.Name `xml:"project"`
				Project string `xml:",innerxml"`
				Id int `xml:"id,attr"`
				Key string `xml:"key,attr"`
			} `xml:"project"`
			Description string `xml:"description"`
			Environment string `xml:"environment"`
			Key struct {
				XmlName xml.Name `xml:"key"`
				Id int `xml:"id,attr"`
				Key string `xml:",innerxml"`
			} `xml:"key"`
			Summary string `xml:"summary"`
			Type struct {
				XmlName xml.Name `xml:"type"`
				Id int `xml:"id,attr"`
				IconUrl string `xml:"iconUrl,attr"`
				Type string `xml:",innerxml"`
			} `xml:"type"`
			Priority struct {
				XmlName xml.Name `xml:"priority"`
				Id int `xml:"id,attr"`
				IconUrl string `xml:"iconUrl,attr"`
				Priority string `xml:",innerxml"`
			} `xml:"priority"`
			Status struct {
				XmlName xml.Name `xml:"status"`
				Id int `xml:"id,attr"`
				IconUrl string `xml:"iconUrl,attr"`
				Status string `xml:",innerxml"`
			} `xml:"status"`
			Resolution struct {
				XmlName xml.Name `xml:"resolution"`
				Id int `xml:"id,attr"`
				Resolution string `xml:",innerxml"`
			} `xml:"resolution"`
			Assignee struct {
				XmlName xml.Name `xml:"assignee"`
				Username string `xml:"username,attr"`
				Assignee string `xml:",innerxml"`
			} `xml:"assignee"`
			Reporter struct {
				XmlName xml.Name `xml:"reporter"`
				Username string `xml:"username,attr"`
				Reporter string `xml:",innerxml"`
			} `xml:"reporter"`
			Labels struct {
				XmlName xml.Name `xml:"labels"`
				Labels []string `xml:"label"`
			} `xml:"labels"`
			Created string `xml:"created"`
			Updated string `xml:"updated"`
			Resolved string `xml:"resolved"`
			Due string `xml:"due"`
			Votes int `xml:"votes"`
			Watches int `xml:"watches"`
			Comments []struct {
				XmlName xml.Name `xml:"comment"`
				Id int `xml:"id,attr"`
				Author string `xml:"author,attr"`
				Created string `xml:"created,attr"`
				Comment string `xml:",innerxml"`
			} `xml:"comments>comment"`
			Attachments []struct {
				XmlName xml.Name `xml:"attachment"`
				Id int `xml:"id,attr"`
				Name string `xml:"name,attr"`
				Size int `xml:"size,attr"`
				Author string `xml:"author,attr"`
				Created string `xml:"created"`
			} `xml:"attachments>attachment"`
			Subtasks []struct {
				XmlName xml.Name `xml:"subtask"`
				Id int `xml:"id,attr"`
				Subtask string `xml:",innerxml"`
			} `xml:"subtasks>subtask"`
			CustomFields []struct {
				XmlName xml.Name `xml:"customfield"`
				Id int `xml:"id,attr"`
				Key string `xml:"key,attr"`
				CustomFieldName string `xml:"customfieldname"`
				CustomFieldValues []string `xml:"customfieldvalues>customfieldvalue"`
			} `xml:"customfields>customfield"`
		} `xml:"item"`
	} `xml:"channel"`
}

type MarkovChainPrefix []string

type MarkovChain struct {
	chain map[string][]string
	prefixLen uint
	rudeWords map[string]struct{}
	minRudeWords uint
}

func (prefix MarkovChainPrefix) String() string {
	return strings.TrimSpace(strings.Join(prefix, " "))
}

func (prefix MarkovChainPrefix) Shift(word string) {
	copy(prefix, prefix[1 :])
	prefix[len(prefix) - 1] = word
}

func (markovChain *MarkovChain) addMessage(msg string) bool {
	word := ""
	key := ""
	msgReader := strings.NewReader(msg)
	prefix := make(MarkovChainPrefix, markovChain.prefixLen)

	for {
		if _, err := fmt.Fscan(msgReader, &word); err != nil {
			break
		}

		key = prefix.String()
		markovChain.chain[key] = append(markovChain.chain[key], word)
		prefix.Shift(word)
	}

	return true
}

func (markovChain *MarkovChain) generateMessage() (msg string) {
	sentence := []string{}

	for {
		word := ""
		key := ""
		sentence = []string{}
		prefix := make(MarkovChainPrefix, markovChain.prefixLen)
		badSentence := false

		for {
			key = prefix.String()
			words := markovChain.chain[key]

			if len(words) == 0 {
				break
			}

			word = words[rand.Intn(len(words))]
			sentence = append(sentence, word)

			if strings.HasSuffix(word, ".") {
				break
			}

			prefix.Shift(word)

			if len(sentence) > 1000 {
				badSentence = true
				break
			}
		}

		rudeWordsFound := uint(0)

		for _, word := range sentence {
			if _, ok := markovChain.rudeWords[strings.ToLower(word)]; ok {
				rudeWordsFound += 1
			}
		}

		if rudeWordsFound < markovChain.minRudeWords {
			badSentence = true
		}

		if !badSentence {
			break
		}
	}

	return strings.Join(sentence, " ")
}

func NewMarkovChain(semanticAccuracy uint, rudeWords []string, minRudeWords uint) *MarkovChain {
	rudeWordSet := make(map[string]struct{})

	for _, word := range rudeWords {
		rudeWordSet[word] = struct{}{}
	}

	return &MarkovChain{
		chain: make(map[string][]string),
		prefixLen: semanticAccuracy,
		rudeWords: rudeWordSet,
		minRudeWords: minRudeWords,
	}
}

func textCleaner() func(string) string {
	stripper, err := regexp.Compile("</?(div|p|img|a|li|ul|br|b|i|hr|span|ins|h\\d|tt|font|em|ol|table|tr|td|tbody|thead|th|del).*?>")

	if err != nil {
		log.Fatal(err)
	}

	cleaner, err := regexp.Compile("\\s+")

	if err != nil {
		log.Fatal(err)
	}

	return func(text string) string {
		return html.UnescapeString(strings.TrimSpace(cleaner.ReplaceAllString(stripper.ReplaceAllString(html.UnescapeString(text), ""), " ")))
	}
}

func main() {
	username := flag.String("u", "", "username to generate messages for")
	wait := flag.Bool("w", false, "wait for a little while after generating each message")
	semanticAccuracy := flag.Uint("a", 2, "semantic correctness target (higher = more grammatical message but less creative)")
	minRudeWords := flag.Uint("r", 0, "minimum number of rude words to appear in generated messages")

	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	cleaner := textCleaner()
	jira := Jira{}
	rudeWordsBytes, err := Asset("input/rudewords.txt")

	if err != nil {
		log.Fatal(err)
	}

	rudeWords := strings.Split(strings.ToLower(string(rudeWordsBytes)), "\n")
	markovChain := NewMarkovChain(*semanticAccuracy, rudeWords, *minRudeWords)

	for binKey, binValue := range _bindata {
		if !strings.HasPrefix(binKey, "input/xml/") {
			continue
		}

		bytes, err := binValue()

		if err != nil {
			log.Fatal(err)
		}

		if err = xml.Unmarshal(bytes, &jira); err != nil {
			log.Fatal(err)
		}

		for _, channel := range jira.Channels {
			for _, item := range channel.Items {
				if *username == "" || item.Reporter.Username == *username {
					markovChain.addMessage(cleaner(item.Summary))
					markovChain.addMessage(cleaner(item.Description))
				}

				for _, comment := range item.Comments {
					if *username == "" || comment.Author == *username {
						markovChain.addMessage(cleaner(comment.Comment))
					}
				}
			}
		}
	}

	for {
		msg := markovChain.generateMessage()

		fmt.Printf("%s\n\n", msg)

		if *wait {
			time.Sleep(time.Duration(len(msg) * 45) * time.Millisecond)
		}
	}
}
