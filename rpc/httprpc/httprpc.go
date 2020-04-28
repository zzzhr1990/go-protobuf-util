package httprpc

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var (
	jsm = &jsonpb.Marshaler{
		EmitDefaults: true,
	}
	jsu = &jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
)

// TryAndUnMarshalStandard post json
func TryAndUnMarshalStandard(method string, url string, postData proto.Message, recvData proto.Message, tryTime int) error {
	var err error = nil
	str := ""
	if postData != nil {
		str, err = jsm.MarshalToString(postData)
		if err != nil {
			return err
		}
	}

	//dataForPost := jsp.Marshal()
	for tryTime > 0 {
		bf := bytes.NewBufferString(str)
		req, err := http.NewRequest(method, url, bf)
		if err != nil {
			tryTime--
			if tryTime > 0 {
				time.Sleep(time.Second * 1)
			}
			err = errors.New("create request for " + url + " failed, error: " + err.Error())
		} else {
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				err = errors.New("request for " + url + " failed, error: " + err.Error())
				tryTime--
				if tryTime > 0 {
					time.Sleep(time.Second * 1)
				}
			} else {
				if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
					// We got normal entity...
					if recvData != nil {
						err = jsu.Unmarshal(resp.Body, recvData)
						if err != nil {
							err = errors.New("unmarshal request for:" + url + " failed, error: " + err.Error())
							tryTime--
							if tryTime > 0 {
								time.Sleep(time.Second * 1)
							}
						} else {
							// Do it ok...
							return nil
						}
					} else {
						// Process OK
						return nil
					}
				} else {
					// trying..
					errorResponse := &ErrorResponse{}
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						err = errors.New("unmarshal request for:" + url + " failed, read body failed error: " + err.Error())
					} else {
						err = json.Unmarshal(body, errorResponse)
						if err != nil {
							err = errors.New("unmarshal request for:" + url + " failed, cannot decode error message, error: " + err.Error())
						} else {
							err = errors.New("request for:" + url + " failed, error: [" + errorResponse.Reference + "] " + errorResponse.Message)
						}
					}
					tryTime--
					if tryTime > 0 {
						time.Sleep(time.Second * 1)
					}
				}
			}
		}
	}
	return err
}
