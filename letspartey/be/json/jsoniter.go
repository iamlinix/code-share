package json

import jsoniter "github.com/json-iterator/go"

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal is exported by gin/json package.
	JsonMarshal = json.Marshal
	// Unmarshal is exported by gin/json package.
	JsonUnmarshal = json.Unmarshal
	// MarshalIndent is exported by gin/json package.
	JsonMarshalIndent = json.MarshalIndent
	// NewDecoder is exported by gin/json package.
	JsonNewDecoder = json.NewDecoder
	// NewEncoder is exported by gin/json package.
	JsonNewEncoder = json.NewEncoder
)
