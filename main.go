package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
	"github.com/urfave/cli"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			startingPath := c.Args().Get(0)
			metaDataFileContent := createMetaDataStrings(2)
			originalFilePaths := getFilePaths(startingPath)
			metaDataFilePaths := createMetaDataFilePaths(originalFilePaths)
			writeMetaDataFiles(metaDataFilePaths, metaDataFileContent)

			return nil
		},
		Name:    "unity-metadata-gen",
		Author:  "Recluse Games",
		Version: "0.0.1",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func writeMetaDataFiles(metaDataFilePaths []string, metaDataFileContent []string) {
	for _, metaDataFilePath := range metaDataFilePaths {
		writeMetaDataFile(metaDataFilePath, metaDataFileContent)
	}
}

func writeMetaDataFile(path string, metaDataStrings []string) {
	f, err := os.Create(path)
	check(err)

	defer f.Close()

	for _, metaDataString := range metaDataStrings {
		_, err := f.WriteString(metaDataString)
		check(err)
	}

	f.Sync()

	log.Output(1, "Metadata File Succesfully Created At Path: "+path)
}

func createMetaDataFilePaths(paths []string) []string {
	metaDataFilePaths := []string{}

	for _, pathName := range paths {
		metaDataFilePaths = append(metaDataFilePaths, createMetaDataFilePath(pathName))
	}

	return metaDataFilePaths
}

func createMetaDataFilePath(path string) string {
	parentFilePathRegex := *regexp.MustCompile(`(.*/|\\)(?:.*)$`)
	parentDirectoryPath := parentFilePathRegex.FindStringSubmatch(path)
	if len(parentDirectoryPath) != 2 {
		log.Fatal("Error parsing File Name")
	}
	metaDataFileName := createMetaDataFileName(path)
	metaDataPathName := parentDirectoryPath[1] + metaDataFileName

	return metaDataPathName
}

func createMetaDataFileName(path string) string {
	fileNameRegex := *regexp.MustCompile(`(?:.*/|\\)(.*)$`)
	stringMatches := fileNameRegex.FindStringSubmatch(path)
	if len(stringMatches) != 2 {
		log.Fatal("Error parsing File Name")
	}
	metaDataFileName := stringMatches[1] + ".meta"

	return metaDataFileName
}

func createMetaDataStrings(metaDataVersion int) []string {
	metaDataFile := []string{
		fmt.Sprintf("fileFormattedVersion: %d\n", metaDataVersion),
		fmt.Sprintf("guid: %s\n", uuid.New()),
		fmt.Sprintf("MonoImporter:\n"),
		fmt.Sprintf("\texternalObjects: {}\n"),
		fmt.Sprintf("\tserializedVersion: %d\n", metaDataVersion),
		fmt.Sprintf("\tdefaultReferences: []\n"),
		fmt.Sprintf("\texecutionOrder: 0\n"),
		fmt.Sprintf("\ticon: {instanceID: 0}\n"),
		fmt.Sprintf("\tuserData: \n"),
		fmt.Sprintf("\tassetBundleName: \n"),
		fmt.Sprintf("\tassetBundleVariant: \n"),
	}

	return metaDataFile
}

func getFilePaths(startingPath string) []string {
	var files []string

	err := filepath.Walk(startingPath, func(path string, info os.FileInfo, err error) error {
		// Skip directories
		if filepath.Ext(path) == ".cs" {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}
