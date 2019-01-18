package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

var extensions = []string{
	".png",
	".bmp",
	".jpeg",
	".jpg",
	".tiff",
}

func indexWallpapers(path string) ([]string) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			for _, ext := range extensions {
				r, err := regexp.MatchString(ext, info.Name())
				if err == nil && r {
					files = append(files, path)
				}
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func main() {
	path := flag.String("path", ".", "absolute path to the wallpapers")
	duration := flag.Int("duration", 30, "duration for each wallpaper to be shown in seconds")

	flag.Parse()

	if 0 == len(*path) {
		os.Exit(1)
	}

	file, err := os.Stat(*path)
	if err != nil {
		panic(err)
	}

	if !file.IsDir() {
		log.Fatal("Incorrect path")
	}

	files := indexWallpapers(*path)
	fmt.Println("Found", len(files), "wallpapers")

	if len(files) == 0 {
		os.Exit(1)
	}

	for {
		for _, wp := range files {
			cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+wp)
			_ = cmd.Run()
			time.Sleep(time.Second * time.Duration(*duration))
		}

		// reindex
		files = indexWallpapers(*path)
	}

}
