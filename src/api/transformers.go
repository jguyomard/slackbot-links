package api

import (
	"fmt"

	"../links"
)

func linkTransformer(item interface{}) interface{} {
	link, isLink := item.(*links.Link)
	if !isLink {
		return nil
	}

	return map[string]interface{}{
		"id":           link.ID,
		"url":          link.URL,
		"title":        link.Title,
		"excerpt":      stringOrNil(link.Excerpt),
		"author":       stringOrNil(link.Author),
		"published_at": dateOrNil(link.DatePublished),
		"image_url":    stringOrNil(link.ImageURL),

		"shared_at": dateOrNil(link.SharedAt),
		"shared_by": map[string]string{
			"id":   link.SharedBy.ID,
			"name": link.SharedBy.Name,
		},
		"shared_on": map[string]string{
			"id":   link.SharedOn.ID,
			"name": link.SharedOn.Name,
		},

		"links": map[string]string{
			"self": fmt.Sprintf("/links/%s", link.ID),
		},
	}
}
