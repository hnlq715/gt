package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/samber/lo"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Usage = "Simple git repo clone tool with workspace support"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "workspace",
			Value: filepath.Join(lo.Must(os.UserHomeDir()), "workspace"),
			Usage: "workspace path",
		},
	}
	app.Action = func(c *cli.Context) error {
		gitpath := c.Args().First()
		if gitpath == "" {
			return errors.New("empty git path is invalid")
		}

		parsedPath := lo.Must(url.Parse(gitpath))
		clonePath := filepath.Join(
			lo.Must(filepath.Abs(c.String("workspace"))), parsedPath.Host, parsedPath.Path)

		_, err := os.Stat(clonePath)
		if os.IsNotExist(err) {
			fmt.Println("git", "clone", gitpath, clonePath)
			lo.Must0(exec.Command("git", "clone", gitpath, clonePath).Run())
		} else {
			fmt.Println(clonePath, "already exists")
		}

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
