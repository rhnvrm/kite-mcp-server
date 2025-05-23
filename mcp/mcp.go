package mcp

import (
	"fmt"
	"log"

	gomcp "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/zerodha/kite-mcp-server/kc"
)

// TODO: add destructive, openworld and readonly hints where applicable.

var (
	ToolList []Tool = []Tool{ // TODO: this does not need to be global, it could be moved inside the function that calls RegisterTools
		// Tools for setting up the client
		&LoginTool{},

		// Tools that get data from Kite Connect
		&ProfileTool{},
		&MarginsTool{},
		&HoldingsTool{},
		&PositionsTool{},
		&TradesTool{},
		&OrdersTool{},
		&GTTOrdersTool{},
		&MFHoldingsTool{},

		// Tools for market data
		&QuotesTool{},
		&InstrumentsSearchTool{},
		&HistoricalDataTool{},

		// Tools that post data to Kite Connect
		&PlaceOrderTool{},
		&ModifyOrderTool{},
		&CancelOrderTool{},
		&PlaceGTTOrderTool{},
		&ModifyGTTOrderTool{},
		&DeleteGTTOrderTool{},
	}
)

type Tool interface {
	Tool() gomcp.Tool
	Handler(*kc.Manager) server.ToolHandlerFunc
}

func RegisterTools(srv *server.MCPServer, manager *kc.Manager) {
	for _, tool := range ToolList {
		srv.AddTool(tool.Tool(), tool.Handler(manager))
	}
}

// Utilities for assertions

func assertString(v any) string {
	if v == nil {
		return ""
	}

	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}

	return s
}

func assertInt(v any) int {
	if v == nil {
		return 0
	}

	i, ok := v.(int)
	if !ok {
		// Try to assert if it is float64, if so, convert it to int
		return int(assertFloat64(v))
	}

	return i
}

func assertFloat64(v any) float64 {
	if v == nil {
		return 0.0
	}

	f, ok := v.(float64)
	if !ok {
		return 0.0
	}

	return f
}

func assertStringArray(v any) []string {
	if v == nil {
		return nil
	}

	arr, ok := v.([]any)
	if !ok {
		log.Printf("debug actual type: %T", v)
		return nil
	}

	out := make([]string, len(arr))
	for i, item := range arr {
		out[i] = assertString(item)
	}

	return out
}

func assertBool(v any) bool {
	if v == nil {
		return false
	}

	b, ok := v.(bool)
	if !ok {
		// Check if it is a string and convert it to bool
		s := assertString(v)
		if s == "true" {
			return true
		} else if s == "false" {
			return false
		}

		return false
	}

	return b
}
