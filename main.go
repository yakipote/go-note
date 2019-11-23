package main

import (
	"fmt"
	"go-note/storage"
	"io/ioutil"
	"os/exec"

	"log"
	"os"

	"github.com/urfave/cli/v2"
)

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
	//listFlg := c.Bool("list")
	//if listFlg {
	//	o, _ := exec.Command("ls", NoteBook).Output()
	//	os.Stdout.Write(o)
	//	return nil
	//}

	fileName := c.Args().Get(0)
	// ファイルが指定されているかチェック
	if fileName == "" {
		return fmt.Errorf("no filename")
	}

	// make tmp file
	fp := "./tmp/" + fileName
	err := ioutil.WriteFile(fp, []byte(""), 0664)
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
		fmt.Printf("delete tmp file")
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
