package celeritas

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
)

func (c *Celeritas) ReadJSON(rw http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1 mb
	r.Body = http.MaxBytesReader(rw, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single json value")
	}

	return nil
}

// WriteJSON writes json from arbitrary data
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
