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

var (
	SendRawTxPath,
	GetNodePath,
	GetACLPath,
	GetUpgradePath,
	GetDAOOwnerPath,
	GetHeightPath,
	GetAccountPath,
	GetAppPath,
	GetTxPath,
	GetBlockPath,
	GetSupportedChainsPath,
	GetBalancePath,
	GetAccountTxsPath,
	GetNodeParamsPath,
	GetNodesPath,
	GetSigningInfoPath,
	GetAppsPath,
	GetAppParamsPath,
	GetPocketParamsPath,
	GetNodeClaimsPath,
	GetNodeClaimPath,
	GetBlockTxsPath,
	GetSupplyPath,
	GetAllParamsPath,
	GetParamPath,
	GetStopPath,
	GetQueryChains,
	GetAccountsPath string
)

func init() {
	routes := rpc.GetRoutes()
	for _, route := range routes {
		switch route.Name {
		case "SendRawTx":
			SendRawTxPath = route.Path
		case "QueryNode":
			GetNodePath = route.Path
		case "QueryACL":
			GetACLPath = route.Path
		case "QueryUpgrade":
			GetUpgradePath = route.Path
		case "QueryDAOOwner":
			GetDAOOwnerPath = route.Path
		case "QueryHeight":
			GetHeightPath = route.Path
		case "QueryAccount":
			GetAccountPath = route.Path
		case "QueryAccounts":
			GetAccountsPath = route.Path
		case "QueryApp":
			GetAppPath = route.Path
		case "QueryTX":
			GetTxPath = route.Path
		case "QueryBlock":
			GetBlockPath = route.Path
		case "QuerySupportedChains":
			GetSupportedChainsPath = route.Path
		case "QueryBalance":
			GetBalancePath = route.Path
		case "QueryAccountTxs":
			GetAccountTxsPath = route.Path
		case "QueryNodeParams":
			GetNodeParamsPath = route.Path
		case "QueryNodes":
			GetNodesPath = route.Path
		case "QuerySigningInfo":
			GetSigningInfoPath = route.Path
		case "QueryApps":
			GetAppsPath = route.Path
		case "QueryAppParams":
			GetAppParamsPath = route.Path
		case "QueryPocketParams":
			GetPocketParamsPath = route.Path
		case "QueryBlockTxs":
			GetBlockTxsPath = route.Path
		case "QuerySupply":
			GetSupplyPath = route.Path
		case "QueryNodeClaim":
			GetNodeClaimPath = route.Path
		case "QueryNodeClaims":
			GetNodeClaimsPath = route.Path
		case "QueryAllParams":
			GetAllParamsPath = route.Path
		case "QueryParam":
			GetParamPath = route.Path
		case "Stop":
			GetStopPath = route.Path
		case "QueryChains":
			GetQueryChains = route.Path
		default:
			continue
		}
	}
}

func QueryRPC(path string, jsonArgs []byte) (string, error) {
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

func GetBlock() (*BlockResponse, error) {
	res, err := QueryRPC(GetBlockPath, []byte{})
	if err != nil {
		return nil, err
	}
	var obj BlockResponse
	err = json.Unmarshal([]byte(res), &obj)
	return &obj, err
}

func GetHeight() (*HeightResponse, error) {
	res, err := QueryRPC(GetHeightPath, []byte{})
	if err != nil {
		return nil, err
	}
	var obj HeightResponse
	err = json.Unmarshal([]byte(res), &obj)
	return &obj, err
}

func getAccountTxsPerPage(page uint32) (*AccountTxsResponse, error) {
	params := rpc.PaginateAddrParams{
		Address:  app.Config.Copper.VaultAddress,
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
	res, err := QueryRPC(GetAccountTxsPath, j)
	if err != nil {
		return nil, err
	}
	var obj AccountTxsResponse
	err = json.Unmarshal([]byte(res), &obj)
	return &obj, err
}

func GetAccountTransferTxs(height int64) ([]*ResultTx, error) {
	var txs []*ResultTx
	var page uint32 = 1
	for {
		res, err := getAccountTxsPerPage(page)
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

func ValidateNetwork() {
	log.Debugln("Connecting to pocket network", "url", app.Config.Pocket.RPCURL)
	res, err := GetBlock()
	if err != nil {
		panic(err)
	}
	log.Debugln("Connected to pocket network", "height", res.Block.Header.Height)

	if res.Block.Header.ChainID != app.Config.Pocket.ChainId {
		log.Debugln("pocket chainId mismatch", "expected", app.Config.Pocket.ChainId, "got", res.Block.Header.ChainID)
		panic("pocket chain id mismatch")
	}
	log.Debugln("Connected to pocket network", "chainId", app.Config.Pocket.ChainId)
}
