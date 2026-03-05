/*
Copyright © 2025 Yi-Chen (Jean) Li jean841115@gmail.com
*/
package main

import (
	"github.com/jeanyichenli/FileUploadSystem/cmd"
	"github.com/jeanyichenli/FileUploadSystem/redis"
)

func main() {
	// init redis
	redis.InitRedisClient()

	// execute cmds
	cmd.Execute()
}
