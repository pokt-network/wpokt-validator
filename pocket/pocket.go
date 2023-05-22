package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

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

const RemoteCLIURL string = "https://node2.testnet.pokt.network"

// const RemoteCLIURL string = "https://mainnet.gateway.pokt.network/v1/lb/e81ff8d231bc754d1e7e5cd8"

func QueryRPC(path string, jsonArgs []byte) (string, error) {
	cliURL := RemoteCLIURL + path
	// fmt.Println(cliURL)

	req, err := http.NewRequest("POST", cliURL, bytes.NewBuffer(jsonArgs))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 30 * time.Second,
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
	return "", fmt.Errorf("the http status code was not okay: %d, and the status was: %s, with a response of %v", resp.StatusCode, resp)
}

func Height() (*HeightResponse, error) {
	res, err := QueryRPC(GetHeightPath, []byte{})
	if err != nil {
		return nil, err
	}
	var obj HeightResponse
	err = json.Unmarshal([]byte(res), &obj)
	return &obj, err
}

const VaultAddress = "0bee0822d5252eaebf7ae37cf9a6e197202230e5"

func AccountTxs() (*AccountTxsResponse, error) {
	params := rpc.PaginateAddrParams{
		Address:  VaultAddress,
		Page:     1,
		PerPage:  1000,
		Received: true,
		Prove:    true,
		Sort:     "desc",
		Height:   0,
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
