package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	heads := headers.NewHeaders()
	heads.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	heads.Set("Connection", "close")
	heads.Set("Content-Type", "text/plain")
	return heads
}
