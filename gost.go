package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

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

	// flag parsing
	add := flag.Bool("a", false, "add gist")
	browser := flag.Bool("b", false, "view gist browser")
	description := flag.String("d", "", "description")
	list := flag.Bool("l", false, "list gists")
	remove := flag.Bool("rm", false, "remove gist")
	save := flag.Bool("s", false, "save gist")
	update := flag.Bool("u", false, "update gist")
	view := flag.Bool("v", false, "view gist")
	flag.Parse()

	if *list {
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
		return
	}

	if *add {
		if len(flag.Args()) == 0 {
			log.Fatal("no files given")
		}
		err := client.Post(*description, flag.Args())
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *update {
		if len(flag.Args()) < 2 {
			log.Fatal("gist id and file minimum required")
		}
		err := client.Update(flag.Args()[0], flag.Args()[1:])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *remove {
		if len(flag.Args()) == 0 {
			log.Fatal("no id given")
		}
		err := client.Delete(flag.Args()[0])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *save {
		if len(flag.Args()) == 0 {
			log.Fatal("no id given")
		}
		err := client.Download(flag.Args()[0])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *view {
		name := ""
		if len(flag.Args()) == 0 {
			log.Fatal("no id given")
		}
		if len(flag.Args()) > 1 {
			name = flag.Args()[1]
		}
		err := client.View(flag.Args()[0], name)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *browser {
		if len(flag.Args()) == 0 {
			log.Fatal("no id given")
		}
		err := client.ViewBrowser(flag.Args()[0])
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}
