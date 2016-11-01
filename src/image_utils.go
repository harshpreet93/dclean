package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"time"
)


func getAllImgs() {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}
	for _, img := range imgs {
		//i, _ := strconv.ParseInt(img.Created, 10, 64)

		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", time.Unix(img.Created, 0))
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}
}

func main() {
	fmt.Println("hello, world")
	getAllImgs()
}
