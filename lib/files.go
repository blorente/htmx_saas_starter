package lib

import (
	"fmt"

	"github.com/pocketbase/pocketbase/models"
)

func GetFileUrl(model *models.Record, field string) (string, error) {
	collection_id := model.TableName()
	record_id := model.Id
	filename := model.GetString(field)
	if filename == "" {
		return "", fmt.Errorf("failed to get data from field")
	}
	rendered := fmt.Sprintf("/api/files/%s/%s/%s", collection_id, record_id, filename)
	return rendered, nil

}
