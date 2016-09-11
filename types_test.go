package osc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestV1Json(t *testing.T) {
	v1 := []byte(`
    {
        "manufacturer": "AAA",
        "model": "BBB",
        "serialNumber": "CCC",
        "firmwareVersion": "DDD",
        "supportUrl": "EEE",
        "endpoints": {
            "httpPort": 80,
            "httpUpdatesPort": 10080,
            "httpsPort": 443,
            "httpsUpdatesPort": 10443            
        },
        "gps": true,
        "gyro": false,
        "uptime": 600,
        "api": [
            "/osc/info",
            "/osc/state",
            "/osc/checkForUpdates",
            "/osc/commands/execute",
            "/osc/commands/status"
        ],
        "_v0": true,
        "_v1": "string"        
    }
    `)

	info, err := parseInfo(bytes.NewReader(v1))

	require.NoError(t, err, "error returned reading V1 JSON", err)
	require.NotNil(t, info)

	require.Equal(t, info.Manufacturer, "AAA")
	require.Equal(t, info.Model, "BBB")
	require.Equal(t, info.Serial, "CCC")
	require.Equal(t, info.Firmware, "DDD")
	require.Equal(t, info.Support, "EEE")
	require.Equal(t, info.Endpoints.HTTP.Server, 80)
	require.Equal(t, info.Endpoints.HTTP.Updates, 10080)
	require.Equal(t, info.Endpoints.HTTPS.Server, 443)
	require.Equal(t, info.Endpoints.HTTPS.Updates, 10443)
	require.True(t, info.GPS)
	require.False(t, info.Gyro)
	require.Equal(t, info.Uptime, 600)
	require.Equal(t, len(info.API), 5)
	require.Contains(t, info.API, "/osc/info")
	require.Contains(t, info.API, "/osc/state")
	require.Contains(t, info.API, "/osc/checkForUpdates")
	require.Contains(t, info.API, "/osc/commands/execute")
	require.Contains(t, info.API, "/osc/commands/status")
	require.Equal(t, info.Vendor["v0"], true)
	require.Equal(t, info.Vendor["v1"], "string")

}
