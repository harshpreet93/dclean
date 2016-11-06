package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"time"
	"regexp"
	"strings"
	"strconv"
	"github.com/pkg/errors"
)

func getAllImgs() []docker.APIImages {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}
	return imgs
}

func printImg(img docker.APIImages) {
	fmt.Println("ID: ", img.ID)
	fmt.Println("RepoTags: ", img.RepoTags)
	fmt.Println("Created: ", time.Unix(img.Created, 0))
	fmt.Println("Size: ", img.Size)
	fmt.Println("VirtualSize: ", img.VirtualSize)
	fmt.Println("ParentId: ", img.ParentID)
}

func printImgs(imgs []docker.APIImages)  {
	for _, img := range imgs {
		printImg(img)
	}
}

func imgsWithTag(imgs []docker.APIImages, tagRegex string) []docker.APIImages {
	var imgsWithTag []docker.APIImages
	for _, img := range imgs {
		matched, _ := regexp.MatchString(tagRegex, img.RepoTags[0])
		if matched {
			imgsWithTag = append(imgsWithTag, img)
		}
	}
	return imgsWithTag
}

func convertRegexMatchToMap(toMatch string, regex *regexp.Regexp) (map[string] string, error ){
	match := regex.FindStringSubmatch(toMatch)
	result := make(map[string]string)
	if len(match) == 0 {
		return nil, errors.New("did not find match for regex "+regex.String())
	}
	for i, name := range regex.SubexpNames() {
		// result[name] = match[i]
		if i != 0 {
			fmt.Println("name "+name)
			result[name] = match[i]
			fmt.Println(result[name])
		}
	}
	return result, nil
}

func convertToNumSecs(timeScalar int, timeUnit string) int {
	//var result int64
	if timeUnit == "m" {
		return timeScalar * 60
	}
	if timeUnit == "h" {
		return timeScalar * 60 * 60
	}
	if timeUnit == "d" {
		return timeScalar * 24 * 60 * 60
	}
	if timeUnit == "w" {
		return timeScalar * 7 * 24 * 60 * 60
	}
	return 0
}

func isOlderThan(img docker.APIImages, filterExp string) bool {
	//createdAt := time.Unix(img.Created, 0)
	trimmed := strings.TrimSpace(filterExp)
	timeFilterExp := regexp.MustCompile(`(?P<num>\d+)(?P<unit>m|h|d|w|y)`)
	timeAndUnit, err := convertRegexMatchToMap(trimmed, timeFilterExp)
	if err != nil {
		return false
	}

	timeScalar, err := strconv.Atoi(timeAndUnit["num"])

	if err != nil {
		return false
	}

	timeUnit := timeAndUnit["unit"]

	querySecs := convertToNumSecs(timeScalar, timeUnit)

	return img.Created < (time.Now().Unix() - int64(querySecs))
	//return true
}


func main() {
	//printImgs( imgsWithTag( getAllImgs(), "^<no.*" ) )
	for _, img := range getAllImgs() {
		fmt.Println("############ img ############")
		printImg(img)
		fmt.Println(isOlderThan(img, "60d"))
	}
	fmt.Println( isOlderThan(getAllImgs()[0], "60d") )

}
