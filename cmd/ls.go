package cmd

import (
	"fmt"
	"github.com/g82411/gd/utils/googleDrive"
	"github.com/urfave/cli/v2"
	"log"
	"sync"
	"time"
)

type ListTask struct {
	FolderID string
	Depth    int
	Prefix   string
	Retry    int
	Query    string
}

const (
	poolSize               = 15
	chanSize               = int(1e5)
	maxRetry               = 3
	maxCheck               = 2
	checkFrequencyInSecond = 3
)

func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:  "ls",
		Usage: "List files and folder in from some folder in google drive",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "depth",
				Aliases:     []string{"d"},
				Usage:       "Depth of recursive",
				DefaultText: "15",
				Value:       15,
			},
		},
		Action: func(c *cli.Context) error {
			if !c.Args().Present() {
				log.Fatalf("Insuffient command")
			}
			folderName := c.Args().First()
			return List{
				FolderName: folderName,
				Recursive:  c.Bool("recursive"),
				Depth:      c.Int("depth"),
			}.Run(c)
		},
	}
}

type List struct {
	Fields     string
	FolderName string
	Recursive  bool
	Depth      int
	Output     string
}

func (l List) Run(c *cli.Context) error {
	tasks := make(chan ListTask, chanSize)
	srv, err := googleDrive.GetService(c.Context)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	checkEmptyTicker := time.NewTicker(checkFrequencyInSecond * time.Second)
	defer checkEmptyTicker.Stop()
	worker := func() {
		defer wg.Done()
		noTaskCounter := 0
		for {
			select {
			case task, ok := <-tasks:
				if !ok {
					return
				}
				query := task.Query
				err := googleDrive.QueryFiles(c.Context, srv, query, func(files []*googleDrive.File) error {
					for _, file := range files {
						if file.IsFolder && task.Depth < l.Depth {
							tasks <- ListTask{
								FolderID: file.ID,
								Depth:    task.Depth + 1,
								Prefix:   task.Prefix + file.Name + "/",
								Query:    fmt.Sprintf("'%s' in parents and trashed = false", file.ID),
							}
						}
						fullPath := task.Prefix + file.Name
						log.Printf("%s, %4s, %s", fullPath, file.ReadableSize, file.Link)
					}
					return nil
				})
				if err != nil {
					log.Printf("Failed to query files: %v", err)
					if task.Retry < maxRetry {
						task.Retry++
						tasks <- task
					}
				}
			case <-checkEmptyTicker.C:
				if len(tasks) == 0 {
					noTaskCounter++
				}
				if noTaskCounter >= maxCheck {
					return
				}
			}
		}
	}
	tasks <- ListTask{
		Query:  fmt.Sprintf("name = '%s'", l.FolderName),
		Depth:  0,
		Prefix: "",
		Retry:  0,
	}
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go worker()
	}
	wg.Wait()
	close(tasks)
	return nil
}
