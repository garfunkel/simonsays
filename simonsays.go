package main

import (
	"log"
	"encoding/json"
	"fmt"
	"regexp"
	"html"
	"strings"
	"math/rand"
	"io/ioutil"
	"time"
	"flag"
	"os"
	"net/http"
	"github.com/gcmurphy/getpass"
)

type Issue struct {
	Expand string `json:"expand"`
	ID string `json:"id"`
	Self string `json:"self"`
	Key string `json:"key"`
	Fields struct {
		Progress struct {
			Progress int `json:"progress"`
			Total int `json:"total"`
		} `json:"progress"`
		Summary string `json:"summary"`
		Timetracking struct {
		} `json:"timetracking"`
		Issuetype struct {
			Self string `json:"self"`
			ID string `json:"id"`
			Description string `json:"description"`
			IconURL string `json:"iconUrl"`
			Name string `json:"name"`
			Subtask bool `json:"subtask"`
		} `json:"issuetype"`
		Votes struct {
			Self string `json:"self"`
			Votes int `json:"votes"`
			HasVoted bool `json:"hasVoted"`
		} `json:"votes"`
		FixVersions []struct {
			Self string `json:"self"`
			ID string `json:"id"`
			Description string `json:"description"`
			Name string `json:"name"`
			Archived bool `json:"archived"`
			Released bool `json:"released"`
			ReleaseDate string `json:"releaseDate"`
		} `json:"fixVersions"`
		Resolution struct {
			Self string `json:"self"`
			ID string `json:"id"`
			Description string `json:"description"`
			Name string `json:"name"`
		} `json:"resolution"`
		Resolutiondate string `json:"resolutiondate"`
		Timespent int `json:"timespent"`
		Reporter struct {
			Self string `json:"self"`
			Name string `json:"name"`
			EmailAddress string `json:"emailAddress"`
			AvatarUrls struct {
				One6X16 string `json:"16x16"`
				Two4X24 string `json:"24x24"`
				Three2X32 string `json:"32x32"`
				Four8X48 string `json:"48x48"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active bool `json:"active"`
		} `json:"reporter"`
		Aggregatetimeoriginalestimate int `json:"aggregatetimeoriginalestimate"`
		Updated string `json:"updated"`
		Created string `json:"created"`
		Description string `json:"description"`
		Priority struct {
			Self string `json:"self"`
			IconURL string `json:"iconUrl"`
			Name string `json:"name"`
			ID string `json:"id"`
		} `json:"priority"`
		Duedate string `json:"duedate"`
		Issuelinks []struct {
			ID string `json:"id"`
			Self string `json:"self"`
			Type struct {
				ID string `json:"id"`
				Name string `json:"name"`
				Inward string `json:"inward"`
				Outward string `json:"outward"`
				Self string `json:"self"`
			} `json:"type"`
			InwardIssue struct {
				ID string `json:"id"`
				Key string `json:"key"`
				Self string `json:"self"`
				Fields struct {
					Summary string `json:"summary"`
					Status struct {
						Self string `json:"self"`
						Description string `json:"description"`
						IconURL string `json:"iconUrl"`
						Name string `json:"name"`
						ID string `json:"id"`
					} `json:"status"`
					Priority struct {
						Self string `json:"self"`
						IconURL string `json:"iconUrl"`
						Name string `json:"name"`
						ID string `json:"id"`
					} `json:"priority"`
					Issuetype struct {
						Self string `json:"self"`
						ID string `json:"id"`
						Description string `json:"description"`
						IconURL string `json:"iconUrl"`
						Name string `json:"name"`
						Subtask bool `json:"subtask"`
					} `json:"issuetype"`
				} `json:"fields"`
			} `json:"inwardIssue"`
		} `json:"issuelinks"`
		Watches struct {
			Self string `json:"self"`
			WatchCount int `json:"watchCount"`
			IsWatching bool `json:"isWatching"`
		} `json:"watches"`
		Worklog struct {
			StartAt int `json:"startAt"`
			MaxResults int `json:"maxResults"`
			Total int `json:"total"`
			Worklogs []struct {
				Self string `json:"self"`
				Author struct {
					Self string `json:"self"`
					Name string `json:"name"`
					EmailAddress string `json:"emailAddress"`
					AvatarUrls struct {
						One6X16 string `json:"16x16"`
						Two4X24 string `json:"24x24"`
						Three2X32 string `json:"32x32"`
						Four8X48 string `json:"48x48"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active bool `json:"active"`
				} `json:"author"`
				UpdateAuthor struct {
					Self string `json:"self"`
					Name string `json:"name"`
					EmailAddress string `json:"emailAddress"`
					AvatarUrls struct {
						One6X16 string `json:"16x16"`
						Two4X24 string `json:"24x24"`
						Three2X32 string `json:"32x32"`
						Four8X48 string `json:"48x48"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active bool `json:"active"`
				} `json:"updateAuthor"`
				Comment string `json:"comment"`
				Created string `json:"created"`
				Updated string `json:"updated"`
				Started string `json:"started"`
				TimeSpent string `json:"timeSpent"`
				TimeSpentSeconds int `json:"timeSpentSeconds"`
				ID string `json:"id"`
			} `json:"worklogs"`
		} `json:"worklog"`
		Subtasks []struct {
			ID string `json:"id"`
			Key string `json:"key"`
			Self string `json:"self"`
			Fields struct {
				Summary string `json:"summary"`
				Status struct {
					Self string `json:"self"`
					Description string `json:"description"`
					IconURL string `json:"iconUrl"`
					Name string `json:"name"`
					ID string `json:"id"`
				} `json:"status"`
				Priority struct {
					Self string `json:"self"`
					IconURL string `json:"iconUrl"`
					Name string `json:"name"`
					ID string `json:"id"`
				} `json:"priority"`
				Issuetype struct {
					Self string `json:"self"`
					ID string `json:"id"`
					Description string `json:"description"`
					IconURL string `json:"iconUrl"`
					Name string `json:"name"`
					Subtask bool `json:"subtask"`
				} `json:"issuetype"`
			} `json:"fields"`
		} `json:"subtasks"`
		Status struct {
			Self string `json:"self"`
			Description string `json:"description"`
			IconURL string `json:"iconUrl"`
			Name string `json:"name"`
			ID string `json:"id"`
		} `json:"status"`
		Labels []string `json:"labels"`
		Workratio int `json:"workratio"`
		Assignee struct {
			Self string `json:"self"`
			Name string `json:"name"`
			EmailAddress string `json:"emailAddress"`
			AvatarUrls struct {
				One6X16 string `json:"16x16"`
				Two4X24 string `json:"24x24"`
				Three2X32 string `json:"32x32"`
				Four8X48 string `json:"48x48"`
			} `json:"avatarUrls"`
			DisplayName string `json:"displayName"`
			Active bool `json:"active"`
		} `json:"assignee"`
		Attachment []struct {
			Self string `json:"self"`
			ID string `json:"id"`
			Filename string `json:"filename"`
			Author struct {
				Self string `json:"self"`
				Name string `json:"name"`
				EmailAddress string `json:"emailAddress"`
				AvatarUrls struct {
					One6X16 string `json:"16x16"`
					Two4X24 string `json:"24x24"`
					Three2X32 string `json:"32x32"`
					Four8X48 string `json:"48x48"`
				} `json:"avatarUrls"`
				DisplayName string `json:"displayName"`
				Active bool `json:"active"`
			} `json:"author"`
			Created string `json:"created"`
			Size int `json:"size"`
			MimeType string `json:"mimeType"`
			Content string `json:"content"`
		} `json:"attachment"`
		Aggregatetimeestimate int `json:"aggregatetimeestimate"`
		Project struct {
			Self string `json:"self"`
			ID string `json:"id"`
			Key string `json:"key"`
			Name string `json:"name"`
			AvatarUrls struct {
				One6X16 string `json:"16x16"`
				Two4X24 string `json:"24x24"`
				Three2X32 string `json:"32x32"`
				Four8X48 string `json:"48x48"`
			} `json:"avatarUrls"`
			ProjectCategory struct {
				Self string `json:"self"`
				ID string `json:"id"`
				Description string `json:"description"`
				Name string `json:"name"`
			} `json:"projectCategory"`
		} `json:"project"`
		Environment string `json:"environment"`
		Timeestimate int `json:"timeestimate"`
		LastViewed string `json:"lastViewed"`
		Aggregateprogress struct {
			Progress int `json:"progress"`
			Total int `json:"total"`
		} `json:"aggregateprogress"`
		Components []struct {
			Self string `json:"self"`
			ID string `json:"id"`
			Name string `json:"name"`
			Description string `json:"description"`
		} `json:"components"`
		Comment struct {
			StartAt int `json:"startAt"`
			MaxResults int `json:"maxResults"`
			Total int `json:"total"`
			Comments []struct {
				Self string `json:"self"`
				ID string `json:"id"`
				Author struct {
					Self string `json:"self"`
					Name string `json:"name"`
					EmailAddress string `json:"emailAddress"`
					AvatarUrls struct {
						One6X16 string `json:"16x16"`
						Two4X24 string `json:"24x24"`
						Three2X32 string `json:"32x32"`
						Four8X48 string `json:"48x48"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active bool `json:"active"`
				} `json:"author"`
				Body string `json:"body"`
				UpdateAuthor struct {
					Self string `json:"self"`
					Name string `json:"name"`
					EmailAddress string `json:"emailAddress"`
					AvatarUrls struct {
						One6X16 string `json:"16x16"`
						Two4X24 string `json:"24x24"`
						Three2X32 string `json:"32x32"`
						Four8X48 string `json:"48x48"`
					} `json:"avatarUrls"`
					DisplayName string `json:"displayName"`
					Active bool `json:"active"`
				} `json:"updateAuthor"`
				Created string `json:"created"`
				Updated string `json:"updated"`
			} `json:"comments"`
		} `json:"comment"`
		Timeoriginalestimate int `json:"timeoriginalestimate"`
		Aggregatetimespent int `json:"aggregatetimespent"`
	} `json:"fields"`
}

type Issues struct {
	Expand string `json:"expand"`
	StartAt int `json:"startAt"`
	MaxResults int `json:"MaxResults"`
	Total int `json:"total"`
	Issues []Issue `json:"issues"`
}

type MarkovChainPrefix []string

type MarkovChain struct {
	chain map[string]map[string][]string
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

func (markovChain *MarkovChain) addMessage(username, msg string) bool {
	word := ""
	key := ""
	msgReader := strings.NewReader(msg)
	prefix := make(MarkovChainPrefix, markovChain.prefixLen)

	for {
		if _, err := fmt.Fscan(msgReader, &word); err != nil {
			break
		}

		key = prefix.String()

		_, ok := markovChain.chain[username]

		if !ok {
			markovChain.chain[username] = make(map[string][]string)
		}

		markovChain.chain[username][key] = append(markovChain.chain[username][key], word)
		prefix.Shift(word)
	}

	return true
}

func (markovChain *MarkovChain) generateMessage(username string) (msg string) {
	sentence := []string{}

	for {
		word := ""
		key := ""
		sentence = []string{}
		prefix := make(MarkovChainPrefix, markovChain.prefixLen)
		badSentence := false

		for {
			key = prefix.String()

			_, ok := markovChain.chain[username]

			if !ok {
				msg = ""

				return
			}

			words := markovChain.chain[username][key]

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

func NewMarkovChain(semanticAccuracy uint) *MarkovChain {
	return &MarkovChain{
		chain: make(map[string]map[string][]string),
		prefixLen: semanticAccuracy,
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

func refreshCache(username string) (issues Issues, err error) {
	var client http.Client
	var startAt int
	var newIssues Issues
	var request *http.Request
	var response *http.Response
	var data []byte
	var urlTemplate = "https://jira.appen.com/rest/api/2/search?fields=*all&maxResults=1000&startAt=%v"

	fmt.Printf("JIRA ")

	password, err := getpass.GetPass()

	if err != nil {
		return
	}

	for {
		request, err = http.NewRequest("GET", fmt.Sprintf(urlTemplate, startAt), nil)

		if err != nil {
			return
		}

		request.SetBasicAuth(username, string(password))

		response, err = client.Do(request)

		if err != nil {
			return
		}

		data, err = ioutil.ReadAll(response.Body)

		if err != nil {
			return
		}

		if err = json.Unmarshal(data, &newIssues); err != nil {
			return
		}

		if len(newIssues.Issues) == 0 {
			break
		}

		issues.Issues = append(issues.Issues, newIssues.Issues...)

		fmt.Printf("%v issues cached\n", len(issues.Issues))

		startAt += len(newIssues.Issues)
	}

	issues.MaxResults = len(issues.Issues)
	issues.Total = len(issues.Issues)

	data, err = json.Marshal(issues)

	if err != nil {
		return
	}

	// Save cache.
	if err = ioutil.WriteFile("/tmp/simonsays_cache.json", data, 0664); err != nil {
		return
	}

	return
}

func loadCache() (issues Issues, err error) {
	if _, err = os.Stat("/tmp/simonsays_cache.json"); err != nil {
		return
	}

	data, err := ioutil.ReadFile("/tmp/simonsays_cache.json")

	err = json.Unmarshal(data, &issues)

	return
}

func main() {
	username := flag.String("u", "", "username to generate messages for")
	wait := flag.Bool("w", false, "wait for a little while after generating each message")
	semanticAccuracy := flag.Uint("a", 2, "semantic correctness target (higher = more grammatical message but less creative)")
	refresh := flag.Bool("r", false, "refresh locally cached JIRA issue database which is used by simonsays")
	refreshUsername := flag.String("j", "", "JIRA username to use when refreshing local issue cache")
	server := flag.Bool("s", false, "start web server to service HTTP clients")

	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	cleaner := textCleaner()

	var issues Issues
	var err error

	if *refresh {
		issues, err = refreshCache(*refreshUsername)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		issues, err = loadCache()

		if os.IsNotExist(err) {
			issues, err = refreshCache(*refreshUsername)

			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Fatal(err)
		}
	}

	if *server {
		markovChain := NewMarkovChain(*semanticAccuracy)

		for _, issue := range issues.Issues {
			markovChain.addMessage(issue.Fields.Reporter.Name, cleaner(issue.Fields.Description))
			markovChain.addMessage("", cleaner(issue.Fields.Description))

			for _, comment := range issue.Fields.Comment.Comments {
				markovChain.addMessage(comment.Author.Name, cleaner(comment.Body))
				markovChain.addMessage("", cleaner(comment.Body))
			}
		}

		http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
			username := strings.TrimLeft(r.URL.Path, "/")

			if username == "favicon.ico" {
				return
			}

			fmt.Printf("Request from %v for user %v\n", r.RemoteAddr, username)

			for i := 0; i < 1000; i++ {
				msg := markovChain.generateMessage(username)

				if msg == "" {
					fmt.Fprintf(w,
								"<html><body><marquee scrollamount=\"90\">" +
					            "<span style=\"background: red; font-size: 300px; color: white\">WHOOP! WHOOP! ... USER DOES NOT EXIST IN JIRA! ... WHOOP! WHOOP!</span>" +
					            "</marquee></body></html>")
					return
				}

				fmt.Fprintf(w, "%s\n\n", msg)
			}
		})

		fmt.Println("Listening on http://0.0.0.0:49001/")

		log.Fatal(http.ListenAndServe("0.0.0.0:49001", nil))
	} else {
		markovChain := NewMarkovChain(*semanticAccuracy)

		for _, issue := range issues.Issues {
			if *username == "" || issue.Fields.Reporter.Name == *username {
				markovChain.addMessage(*username, cleaner(issue.Fields.Description))
			}

			for _, comment := range issue.Fields.Comment.Comments {
				if *username == "" || comment.Author.Name == *username {
					markovChain.addMessage(*username, cleaner(comment.Body))
				}
			}
		}

		for {
			msg := markovChain.generateMessage(*username)

			if msg == "" {
				fmt.Println("User does not exist in JIRA")

				break
			}

			fmt.Printf("%s\n\n", msg)

			if *wait {
				time.Sleep(time.Duration(len(msg) * 45) * time.Millisecond)
			}
		}
	}
}
