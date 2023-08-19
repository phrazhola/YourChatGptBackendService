package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// return user id. if user doesn't exist in db, return empty string
// email is encoded email string
func GetUser(email string, container *azcosmos.ContainerClient) (string, error) {
	if container == nil {
		ctn, err := ConnectDatabase("", "", "")
		if err != nil {
			return "", err
		}

		container = ctn
	}

	query := fmt.Sprintf("SELECT * FROM c WHERE c.id = '%s'", email)
	queryPage := container.NewQueryItemsPager(query, azcosmos.NewPartitionKeyString(string(email[0])), nil)

	if !queryPage.More() {
		return "", nil
	} else {
		queryResponse, err := queryPage.NextPage(context.Background())
		if err != nil {
			return "", err
		}
	
		if len(queryResponse.Items) == 0 {
			return "", nil
		}
	}

	itemResponse, err := container.ReadItem(context.Background(), azcosmos.NewPartitionKeyString(string(email[0])), email, nil)
	if err != nil {
		return "", err
	}

	var itemResponseBody map[string]interface{}
	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)

	return itemResponseBody["user_id"].(string), nil
}

