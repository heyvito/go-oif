package main

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/apex/log"
	cliHandler "github.com/apex/log/handlers/cli"
	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"

	"github.com/heyvito/go-oif/formatter"
)

func findProjectName() (*string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Find our go.mod
	for {
		log.Debugf("Looking for go.mod in %s", cwd)
		fn := filepath.Join(cwd, "go.mod")
		stat, err := os.Stat(fn)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if os.IsNotExist(err) || stat.IsDir() {
			nextCwd := filepath.Dir(cwd)
			if nextCwd == cwd {
				// Root path. Just bail.
				return nil, nil
			}
			cwd = nextCwd
			continue
		}

		// Found it.
		data, err := os.ReadFile(fn)
		if err != nil {
			return nil, err
		} else {
			file, err := modfile.Parse("go.mod", data, nil)
			if err != nil {
				return nil, err
			} else {
				return &file.Module.Mod.Path, nil
			}
		}
	}
}

func main() {
	log.SetHandler(cliHandler.Default)
	projNamePtr, projNameErr := findProjectName()
	projName := ""
	if projNamePtr != nil {
		projName = *projNamePtr
	}

	app := cli.NewApp()
	app.Name = "Opinionated Imports Formatter"
	app.Usage = "Formats imports grouping by source"
	app.Authors = []*cli.Author{
		{Name: "Victor \"Vito\" Gama", Email: "hey@vito.io"},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "project-name",
			Aliases:     []string{"n"},
			Usage:       "Indicates the base project name to sort imports. By default, oif detects the project name using data from your go.mod",
			Value:       projName,
			DefaultText: projName,
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Prints extra information during execution",
		},
	}

	app.Action = func(context *cli.Context) error {
		log.SetLevel(log.ErrorLevel)
		if context.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
			log.Debug("Running in verbose mode")
		}
		switch {
		case projName == "" && projNameErr != nil:
			log.Errorf("Could not read project name from go.mod: %s", projNameErr)
			log.Error("Either use --project-name, or fix your go.mod")
			os.Exit(1)
		case projName == "":
			log.Error("Could not determine project name. Please set it using --project-name, or use go mod")
			os.Exit(1)
		}

		args := context.Args()
		if !args.Present() {
			log.Error("Either provide a list of files to be processed or ./...")
			os.Exit(1)
		}

		var files []string
		if args.Get(0) == "./..." {
			// Recur everything from cwd onwards
			cwd, err := os.Getwd()
			if err != nil {
				log.Errorf("Could not detect cwd: %s", err)
				os.Exit(1)
			}
			err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
				if err == nil && filepath.Ext(info.Name()) == ".go" && !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				log.Errorf("Error listing files: %s", err)
				os.Exit(1)
			}
		} else {
			rawFiles := args.Slice()
			files = make([]string, len(rawFiles))
			var err error
			for i, p := range rawFiles {
				if files[i], err = filepath.Abs(p); err != nil {
					log.Errorf("Error processing file %s: %s", p, err)
					os.Exit(1)
				}
			}
		}

		log.Debugf("Processing %d files", len(files))

		processFiles(files, projName)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("Found errors during execution")
	}
}

func processFiles(paths []string, projName string) {
	maxCPUs := runtime.NumCPU()
	input := make(chan string, 1000)
	toProcess := sync.WaitGroup{}
	toProcess.Add(len(paths))

	log.Debugf("Running with %d goroutines", maxCPUs)

	t := time.Now()
	for i := 0; i < maxCPUs; i++ {
		go process(input, &toProcess, projName)
	}
	for _, p := range paths {
		input <- p
	}

	toProcess.Wait()
	close(input)

	end := time.Now().Sub(t)

	log.WithField("Duration", end.String()).Debug("Done processing")
}

var autoGeneratedPattern = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

func process(f <-chan string, wd *sync.WaitGroup, projName string) {
consumer:
	for path := range f {
		stat, err := os.Stat(path)
		if err != nil {
			panic(err)
		}
		file, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		ll := log.WithField("file", path)

		for i, l := range strings.Split(string(file), "\n") {
			if autoGeneratedPattern.MatchString(l) {
				// File is auto generated. Do not edit.
				ll.Debugf("Detected as autogenerated due to line %d", i+1)
				wd.Done()
				continue consumer
			}
		}

		file = []byte(formatter.FormatImports(projName, string(file)))
		err = os.WriteFile(path, file, stat.Mode())
		if err != nil {
			panic(err)
		}
		ll.Debugf("Done")
		wd.Done()
	}
}
