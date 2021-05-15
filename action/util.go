package action

import (
	"encoding/base64"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	DisableColors bool
)

func PrintFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(buffer))
}

func base64Decode(str string) []byte {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println("Cannot decode certificate: ", err)
	}
	return data
}

func createFileIfNotExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("Could not find file %s. Will create one", path)
		parent := filepath.Dir(path)
		err = os.MkdirAll(parent, 0700)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, nil, 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func printColorized(message string, colorAttr color.Attribute, args ...interface{}) {
	color.NoColor = DisableColors
	color.Set(colorAttr)
	defer color.Unset()
	fmt.Printf(message, args...)
}

func PrintWarn(message string, args ...interface{}) {
	printColorized(message, color.FgYellow, args...)
}

func PrintError(message string, args ...interface{}) {
	printColorized(message, color.FgRed, args...)
}
func PrintSuccess(message string, args ...interface{}) {
	printColorized(message, color.FgGreen, args...)
}
