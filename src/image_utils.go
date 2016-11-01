package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	//"time"
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
	//for _, img := range imgs {
	//	//i, _ := strconv.ParseInt(img.Created, 10, 64)
	//
	//	fmt.Println("ID: ", img.ID)
	//	fmt.Println("RepoTags: ", img.RepoTags)
	//	fmt.Println("Created: ", time.Unix(img.Created, 0))
	//	fmt.Println("Size: ", img.Size)
	//	fmt.Println("VirtualSize: ", img.VirtualSize)
	//	fmt.Println("ParentId: ", img.ParentID)
	//}
}

func imgsWithTag(imgs []docker.APIImages, tagRegex string) {
	//fmt.Println()
	
}


func main() {
	fmt.Println("hello, world")
	getAllImgs()
}
