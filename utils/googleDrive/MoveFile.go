package googleDrive

import (
	"context"
	"google.golang.org/api/drive/v3"
)

func MoveFile(ctx context.Context, service *drive.Service, fileId string, folderId string) error {
	_, err := service.Files.Update(fileId, nil).AddParents(folderId).RemoveParents("root").Do()
	return err

}
