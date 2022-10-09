package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func EncodedMarshalJSON(data interface{}, queryVal string, diagnostics io.Writer) ([]byte, error) {
	indent, errIndent := strconv.Atoi(queryVal)
	var marshalData []byte
	var marshalErr error
	if errIndent != nil {
		fmt.Printf("Can not read indent %d, default value will be 0", indent)
	}
	if indent > 0 && indent < 15 && errIndent == nil {
		marshalData, marshalErr = json.MarshalIndent(data, "", strings.Repeat(" ", indent))
	} else {
		marshalData, marshalErr = json.Marshal(data)
	}
	if marshalErr != nil {
		fmt.Fprintf(diagnostics, "Couldn't proceed with Marshal: %v\n", marshalErr)
		return nil, marshalErr
	}
	return marshalData, nil
}
