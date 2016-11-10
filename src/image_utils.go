package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"time"
	"regexp"
	"strings"
	"strconv"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"os"
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

func deleteImg(img docker.APIImages)  {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)

	if err != nil {
		panic(err)
	}

	removalOptions := docker.RemoveImageOptions{
		Force: true,
		NoPrune: true,

	}
	client.RemoveImageExtended(img.ID, removalOptions)
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
		if i != 0 {
			fmt.Println("name "+name)
			result[name] = match[i]
			fmt.Println(result[name])
		}
	}
	return result, nil
}

func convertToNumSecs(timeScalar int, timeUnit string) int {
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
}

func imgsOlderThan(imgs []docker.APIImages, filterExp string) []docker.APIImages {
	var result []docker.APIImages
	for _, img := range imgs {
		if isOlderThan(img, filterExp) {
			result = append(result, img)
		}
	}
	return result
}

func deleteAndSummarize(plan bool, imgs []docker.APIImages) {
	heading := "Deleting images"
	if plan {
		heading = "Running in plan mode so no images will be deleted"
		//fmt.Println("printing imgs")
		printImgs(imgs)
		return;
	}

	fmt.Println(heading)
	for _, img := range imgs {
		fmt.Println("Deleting image with id: "+img.ID+" created at: ",time.Unix(img.Created, 0))
		deleteImg(img)
	}
}

func main() {

	var ageQuery string
	var nameQuery string
	var plan bool
	app := cli.NewApp()
	app.Name = "boom"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) error {
		//fmt.Println(c.Args().First())
		imgs := getAllImgs()
		if nameQuery != "" {
			imgs = imgsWithTag(imgs, nameQuery)
		}
		if ageQuery != "" {
			imgs = imgsOlderThan(imgs, ageQuery)
		}
		deleteAndSummarize(plan, imgs)
		return nil
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "olderthan, o",
			Usage: "delete images based on age",
			Destination: &ageQuery,
		},
		cli.StringFlag{
			Name: "name, n",
			Usage: "delete images based on name",
			Destination: &nameQuery,
		},
		cli.BoolFlag{
			Name: "plan, p",
			Usage: "run in plan mode to see what images would be deleted",
			Destination: &plan,
		},
	}

	app.Run(os.Args)

}
