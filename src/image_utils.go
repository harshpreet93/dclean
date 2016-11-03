package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"time"
	"regexp"
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


func main() {
	printImgs( imgsWithTag( getAllImgs(), "^cent.*" ) )

}
