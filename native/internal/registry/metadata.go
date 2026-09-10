package registry

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	metadataID          = "microservice-kit-id"
	metadataRoutePrefix = "microservice-kit-routes-"
	metadataChunkSize   = 480
)

func encodeMetadata(instance ServiceInstance) (map[string]string, error) {
	metadata := make(map[string]string, len(instance.Endpoints)+3)
	metadata[metadataID] = instance.ID
	for name, endpoint := range instance.Endpoints {
		metadata[name] = endpoint
	}
	if len(instance.Routes) == 0 {
		return metadata, nil
	}
	data, err := json.Marshal(instance.Routes)
	if err != nil {
		return nil, fmt.Errorf("encode HTTP routes: %w", err)
	}
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := writer.Write(data); err != nil {
		return nil, fmt.Errorf("compress HTTP routes: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("compress HTTP routes: %w", err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(compressed.Bytes())
	for index := 0; len(encoded) > 0; index++ {
		size := min(metadataChunkSize, len(encoded))
		metadata[fmt.Sprintf("%s%02d", metadataRoutePrefix, index)] = encoded[:size]
		encoded = encoded[size:]
	}
	return metadata, nil
}

func decodeMetadata(metadata map[string]string) (string, map[string]string, []HTTPRoute, error) {
	endpoints := make(map[string]string, 2)
	for _, name := range []string{EndpointHTTP, EndpointGRPC} {
		if endpoint := metadata[name]; endpoint != "" {
			endpoints[name] = endpoint
		}
	}
	var keys []string
	for key := range metadata {
		if strings.HasPrefix(key, metadataRoutePrefix) {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return metadata[metadataID], endpoints, nil, nil
	}
	sort.Strings(keys)
	var encoded strings.Builder
	for _, key := range keys {
		encoded.WriteString(metadata[key])
	}
	compressed, err := base64.RawURLEncoding.DecodeString(encoded.String())
	if err != nil {
		return "", nil, nil, fmt.Errorf("decode HTTP routes: %w", err)
	}
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return "", nil, nil, fmt.Errorf("decompress HTTP routes: %w", err)
	}
	data, err := io.ReadAll(reader)
	if closeErr := reader.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return "", nil, nil, fmt.Errorf("decompress HTTP routes: %w", err)
	}
	var routes []HTTPRoute
	if err := json.Unmarshal(data, &routes); err != nil {
		return "", nil, nil, fmt.Errorf("decode HTTP routes: %w", err)
	}
	return metadata[metadataID], endpoints, routes, nil
}
