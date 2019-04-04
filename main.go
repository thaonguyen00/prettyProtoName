
package main

import (
"github.com/pkg/errors"
"io/ioutil"
"os"
"github.com/urfave/cli"
"log"
"path/filepath"
"regexp"
"strings"
"fmt"
)
var flags struct {
	FileIn     string
	FileOut     string
	Type 		string

}

func main() {
	app := cli.NewApp()
	app.Name = "structnames"
	app.Usage = "REST APIs for getting agent details"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "input, i",
			Usage:       "inputfile",
			Destination: &flags.FileIn,
		},
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "output file",
			Destination: &flags.FileOut,
		},

	}
	app.Action = launch
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}

func launch(_ *cli.Context) error{
	file, _ := filepath.Abs(flags.FileIn)
	fmt.Println("worked:", file)

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "Failed to read file.")
	}

	r, err := regexp.Compile("\\s[A-Za-z0-9]+\\_sub[0-9]+\\s\\w+\\s")
	if err != nil {
		return errors.Wrap(err, "Failed to execute regex.")
	}
	s := r.FindAllString(string(fileContent), -1)



	var structNames []string
	newFileContent := string(fileContent)
	for _, item := range s {

		matched := strings.TrimSpace(item)
		spaceIndex := strings.Index(matched, " ")
		prettyName :=strings.Title(matched[spaceIndex+1 :len(matched)])
		badName := matched[:spaceIndex]
		newFileContent = strings.Replace(newFileContent, badName, prettyName, -1)

		structNames = append(structNames, item)

	}
	var fOut *os.File
	if flags.FileOut == "" {
		fOut, _ = os.Create(flags.FileIn)
	}else {
		fOut, _ = os.Create(flags.FileOut)
	}

	defer fOut.Close()
	fOut.WriteString(newFileContent)

	return nil
}

