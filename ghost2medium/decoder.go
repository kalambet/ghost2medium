package ghost2medium

import (
	"encoding/json"
	"os"
	"sort"
	"time"
)

// DecodeJSONArchive decode provided JSON to array of Post and sort them by Date
func DecodeJSONArchive(path string) (posts []*Post, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)

	archive := new(Archive)
	if err = decoder.Decode(archive); err != nil {
		return nil, err
	}

	tagsPostsMap := make(TagsPostMap)
	for _, pt := range archive.DB[0].Data.PostTags {
		tagsPostsMap[pt.PostID] = append(tagsPostsMap[pt.PostID], pt.TagID)
	}

	tags := make(Tags)
	for _, t := range archive.DB[0].Data.Tags {
		tags[t.ID] = *t
	}

	// Mon Jan 2 15:04:05 -0700 MST 2006
	const iso8601Layout = "2006-01-02T15:04:05.000Z"
	for _, post := range archive.DB[0].Data.Posts {
		post.Date, err = time.Parse(iso8601Layout, post.PublishedAt)
		if err != nil {
			return nil, err
		}

		for _, tagID := range tagsPostsMap[post.ID] {
			post.Tags = append(post.Tags, tags[tagID].Name)
		}
	}

	sort.Sort(ByDate(archive.DB[0].Data.Posts))

	return archive.DB[0].Data.Posts, nil
}
