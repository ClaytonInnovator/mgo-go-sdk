package client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mangonet-labs/mgo-go-sdk/client/httpconn"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/model/response"

	"github.com/tidwall/gjson"
)

// SignAndExecuteTransactionBlock implements the method `mgo_executeTransactionBlock`, signs and executes a transaction.
// The transaction is signed using the Keypair, and the request is sent to the node for execution.
func (c *Client) SignAndExecuteTransactionBlock(ctx context.Context, req request.SignAndExecuteTransactionBlockRequest) (response.MgoTransactionBlockResponse, error) {
	var rsp response.MgoTransactionBlockResponse

	signedTxn, err := req.Keypair.SignTransactionBlock(&req.TxnMetaData)
	if err != nil {
		return rsp, err
	}
	respBytes, err := c.conn.Request(ctx, httpconn.Operation{
		Method: "mgo_executeTransactionBlock",
		Params: []interface{}{
			signedTxn.TxBytes,
			[]string{signedTxn.Signature},
			req.Options,
			req.RequestType,
		},
	})

	if err != nil {
		return rsp, err
	}

	if gjson.ParseBytes(respBytes).Get("error").Exists() {
		return rsp, errors.New(gjson.ParseBytes(respBytes).Get("error").String())
	}

	err = json.Unmarshal([]byte(gjson.ParseBytes(respBytes).Get("result").String()), &rsp)
	if err != nil {
		return rsp, err
	}

	return rsp, nil
}
