package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"reactor-crw"
	"reactor-crw/handler/fs"
	"reactor-crw/parser"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"

	"github.com/spf13/cobra"
)

var (
	search     string
	path       string
	savePath   string
	cookie     string
	maxWorkers int
	allPages   bool

	crawlerCmd = &cobra.Command{
		Use:   "reactor-crw",
		Short: "Grab your favorite content from joyreactor.cc",
		Long:  "Allows to quickly download all content by its direct url or entire" +
				" tag or fandom from joyreactor.cc.\nExample: reactor-crw -d \".\" " +
				"-p \"http://joyreactor.cc/tag/someTag/all\" -w 2 -c \"cookie-string\"",
		Run:   run,
	}
)

func init() {
	hd, _ := os.UserHomeDir()

	crawlerCmd.Flags().StringVarP(&search, "search", "s", "image,gif", "A comma separated list of content types that should be downloaded. Possible values: image,gif,webm,mp4. Example: -s \"image,webm\"")
	crawlerCmd.Flags().StringVarP(&path, "path", "p", "", "Provide a full page URL")
	crawlerCmd.Flags().StringVarP(&savePath, "destination", "d", hd, "Save path for content. User's home folder will be used by default.")
	crawlerCmd.Flags().StringVarP(&cookie, "cookie", "c", "", "User's cookie. Some content may be unavailable without it.")
	crawlerCmd.Flags().IntVarP(&maxWorkers, "workers", "w", 1, "Amount of workers")
	crawlerCmd.Flags().BoolVarP(&allPages, "all", "a", true, "Crawl all related pages using pagination.")

	_ = crawlerCmd.MarkFlagRequired("path")
}

func run(_ *cobra.Command, _ []string) {
	start := time.Now()
	t := reactor_crw.NewHttpTransport(&http.Client{}, reactor_crw.Headers{"Cookie": cookie})

	pathUrl, err := url.Parse(path)
	if err != nil {
		log.Fatalf("invalid path provided: %s", path)
	}

	pr, err := fs.NewPathResolver(savePath)
	if err != nil {
		log.Fatalf("cannot process provided destination: %s", savePath)
	}

	ch, _ := fs.NewFileSaver(pr, t, strings.Replace(pathUrl.Path, "/", "_", -1))
	c := reactor_crw.NewClient(
		&reactor_crw.HtmlCrawler{
			Transport: t,
			Parser:    &parser.Html{},
			MultiPage: allPages,
		},
		maxWorkers,
		ch,
	)

	go func() {
		err = c.Run(path, search)
		if err != nil {
			log.Fatal(err)
		}
	}()

	progress(c.TotalSources, c.Progress, c.Errors)

	fmt.Printf("\n>>> Done in %s\n", time.Since(start).String())
}

func progress(total, task <-chan int, err <-chan error) {
	fmt.Print("\n>>> Trying to count the amount of links. Please wait...\n\n")

	t := <-total
	fmt.Printf(">>> %d links were found. Start downloading\n\n", t)

	prg := mpb.New(mpb.WithWidth(64))

	name := ">>> Progress:"
	bar := prg.AddBar(int64(t),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.CountersNoUnit("%d/%d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(decor.Percentage(decor.WC{W: 5})),
	)

	errBuf := strings.Builder{}

	go func() {
		for {
			select {
			case <-task:
				bar.Increment()
			case e := <-err:
				if e != nil {
					errBuf.WriteString(e.Error() + "\n")
				}
			}
		}
	}()

	prg.Wait()

	if errBuf.Len() != 0 {
		fmt.Printf(">>> Errors during crawler work:\n%s", errBuf.String())
	}
}

func main() {
	err := crawlerCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
