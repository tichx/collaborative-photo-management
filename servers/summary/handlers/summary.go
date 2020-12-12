package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const headerCORS = "Access-Control-Allow-Origin"
const corsAnyOrigin = "*"
const contentTypeJSON = "application/json"
const headerContentType = "Content-Type"
const authorization = "Authorization"

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
	w.Header().Add(headerCORS, corsAnyOrigin)
	name := r.URL.Query().Get("url")
	if len(name) == 0 {
		http.Error(w, "URL needed:", http.StatusBadRequest)
		return
	}
	resp, err := fetchHTML(name)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	ps, err := extractSummary(name, resp)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer resp.Close()
	w.Header().Set(headerContentType, contentTypeJSON)
	psJSON, err := json.Marshal(ps)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(psJSON)
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.
	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML
	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, errors.New("unable to get the page")
	}

	if resp.StatusCode >= 400 {
		return nil, err
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, err
	}
	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.
	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary
	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/
	ps := &PageSummary{}
	tokenizer := html.NewTokenizer(htmlStream)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			return ps, tokenizer.Err()
		}
		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if "head" == token.Data {
				tokenType = tokenizer.Next()
				break
			}
		}

		token := tokenizer.Token()
		attrMap := make(map[string]string)
		for i := 0; i < len(token.Attr); i++ {
			attrMap[token.Attr[i].Key] = token.Attr[i].Val
		}

		if "title" == token.Data {
			tokenType = tokenizer.Next()
			if len(ps.Title) == 0 {
				ps.Title = tokenizer.Token().Data
			}

		}
		// handles icon
		if "link" == token.Data {
			if attrMap["rel"] == "icon" {
				pi := &PreviewImage{}
				if val, ok := attrMap["type"]; ok {
					pi.Type = val
				}
				if val, ok := attrMap["alt"]; ok {
					pi.Alt = val
				}
				if val, ok := attrMap["sizes"]; ok {
					// split sizes such as 30x30 into heights and width
					sizeArr := strings.Split(val, "x")
					if len(sizeArr) == 2 {
						height, err := strconv.Atoi(sizeArr[0])
						if err != nil {
							fmt.Errorf("parsing error: %v", err)
						}
						width, err := strconv.Atoi(sizeArr[1])
						if err != nil {
							fmt.Errorf("parsing error: %v", err)
						}
						pi.Height = height
						pi.Width = width
					} else {
						height, err := strconv.Atoi(sizeArr[0])
						if err != nil {
							fmt.Errorf("parsing error: %v", err)
						}
						pi.Height = height
					}

				}
				if val, ok := attrMap["href"]; ok {
					if strings.HasPrefix(val, "https") {
						pi.SecureURL = val
					} else {
						if strings.HasPrefix(val, "http") {
							pi.URL = val
						} else {
							url, _ := url.Parse(val)
							base, _ := url.Parse(pageURL)
							newURL := (base.ResolveReference(url)).String()
							pi.URL = newURL
						}
					}
				}
				ps.Icon = pi
			}
		}

		if "meta" == token.Data {
			// handles every meta link
			if attrMap["property"] == "og:site_name" {
				ps.SiteName = attrMap["content"]
			}
			if attrMap["property"] == "og:type" {
				ps.Type = attrMap["content"]
			}
			if attrMap["property"] == "og:url" {
				ps.URL = attrMap["content"]
			}
			if attrMap["property"] == "og:title" {
				ps.Title = attrMap["content"]
			}
			if attrMap["name"] == "description" {
				if len(ps.Description) == 0 {
					ps.Description = attrMap["content"]
				}
			}
			if attrMap["property"] == "og:description" {
				ps.Description = attrMap["content"]
			}
			if attrMap["name"] == "author" {
				ps.Author = attrMap["content"]
			}
			if attrMap["name"] == "keywords" {
				s := strings.TrimSpace(attrMap["content"])
				stringArr := strings.Split(s, ",")
				for i := 0; i < len(stringArr); i++ {
					stringArr[i] = strings.TrimSpace(stringArr[i])
				}
				ps.Keywords = stringArr
			}
			// handles Images struct
			imagesArr := []*PreviewImage{}
			if attrMap["property"] == "og:image" {
				pi := &PreviewImage{}
				url, _ := url.Parse(attrMap["content"])
				base, _ := url.Parse(pageURL)
				newURL := (base.ResolveReference(url)).String()
				pi.URL = newURL
				imagesArr = append(imagesArr, pi)
			}
			if attrMap["property"] == "og:image:height" {
				if ps.Images != nil {
					image := ps.Images[len(ps.Images)-1]
					image.Height, _ = strconv.Atoi(attrMap["content"])
				} else {
					image := &PreviewImage{}
					image.Height, _ = strconv.Atoi(attrMap["content"])
					imagesArr = append(imagesArr, image)
				}
			}
			if attrMap["property"] == "og:image:secure_url" {
				if ps.Images != nil {
					image := ps.Images[len(ps.Images)-1]
					image.SecureURL = attrMap["content"]
				} else {
					image := &PreviewImage{}
					image.SecureURL = attrMap["content"]
					imagesArr = append(imagesArr, image)
				}
			}
			if attrMap["property"] == "og:image:width" {
				if len(ps.Images) != 0 {
					image := ps.Images[len(ps.Images)-1]
					width, err := strconv.Atoi(attrMap["content"])
					if err != nil {
						fmt.Errorf("parsing error: %v", err)
					}
					image.Width = width
				} else {
					pi := &PreviewImage{}
					width, err := strconv.Atoi(attrMap["content"])
					if err != nil {
						fmt.Errorf("parsing error: %v", err)
					}
					pi.Width = width
					imagesArr = append(imagesArr, pi)
				}
			}
			if attrMap["property"] == "og:image:alt" {
				if ps.Images != nil {
					image := ps.Images[len(ps.Images)-1]
					image.Alt = attrMap["content"]
				} else {
					image := &PreviewImage{}
					image.Alt = attrMap["content"]
					imagesArr = append(imagesArr, image)
				}
			}
			if attrMap["property"] == "og:image:type" {
				if ps.Images != nil {
					image := ps.Images[len(ps.Images)-1]
					image.Type = attrMap["content"]
				} else {
					image := &PreviewImage{}
					image.Type = attrMap["content"]
					imagesArr = append(imagesArr, image)
				}
			}
			if (len(imagesArr)) != 0 {
				for i := 0; i < len(imagesArr); i++ {
					ps.Images = append(ps.Images, imagesArr[i])
				}
			}
		}
	}
	return ps, nil
}

func trimStringFromSlash(s string) string {
	if idx := strings.LastIndex(s, "/"); idx != -1 {
		return s[:idx]
	}
	return s
}
