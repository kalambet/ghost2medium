package ghost2medium

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	utils "github.com/kalambet/go-utils"
	medium "github.com/kalambet/medium-sdk-go"
)

type publication struct {
	Publication *medium.Publication
	Role        string
}

type destination struct {
	Medium      *medium.Medium
	User        *medium.User
	Publication *publication
}

// Import intercative CLI Publicaiton slector
func Import(token string, test bool, posts []*Post) (err error) {
	if len(posts) == 0 {
		return errors.New("there are no posts to import")
	}

	utils.PrintInColorln(fmt.Sprintf("\n\tNumber of posts in your archive is %d", len(posts)), utils.Green)
	utils.PrintInColorln("\nYou'll need it later for letter 📬  to Medium mailto:yourfriends@medium.com", utils.Yellow)
	utils.PrintInColorln("Please read: https://github.com/Medium/medium-api-docs/issues/28 \n", utils.Magenta)

	d := &destination{}
	err = d.selectPublication(token)

	if err != nil {
		return err
	}

	err = d.importArchive(posts, test)
	if err != nil {
		return err
	}

	return nil
}

func (d *destination) selectPublication(token string) (err error) {
	d.Medium = medium.NewClientWithAccessToken(token)

	d.User, err = d.Medium.GetUser()
	if err != nil {
		return err
	}

	publications, err := d.Medium.GetPublications(d.User.ID)
	if err != nil {
		return err
	}

	if len(*publications) == 0 {
		return errors.New("there are no publications where you can import your archive")
	}

	publicationsSelector := make([]*publication, 0)

	utils.PrintInColorln("Getting list of all Publications you can contribute to 📑 ", utils.Green)
	progressScale := 100 / len(*publications)
	progress := 0
	utils.PrintInColor("⏳ ", utils.White)
	for _, pub := range *publications {
		utils.PrintInColor(fmt.Sprintf("%2d%%...", progress), utils.Yellow)
		progress += progressScale
		contributors, err := d.Medium.GetContributors(pub.ID)
		if err != nil {
			return err
		}
		for _, c := range *contributors {
			if c.UserID == d.User.ID {
				publicationsSelector = append(publicationsSelector, &publication{Publication: pub, Role: c.Role})
			}
		}
	}
	utils.PrintInColorln("100%", utils.Yellow)

	uiIndex := 1 // Start from 1 for user conveniens
	utils.PrintInColorln("\nPlease select Publication you want import your archive to: ", utils.Magenta)
	for _, pub := range publicationsSelector {
		msg := fmt.Sprintf("\t%d: As %s to the %s (URL: %s, ID: %s)",
			uiIndex,
			strings.ToTitle(pub.Role),
			pub.Publication.Name,
			pub.Publication.URL,
			pub.Publication.ID)
		utils.PrintInColorln(msg, utils.Cyan)
		uiIndex++
	}

	utils.PrintInColor("\nPublication -> ", utils.Yellow)
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')

	pIndex, err := strconv.Atoi(strings.Trim(userInput, "\n"))
	if err != nil {
		return err
	} else if pIndex-1 >= len(publicationsSelector) {
		return errors.New("index you chose is not in range")
	}

	d.Publication = publicationsSelector[pIndex-1]
	return nil
}

func (d *destination) importArchive(posts []*Post, test bool) (err error) {
	for _, post := range posts {
		utils.PrintInColorln(fmt.Sprintf("\tImporting post '%s' (UUID: %s) ...", post.Title, post.UUID), utils.Yellow)
		utils.PrintInColorln(fmt.Sprintf("\tPublish Date: %s", post.Date.String()), utils.Cyan)
		if !test {
			_, err := d.Medium.CreatePost(medium.CreatePostOptions{
				UserID:          d.User.ID,
				Title:           post.Title,
				Content:         post.Markdown,
				ContentFormat:   medium.ContentFormatMarkdown,
				PublishStatus:   medium.PublishStatusDraft,
				PublishedAt:     post.CreatedAt,
				Tags:            post.Tags,
				PublicationID:   d.Publication.Publication.ID,
				NotifyFollowers: false,
			})
			if err != nil {
				utils.PrintInColorln(fmt.Sprintf("\tPost with id %s was not imported", post.UUID), utils.Red)
				utils.PrintInColorln(fmt.Sprintf("\t%s", err.Error()), utils.Red)
			}

		}
		utils.PrintInColorln("\t---\n", utils.Yellow)
	}
	return nil
}
