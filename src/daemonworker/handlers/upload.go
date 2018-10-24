package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/inspursoft/wand/src/daemonworker/utils"
)

func (c *Handler) UploadResource(resp http.ResponseWriter, req *http.Request) {
	fullName := req.FormValue("full_name")
	buildNumber := req.FormValue("build_number")
	if strings.TrimSpace(fullName) == "" || strings.TrimSpace(buildNumber) == "" {
		rendStatus(resp, http.StatusBadRequest, fmt.Sprintln("No repo full name or build number provided."))
		return
	}
	f, fh, err := req.FormFile("upload")
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to resolve uploaded file: %+v\n", err))
		return
	}
	uploadTargetPath := filepath.Join(uploadResourcePath, fullName, buildNumber)
	err = utils.CheckFilePath(uploadTargetPath)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to mkdir: %s, error: %+v\n", uploadTargetPath, err))
		return
	}
	if ext := filepath.Ext(fh.Filename); ext == ".tar" {
		err = utils.Untar(f, uploadTargetPath)
		if err != nil {
			rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to untar file: %s, error: %+v\n", fh.Filename, err))
			return
		}
		return
	}
	err = utils.CopyFile(f, filepath.Join(uploadTargetPath, fh.Filename), 0755)
	if err != nil {
		rendStatus(resp, http.StatusInternalServerError, fmt.Sprintf("Failed to write source to target: %+v\n", err))
	}
}
