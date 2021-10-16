Reactor Crawler
==========

[![Go Report Card](https://goreportcard.com/badge/github.com/reactor-joy/reactor-crw)](https://goreportcard.com/report/github.com/reactor-joy/reactor-crw)

Simple CLI content crawler for [Joyreactor](http://joyreactor.cc). He'll find all media content on the 
page you've provided and save it. If there will be any kind of pagination... he'll go through all pages
as well unless you'll tell him to not.

<p>
    <img src="https://i.imgur.com/DjD6cW0.gif" width="800" alt="reactor_crawler_example">
</p>

## Quick start

Here's the quickest way to download something and test the crawler:
* Download a build according to your OS.
* Pick some URL from [Joyreactor](http://joyreactor.cc). For example `http://joyreactor.cc/tag/digital+art`.
* Run the crawler `reactor-crw -p "http://joyreactor.cc/tag/digital+art"`

## What else

There's a list of optional flags that adds a little more control over the crawler.
Here is a full output from `$ reactor-crw --help`: 

```
$ reactor-crw --help

Allows to quickly download all content by its direct url or entire tag or fandom from joyreactor.cc.
Example: reactor-crw -d "." -p "http://joyreactor.cc/tag/someTag/all" -w 2 -c "cookie-string"

Usage:
  reactor-crw [flags]

Flags:
  -c, --cookie string        User's cookie. Some content may be unavailable without it
  -d, --destination string   Save path for content. Default value is a user's home folder
                             (example C:\Users\username for Windows) (default "/home/avpretty")
  -h, --help                 help for reactor-crw
  -p, --path string          Provide a full page URL
  -s, --search string        A comma separated list of content types that should be downloaded.
                             Possible values: image,gif,webm,mp4. Example: -s "image,webm" (default "image,gif")
  -o, --single-page          Crawl only one page
  -w, --workers int          Amount of workers (default 1)
```

From all flags only `-p --path` is required. All other flags can be omitted and default values will be used.

Here's another example that will download only `mp4` content from the post and will save it to the current directory:

```sh
$ reactor-crw -p "http://joyreactor.cc/post/000000" -d "." -s "mp4" -o
```