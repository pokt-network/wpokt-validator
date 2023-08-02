package app

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	log "github.com/sirupsen/logrus"
)

func accessSecretVersion(client *secretmanager.Client, name string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return "", err
	}

	return string(result.Payload.Data), nil
}

func readKeysFromGSM() {
	if Config.GoogleSecretManager.Enabled == true {
		log.Debug("[GSM] Reading keys from Google Secret Manager")
	} else {
		log.Debug("[GSM] Google Secret Manager is disabled")
		return
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("[GSM] Failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	if Config.MongoDB.URI == "" && Config.GoogleSecretManager.MongoSecretName == "" {
		log.Fatalf("[GSM] Mongo secret name is empty")
	}

	if Config.GoogleSecretManager.MongoSecretName != "" {
		log.Debug("[GSM] Reading mongo uri")
		Config.MongoDB.URI, err = accessSecretVersion(client, Config.GoogleSecretManager.MongoSecretName)
		if err != nil {
			log.Fatalf("[GSM] Failed to access mongo uri: %v", err)
		}
		log.Info("[GSM] Successfully read mongo uri")
	}

	if Config.Ethereum.PrivateKey == "" && Config.GoogleSecretManager.EthSecretName == "" {
		log.Fatalf("[GSM] Ethereum secret name is empty")
	}

	if Config.GoogleSecretManager.EthSecretName != "" {
		log.Debug("[GSM] Reading ethereum private key")
		Config.Ethereum.PrivateKey, err = accessSecretVersion(client, Config.GoogleSecretManager.EthSecretName)
		if err != nil {
			log.Fatalf("[GSM] Failed to access ethereum private key: %v", err)
		}
		log.Info("[GSM] Successfully read ethereum private key")

	}

	if Config.Pocket.PrivateKey == "" && Config.GoogleSecretManager.PoktSecretName == "" {
		log.Fatalf("[GSM] Pocket secret name is empty")
	}

	if Config.GoogleSecretManager.PoktSecretName != "" {
		log.Debug("[GSM] Reading pocket private key")
		Config.Pocket.PrivateKey, err = accessSecretVersion(client, Config.GoogleSecretManager.PoktSecretName)
		if err != nil {
			log.Fatalf("[GSM] Failed to access pocket private key: %v", err)
		}
		log.Info("[GSM] Successfully read pocket private key")
	}
}
