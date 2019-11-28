package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"go-note/storage"
	"io/ioutil"
	"os/exec"

	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type Config struct {
	Firebase storage.Config
}

var config Config

func init() {
	// read config file
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		log.Fatalln(err)
	}
	storage.InitStorage(config.Firebase)
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "new",
				Usage:   "new page",
				Aliases: []string{"n"},
			},
			&cli.StringFlag{
				Name:    "edit",
				Usage:   "edit page",
				Aliases: []string{"e"},
			},
			&cli.BoolFlag{
				Name:    "list",
				Usage:   "list pages",
				Aliases: []string{"l"},
			},
		},
		Action: appRun,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func appRun(c *cli.Context) error {
	// display file list
	listFlg := c.Bool("list")
	if listFlg {
		if err := storage.List(); err != nil {
			return err
		}
		return nil
	}

	// edit
	fileName := c.String("edit")
	data := []byte("")
	if fileName != "" {
		fmt.Println("edit file")
		data = storage.Download(fileName)
	} else {
		fileName = c.Args().Get(0)
		// ファイルが指定されているかチェック
		if fileName == "" {
			return fmt.Errorf("no filename")
		}
	}

	// make tmp file
	fp := "./tmp/" + fileName
	err := ioutil.WriteFile(fp, data, 0664)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	fmt.Printf("open tmp file")
	defer func() {
		err = os.Remove(fp)
		if err != nil {
			fmt.Printf(err.Error())
		}
		fmt.Println("delete tmp file")
	}()

	fmt.Printf("open editer\n")
	cmd := exec.Command("vim", fp)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	fmt.Printf("close editer\n")
	content, _ := os.Open(fp)
	defer func() {
		err = content.Close()
		if err != nil {
			fmt.Printf(err.Error())
		}
	}()
	storage.Upload(content)
	return nil
}
