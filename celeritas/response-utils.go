package celeritas

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
)

func (c *Celeritas) WriteJSON(rw http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			rw.Header()[key] = value
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	_, err = rw.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) WriteXML(rw http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			rw.Header()[key] = value
		}
	}

	rw.Header().Set("Content-Type", "application/xml")
	rw.WriteHeader(status)
	_, err = rw.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) DownloadFile(rw http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
	fp := path.Join(pathToFile, fileName)
	fileToServe := filepath.Clean(fp)
	rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; file=\"%s\"", fileName))
	http.ServeFile(rw, r, fileToServe)

	return nil
}

func (c *Celeritas) Error401(rw http.ResponseWriter) {
	c.ErrorStatus(rw, http.StatusUnauthorized)
}

func (c *Celeritas) Error403(rw http.ResponseWriter) {
	c.ErrorStatus(rw, http.StatusForbidden)
}

func (c *Celeritas) Error404(rw http.ResponseWriter) {
	c.ErrorStatus(rw, http.StatusNotFound)
}

func (c *Celeritas) Error500(rw http.ResponseWriter) {
	c.ErrorStatus(rw, http.StatusInternalServerError)
}

func (c *Celeritas) ErrorStatus(rw http.ResponseWriter, status int) {
	http.Error(rw, http.StatusText(status), status)
}
