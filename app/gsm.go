package app

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	log "github.com/sirupsen/logrus"
)

func accessSecretVersion(client *secretmanager.Client, name string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", Config.GoogleSecretManager.ProjectId, name),
	}

	result, err := client.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return "", err
	}

	return string(result.Payload.Data), nil
}

func readKeysFromGSM() {
	if Config.GoogleSecretManager.Enabled == false {
		log.Debug("[GSM] Google Secret Manager is disabled")
		return
	}

	if Config.GoogleSecretManager.ProjectId == "" {
		log.Fatalf("[GSM] ProjectId is empty")
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("[GSM] Failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	if Config.Ethereum.PrivateKey == "" {
		if Config.GoogleSecretManager.EthSecretName == "" {
			log.Fatalf("[GSM] Ethereum secret name is empty")
		}

		log.Debug("[GSM] Reading ethereum private key")
		Config.Ethereum.PrivateKey, err = accessSecretVersion(client, Config.GoogleSecretManager.EthSecretName)
		if err != nil {
			log.Fatalf("[GSM] Failed to access ethereum private key: %v", err)
		}
		log.Info("[GSM] Successfully read ethereum private key")
	}

	if Config.Pocket.PrivateKey == "" {
		if Config.GoogleSecretManager.PoktSecretName == "" {
			log.Fatalf("[GSM] Pocket secret name is empty")
		}

		log.Debug("[GSM] Reading pocket private key")
		Config.Pocket.PrivateKey, err = accessSecretVersion(client, Config.GoogleSecretManager.PoktSecretName)
		if err != nil {
			log.Fatalf("[GSM] Failed to access pocket private key: %v", err)
		}
		log.Info("[GSM] Successfully read pocket private key")
	}
}
