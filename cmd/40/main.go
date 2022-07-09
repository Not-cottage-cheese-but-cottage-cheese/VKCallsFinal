package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"

	vk_api "github.com/SevereCloud/vksdk/v2/api"
)

func main() {
	if len(os.Args) < 4 {
		log.Panic("invalid argument count")
	}

	token, ownerID, videoID := os.Args[1], os.Args[2], os.Args[3]
	patterns := os.Args[4:]

	regexps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		if re, err := regexp.Compile(pattern); err != nil {
			log.Panicf("invalid pattern %s\n", pattern)
		} else {
			regexps[i] = re
		}
	}

	api := vk_api.NewVK(token)

	messages := []string{}

	offset := 0
	for {
		resp, err := api.VideoGetComments(vk_api.Params{
			"owner_id": ownerID,
			"video_id": videoID,
			"count":    100,
			"offset":   offset,
		})

		if err != nil {
			log.Panicln(err)
		}
		if len(resp.Items) == 0 {
			break
		}

		for _, commentInfo := range resp.Items {
			messages = append(messages, commentInfo.Text)
		}

		offset += 100
	}

	commentsByPattern := make(map[string]int, len(patterns))
	for _, pattern := range patterns {
		commentsByPattern[pattern] = 0
	}

	for _, message := range messages {
		for i, re := range regexps {
			if re.MatchString(message) {
				commentsByPattern[patterns[i]]++
			}
		}
	}

	type pair struct {
		pattern string
		count   int
	}

	patternsLeadboard := make([]pair, 0, len(patterns))
	for pattern, count := range commentsByPattern {
		patternsLeadboard = append(patternsLeadboard, pair{
			pattern: pattern,
			count:   count,
		})
	}
	sort.SliceStable(patternsLeadboard, func(i, j int) bool {
		return patternsLeadboard[i].count > patternsLeadboard[j].count
	})

	for i, p := range patternsLeadboard {
		fmt.Printf("%d) %s => %d\n", i+1, p.pattern, p.count)
	}
}
