package ghost2medium

import (
	"encoding/json"
	"os"
)

type Archive struct {
	DB []*Blog `json:"db"`
}

type Blog struct {
	MetaRaw json.RawMessage `json:"meta"`
	Data    *BlogData       `json:"data"`
}

type BlogData struct {
	StatementsRaw       json.RawMessage `json:"pg_stat_statements"`
	GeographyClmnsRaw   json.RawMessage `json:"geography_columns"`
	GeometryClmnsRaw    json.RawMessage `json:"geometry_columns"`
	SpatialRefSysRaw    json.RawMessage `json:"spatial_ref_sys"`
	RasterClmnsRaw      json.RawMessage `json:"raster_columns"`
	RasterOverviewsRaw  json.RawMessage `json:"raster_overviews"`
	RolesRaw            json.RawMessage `json:"roles"`
	PermissionsRaw      json.RawMessage `json:"permissions"`
	PermissionsUsersRaw json.RawMessage `json:"permissions_users"`
	PermissionsAppsRaw  json.RawMessage `json:"permissions_apps"`
	SettingsRaw         json.RawMessage `json:"settings"`
	PostTags            []*PostTags     `json:"posts_tags"`
	Posts               []*Post         `json:"posts"`
	RolesUsersRaw       json.RawMessage `json:"roles_users"`
	PermissionsRolesRaw json.RawMessage `json:"permissions_roles"`
	UsersRaw            json.RawMessage `json:"users"`
	AppSettingsRaw      json.RawMessage `json:"app_settings"`
	Tags                []*Tag          `json:"tags"`
	AppsRaw             json.RawMessage `json:"apps"`
	AppFieldsRaw        json.RawMessage `json:"app_fileds"`
	Subscribers         json.RawMessage `json:"subscribers"`
}

type Post struct {
	ID              int             `json:"id"`
	UUID            string          `json:"uuid"`
	Title           string          `json:"title"`
	Slug            string          `json:"slug"`
	Markdown        string          `json:"markdown"`
	HTML            string          `json:"html"`
	Image           json.RawMessage `json:"image"`
	Featured        bool            `json:"featured"`
	Page            bool            `json:"page"`
	Status          string          `json:"status"`
	Language        string          `json:"language"`
	MetaTitle       json.RawMessage `json:"meta_title"`
	MetaDescription json.RawMessage `json:"meta_description"`
	AuthorID        int             `json:"author_id"`
	CreatedAt       string          `json:"created_at"`
	CreatedBy       int             `json:"created_by"`
	UpdatedAt       string          `json:"updated_at"`
	UpdatedBy       int             `json:"updated_by"`
	PublishedAt     string          `json:"published_at"`
	PublishedBy     int             `json:"published_by"`
	Visibility      string          `json:"visibility"`
	MobileDoc       json.RawMessage `json:"mobiledoc"`
	Amp             json.RawMessage `json:"amp"`

	Tags []string
}

type Tag struct {
	ID              int             `json:"id"`
	UUID            string          `json:"uuid"`
	Name            string          `json:"name"`
	Slug            string          `json:"slug"`
	Description     string          `json:"description"`
	Image           json.RawMessage `json:"image"`
	ParentID        int             `json:"parent_id"`
	MetaTitle       json.RawMessage `json:"meta_title"`
	MetaDescription json.RawMessage `json:"meta_description"`
	CreatedAt       string          `json:"created_at"`
	CreatedBy       int             `json:"created_by"`
	UpdatedAt       string          `json:"updated_at"`
	UpdatedBy       int             `json:"updated_by"`
	Visibility      string          `json:"visibility"`
}

type PostTags struct {
	ID        int `json:"id"`
	PostID    int `json:"post_id"`
	TagID     int `json:"tag_id"`
	SortOrder int `json:"sort_order"`
}

type TagsPostMap map[int][]int
type Tags map[int]Tag

// DecodeJSONArchive decode provided JSON to array of Post
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

	for idx := range archive.DB[0].Data.Posts {
		for _, tagID := range tagsPostsMap[archive.DB[0].Data.Posts[idx].ID] {
			archive.DB[0].Data.Posts[idx].Tags = append(archive.DB[0].Data.Posts[idx].Tags, tags[tagID].Name)
		}
	}

	return archive.DB[0].Data.Posts, nil
}
