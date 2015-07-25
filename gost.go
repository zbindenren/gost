package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/zbindenren/gost/configuration"
	"github.com/zbindenren/gost/gist"
)

func main() {
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
	app.Name = "gost"
	app.Usage = "utility to interact with your gists"
	// app.Flags = []cli.Flag{
	// cli.BoolFlag{
	// Name:  "test",
	// Usage: "test",
	// },
	// }
	app.Commands = []cli.Command{
		{
			Name:  "create",
			Usage: "create new gist",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "description",
					Usage: "gist description",
					Value: "",
				},
				cli.StringFlag{
					Name:  "directory",
					Usage: "use files from directory",
					Value: "",
				},
			},
			Action: func(c *cli.Context) {
				if c.IsSet("directory") {
					log.Fatal("directory option not implemented yet")
				} else {
					if len(c.Args()) == 0 {
						log.Fatal("no files given")
					}
					err := client.Post(c.String("description"), c.Args())
					if err != nil {
						log.Fatal(err)
					}
				}
				return
			},
		},
		{
			Name:  "list",
			Usage: "list your gists or files",
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
						fmt.Printf("%s %20s - %s\n", g.ID, strings.Join(files, ", "), g.Description)
					}
				} else {
					gist, err := client.Get(c.Args().First())
					if err != nil {
						log.Fatal(err)
					}
					for _, f := range gist.Files {
						fmt.Printf("%20s %10d %s\n", f.FileName, f.Size, f.RawURL)
					}
				}
			},
		},
		{
			Name:  "delete",
			Usage: "delete gist or file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Usage: "deletes file",
					Value: "",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					log.Fatal("please specify a gist id")
				}
				if c.IsSet("file") {
					log.Fatal("file delete not implemented yet")
				} else {
					err := client.Delete(c.Args().First())
					if err != nil {
						log.Fatal(err)
					}
				}
				return
			},
		},
		{
			Name:  "view",
			Usage: "view gists",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Usage: "specify file",
					Value: "",
				},
				cli.BoolFlag{
					Name:  "browser",
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
			Name: "save",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "destdir",
					Usage: "destination directory",
					Value: "",
				},
				cli.StringFlag{
					Name:  "file",
					Usage: "specify file",
					Value: "",
				},
			},
			Usage: "save gist or file",
			Action: func(c *cli.Context) {
				if c.IsSet("file") {
					log.Fatal("file delete not implemented yet")
				}
				if c.IsSet("destdir") {
					log.Fatal("destdir not implemented yet")
				}
				if len(c.Args()) == 0 {
					log.Fatal("no id given")
				}
				err := client.Download(c.Args().First())
				if err != nil {
					log.Fatal(err)
				}
				return
			},
		},
		{
			Name:  "update",
			Usage: "updates gists",
			Action: func(c *cli.Context) {
				log.Fatal("not impolemented yet")
			},
		},
	}
	app.Run(os.Args)

	// 	if *update {
	// 		if len(flag.Args()) < 2 {
	// 			log.Fatal("gist id and file minimum required")
	// 		}
	// 		err := client.Update(flag.Args()[0], flag.Args()[1:])
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		return
	// 	}

}
