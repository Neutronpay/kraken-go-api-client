package krakenapi

import (
	"encoding/base64"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var publicAPI = New("", "")

func TestKrakenApi(t *testing.T) {
	var kk interface{} = KrakenApi{
		key:    "key",
		secret: "secret",
	}

	name := reflect.TypeOf(kk).Name()
	if name != "KrakenAPI" {
		t.Errorf("Unexpected struct, got %s want %s", name, "KrakenAPI")
	}
}

func TestCreateSignature(t *testing.T) {
	expectedSig := "Uog0MyIKZmXZ4/VFOh0g1u2U+A0ohuK8oCh0HFUiHLE2Csm23CuPCDaPquh/hpnAg/pSQLeXyBELpJejgOftCQ=="
	urlPath := "/0/private/"
	secret, _ := base64.StdEncoding.DecodeString("SECRET")
	values := url.Values{
		"TestKey": {"TestValue"},
	}

	sig := createSignature(urlPath, values, secret)

	if sig != expectedSig {
		t.Errorf("Expected Signature to be %s, got: %s\n", expectedSig, sig)
	}
}

func TestTime(t *testing.T) {
	resp, err := publicAPI.Time()
	if err != nil {
		t.Errorf("Time() should not return an error, got %s", err)
	}

	if resp.Unixtime <= 0 {
		t.Errorf("Time() should return valid Unixtime, got %d", resp.Unixtime)
	}
}

func TestAssets(t *testing.T) {
	_, err := publicAPI.Assets()
	if err != nil {
		t.Errorf("Assets() should not return an error, got %s", err)
	}
}

func TestAssetPairs(t *testing.T) {
	resp, err := publicAPI.AssetPairs()
	if err != nil {
		t.Errorf("AssetPairs() should not return an error, got %s", err)
	}

	if (*resp)["XXBTZEUR"].Base+(*resp)["XXBTZEUR"].Quote != "XXBTZEUR" {
		t.Errorf("AssetPairs() should return valid response, got %+v", (*resp)["XXBTZEU"])
	}
}

func TestTicker(t *testing.T) {
	resp, err := publicAPI.Ticker("XXBTZEUR", "XXRPZEUR")
	if err != nil {
		t.Errorf("Ticker() should not return an error, got %s", err)
	}

	if (*resp)["XXBTZEUR"].OpeningPrice == 0 {
		t.Errorf("Ticker() should return valid OpeningPrice, got %+v", (*resp)["XXBTZEUR"].OpeningPrice)
	}
}

func TestOHLC(t *testing.T) {
	resp, err := publicAPI.OHLC("XXBTZEUR")
	if err != nil {
		t.Errorf("OHLC() should not return an error, got %s", err)
	}

	if resp.Pair == "" {
		t.Errorf("OHLC() should return valid Pair, got %+v", resp.Pair)
	}
}

func TestQueryTime(t *testing.T) {
	result, err := publicAPI.Query("Time", map[string]string{})
	resultKind := reflect.TypeOf(result).Kind()

	if err != nil {
		t.Errorf("Query should not return an error, got %s", err)
	}
	if resultKind != reflect.Map {
		t.Errorf("Query `Time` should return a Map, got: %s", resultKind)
	}
}

func TestQueryTicker(t *testing.T) {
	result, err := publicAPI.Query("Ticker", map[string]string{
		"pair": "XXBTZEUR",
	})
	resultKind := reflect.TypeOf(result).Kind()

	if err != nil {
		t.Errorf("Query should not return an error, got %s", err)
	}

	if resultKind != reflect.Map {
		t.Errorf("Query `Ticker` should return a Map, got: %s", resultKind)
	}
}

func TestQueryTrades(t *testing.T) {
	result, err := publicAPI.Trades("XXBTZEUR", 1495777604391411290)

	if err != nil {
		t.Errorf("Trades should not return an error, got %s", err)
	}

	if result.Last == 0 {
		t.Errorf("Returned parameter `last` should always have a value...")
	}

	if len(result.Trades) > 0 {
		for _, trade := range result.Trades {
			if trade.Buy == trade.Sell {
				t.Errorf("Trade should be buy or sell")
			}
			if trade.Market == trade.Limit {
				t.Errorf("Trade type should be market or limit")
			}
		}
	}
}

func TestQueryDepth(t *testing.T) {
	pair := "XETHZEUR"
	count := 10
	result, err := publicAPI.Depth(pair, count)
	if err != nil {
		t.Errorf("Depth should not return an error, got %s", err)
	}

	resultType := reflect.TypeOf(result)

	if resultType != reflect.TypeOf(&OrderBook{}) {
		t.Errorf("Depth should return an OrderBook, got %s", resultType)
	}

	if len(result.Asks) > count {
		t.Errorf("Asks length must be less than count , got %d > %d", len(result.Asks), count)
	}

	if len(result.Bids) > count {
		t.Errorf("Bids length must be less than count , got %d > %d", len(result.Bids), count)
	}
}

func TestDepositStatus(t *testing.T) {
	// Create a test API instance with dummy credentials for structure testing
	testAPI := New("test_key", "test_secret")

	// Test with empty struct
	_, err := testAPI.DepositStatus(DepositStatusParams{})
	if err == nil {
		t.Errorf("DepositStatus should return an authentication error with dummy credentials")
	}

	// Verify the error contains authentication-related text
	errorStr := err.Error()
	if !strings.Contains(errorStr, "API") {
		t.Errorf("Expected API-related error, got: %s", errorStr)
	}

	// Test with populated struct using helper methods
	params := DepositStatusParams{
		Asset:  "XXBT",
		Method: "Bitcoin",
	}
	params.SetLimit(5).SetTimeRange(1640995200, 1672531200)

	_, err = testAPI.DepositStatus(params)
	if err == nil {
		t.Errorf("DepositStatus should return an authentication error with dummy credentials")
	}

	// Verify the error contains authentication-related text
	errorStr = err.Error()
	if !strings.Contains(errorStr, "API") {
		t.Errorf("Expected API-related error, got: %s", errorStr)
	}

	// Test with only start time (deposits since timestamp)
	paramsStartOnly := DepositStatusParams{Asset: "XETH"}
	paramsStartOnly.SetStartTime(1640995200).SetLimit(10)

	_, err = testAPI.DepositStatus(paramsStartOnly)
	if err == nil {
		t.Errorf("DepositStatus should return an authentication error with dummy credentials")
	}

	// Test with only end time (deposits before timestamp)
	paramsEndOnly := DepositStatusParams{Asset: "XXBT"}
	paramsEndOnly.SetEndTime(1672531200)

	_, err = testAPI.DepositStatus(paramsEndOnly)
	if err == nil {
		t.Errorf("DepositStatus should return an authentication error with dummy credentials")
	}
}

func TestDepositMethods(t *testing.T) {
	// Create a test API instance with dummy credentials for structure testing
	testAPI := New("test_key", "test_secret")

	// Test with asset parameter
	params := DepositMethodsParams{
		Asset: "XXBT",
	}

	_, err := testAPI.DepositMethods(params)
	// We expect this to fail with authentication error since we're using dummy credentials
	// but it validates that the method structure and parameter handling works
	if err == nil {
		t.Errorf("DepositMethods should return an authentication error with dummy credentials")
	}

	// Verify the error contains authentication-related text
	errorStr := err.Error()
	if !strings.Contains(errorStr, "API") {
		t.Errorf("Expected API-related error, got: %s", errorStr)
	}

	// Test with both asset and aclass parameters
	paramsWithAclass := DepositMethodsParams{
		Asset:  "XXBT",
		Aclass: "currency",
	}

	_, err = testAPI.DepositMethods(paramsWithAclass)
	if err == nil {
		t.Errorf("DepositMethods should return an authentication error with dummy credentials")
	}

	// Verify the error contains authentication-related text
	errorStr = err.Error()
	if !strings.Contains(errorStr, "API") {
		t.Errorf("Expected API-related error, got: %s", errorStr)
	}

	// Test with only aclass parameter
	paramsOnlyAclass := DepositMethodsParams{
		Aclass: "crypto",
	}

	_, err = testAPI.DepositMethods(paramsOnlyAclass)
	if err == nil {
		t.Errorf("DepositMethods should return an authentication error with dummy credentials")
	}

	// Verify the error contains authentication-related text
	errorStr = err.Error()
	if !strings.Contains(errorStr, "API") {
		t.Errorf("Expected API-related error, got: %s", errorStr)
	}

	// Test with empty parameters (should still work as the API will handle validation)
	emptyParams := DepositMethodsParams{}

	_, err = testAPI.DepositMethods(emptyParams)
	if err == nil {
		t.Errorf("DepositMethods should return an authentication error with dummy credentials")
	}

	// Verify the error contains authentication-related text
	errorStr = err.Error()
	if !strings.Contains(errorStr, "API") {
		t.Errorf("Expected API-related error, got: %s", errorStr)
	}
}
