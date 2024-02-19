package googleDrive

import (
	"context"
	"github.com/g82411/gd/utils/strutil"
	"google.golang.org/api/drive/v3"
)

type File struct {
	ID           string
	Name         string
	MimeType     string
	Link         string
	Size         int64
	IsFolder     bool
	ReadableSize string
}

func serializeFileObj(driveFile *drive.File) *File {
	link := driveFile.WebViewLink
	if link == "" {
		link = driveFile.WebContentLink
	}
	isFolder := driveFile.MimeType == "application/vnd.google-apps.folder"

	return &File{
		ID:           driveFile.Id,
		Name:         driveFile.Name,
		MimeType:     driveFile.MimeType,
		Link:         link,
		Size:         driveFile.Size,
		IsFolder:     isFolder,
		ReadableSize: strutil.HumanReadAbleSize(driveFile.Size),
	}
}

func QueryFiles(ctx context.Context, service *drive.Service, query string, f func(files []*File) error) error {
	call := service.Files.List().Q(query).Fields("nextPageToken, files(id, name, mimeType, webViewLink, webContentLink, size)")
	return call.Pages(ctx, func(page *drive.FileList) error {
		files := make([]*File, len(page.Files))
		for i, file := range page.Files {
			link := file.WebViewLink
			if link == "" {
				link = file.WebContentLink
			}
			files[i] = serializeFileObj(file)
		}
		return f(files)
	})
}
