package ghost2medium

import (
	"strings"
	"testing"
	"time"
)

func TestDecodeJSONArchive(t *testing.T) {
	posts, err := DecodeJSONArchive("")
	if err == nil {
		t.Error("Decoder call with empty file should throw an error")
	}

	if len(posts) != 0 {
		t.Error("Empty Decoder call should return empty array!")
	}

	posts, err = DecodeJSONArchive("./testdata/blog_test.json")

	if err != nil {
		t.Error(err.Error())
	}

	if n := len(posts); n != 3 {
		t.Errorf("Number of posts needs to be 3 but actual is %d", n)
	}

	post := posts[0]

	if strings.Compare(post.Title, "First Post") != 0 {
		t.Errorf("First post should have has title `Frist Post` but instead its title `%s`", post.Title)
	}

	if n := len(post.Tags); n != 2 {
		t.Errorf("`First Post` should have has two Tags but instead it has %d", n)
	}

	if strings.Compare(post.Tags[0], "first tag") != 0 && strings.Compare(post.Tags[1], "first tag") != 0 {
		t.Error("One of the tag should be `frist tag`")
	}

	if strings.Compare(post.Tags[0], "third tag") != 0 && strings.Compare(post.Tags[1], "third tag") != 0 {
		t.Error("One of the tag should be `third tag`")
	}

	if post.Date.Year() != 2011 ||
		post.Date.Month() != time.January ||
		post.Date.Day() != 1 ||
		post.Date.Hour() != 1 ||
		post.Date.Minute() != 1 {
		t.Errorf("First post should have publish date 2011-01-01T01:01:00.000Z and not %s", post.Date.String())
	}

	post = posts[1]

	if strings.Compare(post.Title, "Second Post") != 0 {
		t.Errorf("Second post should have has title `Second Post` but instead its title `%s`", post.Title)
	}

	if n := len(post.Tags); n != 0 {
		t.Errorf("`Second Post` should not have has Tags but instead it has %d", n)
	}

	if post.Date.Year() != 2012 ||
		post.Date.Month() != time.February ||
		post.Date.Day() != 2 ||
		post.Date.Hour() != 2 ||
		post.Date.Minute() != 2 {
		t.Errorf("Second post should have publish date 2012-02-02T02:02:00.000Z and not %s", post.Date.String())
	}

	post = posts[2]

	if strings.Compare(post.Title, "Third Post") != 0 {
		t.Errorf("Thid post should have has title `Third Post` but instead its title `%s`", post.Title)
	}

	if n := len(post.Tags); n != 1 {
		t.Errorf("`Third Post` should have has one Tag but instead it has %d", n)
	}

	if strings.Compare(post.Tags[0], "second tag") != 0 {
		t.Error("The tag should be `second tag`")
	}

	if post.Date.Year() != 2013 ||
		post.Date.Month() != time.March ||
		post.Date.Day() != 3 ||
		post.Date.Hour() != 3 ||
		post.Date.Minute() != 3 {
		t.Errorf("Third post should have publish date 2013-03-03T03:03:00.000Z and not %s", post.Date.String())
	}
}
