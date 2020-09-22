// Copyright (c) 2020, Geert JM Vanderkelen

package xgraphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/eventeneer/xkit/xhttp"
)

// Variables defines a key/value mapping which can be used to pass
// variables when executing GraphQL queries.
type Variables map[string]interface{}

type resultData struct {
	Data   json.RawMessage `json:"data"`
	Errors []Error         `json:"errors"`
}

// Client defines a GraphQL client connecting to the endpoint using
// HTTP.
type Client struct {
	*xhttp.Client
}

// NewClient returns a new GraphQL Client connecting to server using
// uri. The options argument can be used to configure the underlying
// HTTP Client.
func NewClient(uri string, options ...xhttp.ClientOption) *Client {
	return &Client{
		Client: xhttp.NewClient(uri, options...),
	}
}

// Payload defines what we send to the GraphQL endpoint.
type Payload struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables,omitempty"`
}

// BlackHoler can be implemented by struct types for which the API result is not
// wanted. This is mostly useful for testing.
type BlackHoler interface {
	DiscardResult()
}

// Execute executes the GraphQL query and stores payload in result using the
// connection information found in c.
// If result implements the xgraphql.BlackHoler interface, the result is
// discarded (not decoded).
func (c Client) Execute(query string, result interface{}) error {
	return c.execute(query, result, nil)
}

// ExecuteWithVars executes the GraphQL query with variables in vars and
// stores payload in result using the connection information found in c.
// If result implements the xgraphql.BlackHoler interface, the result is
// discarded (not decoded).
func (c Client) ExecuteWithVars(query string, result interface{}, vars Variables) error {
	return c.execute(query, result, vars)
}

func (c Client) execute(query string, result interface{}, vars Variables) error {
	payload := &Payload{
		Query:     query,
		Variables: vars,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.Client.Post(bytes.NewReader(body))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed sending request; got %d", resp.StatusCode)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data resultData
	if err := json.Unmarshal(buf, &data); err != nil {
		return err
	}

	if len(data.Errors) > 0 {
		return data.Errors[0]
	}

	if _, isABlackHole := result.(BlackHoler); isABlackHole {
		// apparently, we are not interested in the result
		return nil
	}

	err = json.Unmarshal(data.Data, result)
	return err
}
