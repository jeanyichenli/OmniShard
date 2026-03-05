/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jeanyichenli/FileUploadSystem/api/handlers"
	"github.com/spf13/cobra"
)

const (
	filePathOptionName       = "file-path"
	filePathShortOptionName  = "f"
	chunkSizeOptionName      = "chunk-size"
	chunkSizeShortOptionName = "c"
)

// Transform unit string to unit bytes
var sizeMap = map[string]float64{
	"k": 1024,
	"K": 1024,
	"m": 1024 * 1024,
	"M": 1024 * 1024,
	"g": 1024 * 1024 * 1024,
	"G": 1024 * 1024 * 1024,
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file and split the file into chunks",
	Long: `Upload command is going to get the uploaded file
	and split the file into chunks. In recent, the chunked data
	can only save into disk files.
	You can use options like
	--file-path to define specific file path.
	Or 
	--chunk-size to decide the size to split the file. ex. 1M, 1k, 256k,...`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upload called")

		// Get flags
		filepath, chunkString, err := getUploadOptions(cmd)
		if err != nil {
			fmt.Printf("Error parsing options %s, err: %s\n", filepath, err.Error())
			return
		}

		// Calculate chunk size
		chunkerSize, err := transfromChunkString(chunkString)
		if err != nil {
			fmt.Printf("Error transforming chunk size, err: %s\n", err.Error())
			return
		}

		// Open file
		// TODO: Find other method to handle large file sizes like 100G.
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Error opening file %s, err: %s\n", filepath, err.Error())
			return
		}
		defer file.Close()

		// Get file size
		fi, err := file.Stat()
		if err != nil {
			fmt.Printf("Error getting status of file %s, err: %s\n", filepath, err.Error())
			return
		}

		fileSize := fi.Size()                                     // Total file size
		totalChunks := math.Ceil(float64(fileSize) / chunkerSize) // Chunk number (upper round)

		// POST the file to upload API
		// TODO: program local http module to separate http function calls and main flow

		// Prepare body
		uploadContent := handlers.UploadSession{
			Filename:    filepath,
			TotalSize:   fileSize,
			ChunkSize:   int64(chunkerSize),
			TotalChunks: int64(totalChunks),
		}

		uploadBody, err := json.Marshal(uploadContent)
		if err != nil {
			fmt.Printf("Error json marshaling when uploading, err: %s\n", err.Error())
		}

		// Send request
		url := "http://localhost:8080/upload"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(uploadBody))
		if err != nil {
			fmt.Printf("Error creating new request to %s, err: %s\n", url, err.Error())
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending http request to %s, err: %s\n", url, err.Error())
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body, err: %s\n", err.Error())
			return
		}

		fmt.Println("Test result:", string(body))
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	uploadCmd.Flags().StringP(filePathOptionName, filePathShortOptionName, "", "File path to upload.")
	uploadCmd.MarkFlagRequired(filePathOptionName)
	uploadCmd.Flags().StringP(chunkSizeOptionName, chunkSizeShortOptionName, "256k", "Chunk size to split file. Default is 256kB.")
}

func getUploadOptions(cmd *cobra.Command) (filepath, chunkSize string, err error) {
	filepath, err = cmd.Flags().GetString(filePathOptionName)
	if err != nil {
		fmt.Printf("Error getting %s option: %s", filePathOptionName, err.Error())
		return filepath, chunkSize, err
	}

	chunkSize, err = cmd.Flags().GetString(chunkSizeOptionName)
	if err != nil {
		fmt.Printf("Error getting %s option: %s\n", chunkSizeOptionName, err.Error())
		return filepath, chunkSize, err
	}

	return filepath, chunkSize, err
}

func transfromChunkString(chunkString string) (chunkerSize float64, err error) {
	// Check string format
	if !checkChunkStringFormat(chunkString) {
		return 0, fmt.Errorf("invalid chunker size format")
	}

	unitStr := chunkString[len(chunkString)-1:]
	numberStr := chunkString[:len(chunkString)-1]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		fmt.Printf("Error converting string to int, err: %s\n", err.Error())
		return 0, err
	}

	chunkerSize = float64(number) * sizeMap[unitStr]

	return chunkerSize, nil
}

func checkChunkStringFormat(chunkString string) bool {
	re := regexp.MustCompile(`^\d+[KMGkmg]$`) // Number + Unit
	return re.MatchString(chunkString)
}
