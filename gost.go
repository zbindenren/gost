package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/zbindenren/gost/configuration"
	"github.com/zbindenren/gost/gist"
)

func main() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	c, err := configuration.LoadConfiguration()
	if err != nil {
		if err == configuration.ErrNoConfigFound {
			c, err = configuration.NewConfiguration()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
	err = c.Save()
	if err != nil {
		log.Fatal(err)
	}

	client := gist.New(c)

	app := cli.NewApp()
	app.Version = "1"
	app.Usage = "utility to interact with your gists"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "description, d",
			Usage: "gist description",
		},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			log.Fatal("please specify files to use")
		}
		err := client.Post(c.GlobalString("description"), c.Args())
		if err != nil {
			log.Fatal(err)
		}
	}
	app.Commands = []cli.Command{
		{
			Name:  "ls",
			Usage: "list your gists or files in a gist",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					gists, err := client.List()
					if err != nil {
						log.Fatal(err)
					}
					for _, g := range gists {
						files := []string{}
						for _, f := range g.Files {
							files = append(files, f.FileName)
						}
						fmt.Fprintf(w, "%s\t%s\t- %s\n", g.ID, strings.Join(files, ", "), g.Description)
					}
					w.Flush()
				} else {
					gist, err := client.Get(c.Args().First())
					if err != nil {
						log.Fatal(err)
					}
					for _, f := range gist.Files {
						fmt.Fprintf(w, "%s\t%d\t%s\n", f.FileName, f.Size, f.RawURL)
					}
					w.Flush()
				}
			},
		},
		{
			Name:  "rm",
			Usage: "delete gist or file in a gist",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "deletes file",
					Value: "",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					log.Fatal("please specify a gist id")
				}
				err := client.Delete(c.Args().First(), c.String("file"))
				if err != nil {
					log.Fatal(err)
				}
				return
			},
		},
		{
			Name:  "cat",
			Usage: "view gists",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "specify file name",
					Value: "",
				},
				cli.BoolFlag{
					Name:  "browser, b",
					Usage: "view in browser",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					log.Fatal("please specify a gist id")
				}
				if c.Bool("browser") {
					err := client.ViewBrowser(c.Args().First())
					if err != nil {
						log.Fatal(err)
					}
					return
				}
				err := client.View(c.Args().First(), c.String("file"))
				if err != nil {
					log.Fatal(err)
				}
				return
			},
		},
		{
			Name: "get",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "specify file name",
					Value: "",
				},
			},
			Usage: "download gist or file",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					log.Fatal("gist id missing")
				}
				err := client.Download(c.Args().First(), c.String("file"))
				if err != nil {
					log.Fatal(err)
				}
				return
			},
		},
		{
			Name: "update",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "file, f",
					Usage: "specify file name(s)",
				},
			},
			Usage: "updates gists",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					log.Fatal("gist id missing")
				}
				if len(c.StringSlice("file")) == 0 && !c.GlobalIsSet("description") {
					log.Fatal("file name missing")
				}
				err := client.Update(c.Args().First(), c.GlobalString("description"), c.StringSlice("file"))
				if err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)
}
