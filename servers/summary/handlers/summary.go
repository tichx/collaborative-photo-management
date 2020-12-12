package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const headerCORS = "Access-Control-Allow-Origin"
const corsAnyOrigin = "*"

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("summary handler")
	w.Header().Add(headerCORS, corsAnyOrigin)
	w.Header().Set("Content-Type", "application/json")
	name := r.FormValue("url")
	if len(name) == 0 {

		http.Error(w, "empty url", http.StatusBadRequest)
		return
	}
	resp, err := fetchHTML(name)

	if err != nil {
		errmsg := fmt.Sprintf("bad request: %v", err)
		http.Error(w, errmsg, http.StatusBadRequest)
		return
	}

	defer resp.Close()

	psum, err := extractSummary(name, resp)

	if err != nil {
		errmsg := fmt.Sprintf("internal error: %v", err)
		http.Error(w, errmsg, 500)
		return
	}
	if err := json.NewEncoder(w).Encode(psum); err != nil {
		errmsg := fmt.Sprintf("internal error: %v", err)
		http.Error(w, errmsg, 500)
		return
	}

}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {

	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New("not ok")
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, errors.New("bad content type")
	}

	return resp.Body, nil

}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	u, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}
	pageSummary := new(PageSummary)

	tokenizer := html.NewTokenizer(htmlStream)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return pageSummary, nil
			}
			return nil, tokenizer.Err()
		} else if tokenType == html.EndTagToken {
			if tokenizer.Token().Data == "head" {
				return pageSummary, nil
			}
		} else if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			switch token.Data {
			case "link":
				for _, attribute := range token.Attr {
					if attribute.Key == "rel" && attribute.Val == "icon" {
						img := new(PreviewImage)
						for _, attrKey := range token.Attr {
							if attrKey.Key == "href" {

								v, err := url.Parse(attrKey.Val)
								if err != nil {
									return nil, err
								}
								compURL := u.ResolveReference(v)
								img.URL = compURL.String()
							}
							if attrKey.Key == "type" {
								img.Type = attrKey.Val
							}
							if attrKey.Key == "sizes" {
								sizeArr := strings.Split(attrKey.Val, "x")
								if len(sizeArr) > 1 {
									img.Height, _ = strconv.Atoi(sizeArr[0])
									img.Width, _ = strconv.Atoi(sizeArr[1])
								}
							}

						}
						pageSummary.Icon = img
					}
				}
			case "title":
				tokenizer.Next()
				if len(pageSummary.Title) == 0 {
					pageSummary.Title = tokenizer.Token().Data
				}
			case "meta":
				for _, attr := range token.Attr {
					switch attr.Key {
					case "property":
						var content string
						for _, c := range token.Attr {
							if c.Key == "content" {
								content = c.Val
							}
						}
						if len(content) > 0 {
							lastIndex := len(pageSummary.Images) - 1

							switch attr.Val {
							case "og:url":
								v, err := url.Parse(content)
								if err != nil {
									return nil, err
								}
								compURL := u.ResolveReference(v)
								pageSummary.URL = compURL.String()
							case "og:site_name":
								pageSummary.SiteName = content
							case "og:type":
								pageSummary.Type = content
							case "og:title":
								pageSummary.Title = content
							case "og:description":
								pageSummary.Description = content
							case "og:image:url":
								lastSlice := pageSummary.Images[lastIndex]
								lastSlice.URL = content
							case "og:image:secure_url":
								lastSlice := pageSummary.Images[lastIndex]
								lastSlice.SecureURL = content
							case "og:image":
								image := PreviewImage{}
								v, err := url.Parse(content)
								if err != nil {
									return nil, err
								}
								compURL := u.ResolveReference(v)
								image.URL = compURL.String()
								pageSummary.Images = append(pageSummary.Images, &image)
							case "og:image:type":
								lastSlice := pageSummary.Images[lastIndex]
								lastSlice.Type = content
							case "og:image:width":
								lastSlice := pageSummary.Images[lastIndex]
								w, err := strconv.Atoi(content)
								if err != nil {
									return nil, err
								}
								lastSlice.Width = w
							case "og:image:height":
								lastSlice := pageSummary.Images[lastIndex]

								h, err := strconv.Atoi(content)
								if err != nil {
									return nil, err
								}
								lastSlice.Height = h
							case "og:image:alt":
								lastSlice := pageSummary.Images[lastIndex]
								lastSlice.Alt = content
							}
						}
					case "name":

						var content string
						for _, c := range token.Attr {
							if c.Key == "content" {
								content = c.Val
							}
						}
						if len(content) > 0 {
							switch attr.Val {
							case "keywords":
								word := strings.Split(content, ",")
								for _, eachWord := range word {
									eachWord = strings.TrimSpace(eachWord)
									pageSummary.Keywords = append(pageSummary.Keywords, eachWord)
								}
							case "description":
								if len(pageSummary.Description) == 0 {
									pageSummary.Description = content
								}
							case "author":
								pageSummary.Author = content
							}
						}

					}
				}

			}
		}

	}

}
