package standard

import (
	"context"
	"strings"
	"testing"

	"github.com/projectdiscovery/katana/pkg/engine/common"
	"github.com/projectdiscovery/katana/pkg/navigation"
	"github.com/projectdiscovery/katana/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakeRequest_Error_NewRequestWithContext_003 checks that an error from
// http.NewRequestWithContext (e.g., due to an invalid method) is correctly
// propagated.
func TestMakeRequest_Error_NewRequestWithContext_003(t *testing.T) {
	crawler := &Crawler{
		Shared: &common.Shared{
			Options: &types.CrawlerOptions{
				Options: &types.Options{},
			},
		},
	}
	session := &common.CrawlSession{Ctx: context.Background(), Hostname: "example.com"}
	request := &navigation.Request{
		Method: "INVALID METHOD", // This will cause NewRequestWithContext to fail
		URL:    "http://example.com",
	}

	// Act
	response, err := crawler.makeRequest(session, request)

	// Assert
	require.Error(t, err)
	assert.NotNil(t, response) // Response struct is initialized even on error
	assert.Equal(t, request.Depth+1, response.Depth)
	assert.True(t, strings.Contains(err.Error(), "invalid method"))
}

