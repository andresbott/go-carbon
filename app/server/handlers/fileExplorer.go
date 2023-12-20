package handlers

import (
	"encoding/json"
	"fmt"
	file "git.andresbott.com/Golang/carbon/libs/files/filefs"
	logger "git.andresbott.com/Golang/carbon/libs/log"
	"net/http"
)

type FeHandler struct {
	Logger logger.LeveledStructuredLogger
	FS     file.FS
}

type FileEntry struct {
	Name  string
	IsDir bool
	// Size?
}

type FSResponse struct {
	Count int
	Items []FileEntry
}

func (h *FeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// return JSON payload
	h.Logger.Info("req")

	// Parse query parameters
	queryParams := r.URL.Query()

	// Extract the value of the "name" query parameter
	path := queryParams.Get("path")
	// todo handle url encoding

	items, err := h.FS.List(path)

	if err != nil {
		h.Logger.Warn(fmt.Sprintf("failed to list items: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fsItems := fsEntry2FileEntry(items)

	// Create the response payload
	response := FSResponse{
		Count: len(fsItems),
		Items: fsItems,
	}

	// Set the content type header to JSON
	w.Header().Set("Content-Type", "application/json")
	//TODO: this is only for development!!
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//time.Sleep(1 * time.Second)

	// Convert the response payload to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		h.Logger.Warn(fmt.Sprintf("failed to marshal JSON response: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func fsEntry2FileEntry(in []file.FSEntry) []FileEntry {

	ret := make([]FileEntry, len(in))
	for i := 0; i < len(in); i++ {
		item := FileEntry{
			Name:  in[i].Name(),
			IsDir: in[i].IsDir(),
		}
		ret[i] = item
	}
	return ret
}
