package pocket

import (
	log "github.com/sirupsen/logrus"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dan13ram/wpokt-backend/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
)

type PocketClient interface {
	GetBlock() (*BlockResponse, error)
	GetHeight() (*HeightResponse, error)
	GetAccountTxsByHeight(height int64) ([]*ResultTx, error)
	ValidateNetwork()
}

type pocketClient struct{}

var (
	Client PocketClient = &pocketClient{}
)

var (
	sendRawTxPath,
	getNodePath,
	getACLPath,
	getUpgradePath,
	getDAOOwnerPath,
	getHeightPath,
	getAccountPath,
	getAppPath,
	getTxPath,
	getBlockPath,
	getSupportedChainsPath,
	getBalancePath,
	getAccountTxsPath,
	getNodeParamsPath,
	getNodesPath,
	getSigningInfoPath,
	getAppsPath,
	getAppParamsPath,
	getPocketParamsPath,
	getNodeClaimsPath,
	getNodeClaimPath,
	getBlockTxsPath,
	getSupplyPath,
	getAllParamsPath,
	getParamPath,
	getStopPath,
	getQueryChains,
	getAccountsPath string
)

func init() {
	routes := rpc.GetRoutes()
	for _, route := range routes {
		switch route.Name {
		case "SendRawTx":
			sendRawTxPath = route.Path
		case "QueryNode":
			getNodePath = route.Path
		case "QueryACL":
			getACLPath = route.Path
		case "QueryUpgrade":
			getUpgradePath = route.Path
		case "QueryDAOOwner":
			getDAOOwnerPath = route.Path
		case "QueryHeight":
			getHeightPath = route.Path
		case "QueryAccount":
			getAccountPath = route.Path
		case "QueryAccounts":
			getAccountsPath = route.Path
		case "QueryApp":
			getAppPath = route.Path
		case "QueryTX":
			getTxPath = route.Path
		case "QueryBlock":
			getBlockPath = route.Path
		case "QuerySupportedChains":
			getSupportedChainsPath = route.Path
		case "QueryBalance":
			getBalancePath = route.Path
		case "QueryAccountTxs":
			getAccountTxsPath = route.Path
		case "QueryNodeParams":
			getNodeParamsPath = route.Path
		case "QueryNodes":
			getNodesPath = route.Path
		case "QuerySigningInfo":
			getSigningInfoPath = route.Path
		case "QueryApps":
			getAppsPath = route.Path
		case "QueryAppParams":
			getAppParamsPath = route.Path
		case "QueryPocketParams":
			getPocketParamsPath = route.Path
		case "QueryBlockTxs":
			getBlockTxsPath = route.Path
		case "QuerySupply":
			getSupplyPath = route.Path
		case "QueryNodeClaim":
			getNodeClaimPath = route.Path
		case "QueryNodeClaims":
			getNodeClaimsPath = route.Path
		case "QueryAllParams":
			getAllParamsPath = route.Path
		case "QueryParam":
			getParamPath = route.Path
		case "Stop":
			getStopPath = route.Path
		case "QueryChains":
			getQueryChains = route.Path
		default:
			continue
		}
	}
}

func queryRPC(path string, jsonArgs []byte) (string, error) {
	cliURL := app.Config.Pocket.RPCURL + path

	req, err := http.NewRequest("POST", cliURL, bytes.NewBuffer(jsonArgs))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * time.Duration(app.Config.Pocket.RPCTimeOutSecs),
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res, err := strconv.Unquote(string(bz))
	if err == nil {
		bz = []byte(res)
	}

	if resp.StatusCode == http.StatusOK {
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, bz, "", "    ")
		if err == nil {
			return prettyJSON.String(), nil
		}
		return string(bz), nil
	}
	return "", fmt.Errorf("the http status code was not okay: %d, with a response of %+v", resp.StatusCode, resp)
}

func (c *pocketClient) GetBlock() (*BlockResponse, error) {
	res, err := queryRPC(getBlockPath, []byte{})
	if err != nil {
		return nil, err
	}
	var obj BlockResponse
	err = json.Unmarshal([]byte(res), &obj)
	return &obj, err
}

func (c *pocketClient) GetHeight() (*HeightResponse, error) {
	res, err := queryRPC(getHeightPath, []byte{})
	if err != nil {
		return nil, err
	}
	var obj HeightResponse
	err = json.Unmarshal([]byte(res), &obj)
	return &obj, err
}

func (c *pocketClient) getAccountTxsPerPage(page uint32) (*AccountTxsResponse, error) {
	params := rpc.PaginateAddrParams{
		Address:  app.Config.Pocket.VaultAddress,
		Page:     int(page),
		PerPage:  1000,
		Received: true,
		Prove:    true,
		Sort:     "asc",
	}
	j, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	res, err := queryRPC(getAccountTxsPath, j)
	if err != nil {
		return nil, err
	}
	var obj AccountTxsResponse
	err = json.Unmarshal([]byte(res), &obj)
	return &obj, err
}

func (c *pocketClient) GetAccountTxsByHeight(height int64) ([]*ResultTx, error) {
	var txs []*ResultTx
	var page uint32 = 1
	for {
		res, err := c.getAccountTxsPerPage(page)
		if err != nil {
			return nil, err
		}
		lastHeight := res.Txs[len(res.Txs)-1].Height
		// filter only type pos/Send
		for _, tx := range res.Txs {
			if tx.StdTx.Msg.Type == "pos/Send" && tx.Height >= height {
				txs = append(txs, tx)
			}
		}
		if len(txs) >= int(res.TotalTxs) || lastHeight < height || len(res.Txs) == 0 {
			break
		}
		page++
	}

	return txs, nil
}

func (c *pocketClient) ValidateNetwork() {
	log.Debugln("[POCKET] Validating network")
	log.Debugln("[POCKET] URL", app.Config.Pocket.RPCURL)
	res, err := c.GetBlock()
	if err != nil {
		log.Errorln("[POCKET] Error getting block", err)
		panic(err)
	}
	log.Debugln("[POCKET] Validating network", "chainId", res.Block.Header.ChainID)
	log.Debugln("[POCKET] Validating network", "height", res.Block.Header.Height)
	if res.Block.Header.ChainID != app.Config.Pocket.ChainId {
		log.Debugln("[POCKET] Chain ID mismatch", "expected", app.Config.Pocket.ChainId, "got", res.Block.Header.ChainID)
		panic("[POCKET] Chain ID mismatch")
	}
	log.Debugln("[POCKET] Validated network")
}
