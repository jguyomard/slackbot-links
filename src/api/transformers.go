package api

import "../links"

func linkTransformer(item interface{}) interface{} {
	link, isLink := item.(*links.Link)
	if !isLink {
		return nil
	}

	return map[string]interface{}{
		"url":            link.URL,
		"title":          link.Title,
		"excerpt":        stringOrNil(link.Excerpt),
		"author":         stringOrNil(link.Author),
		"date_published": dateOrNil(link.DatePublished),
		"image_url":      stringOrNil(link.ImageURL),
	}
}
