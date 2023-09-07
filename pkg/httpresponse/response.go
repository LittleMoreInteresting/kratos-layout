package httpresponse

import (
	"fmt"
	stdhttp "net/http"
	"time"

	codeproto "github.com/go-kratos/kratos/v2/encoding/proto"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	res "kratos-layout/api/gen/response"
)

// HTTPError is an HTTP error.
type HTTPError struct {
	res.HttpResponse
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError code: %d message: %s", e.Code, e.Message)
}

// FromError try to convert an error to *HTTPError.
func FromError(err error) *HTTPError {

	if err == nil {
		return nil
	}
	e := errors.FromError(err)
	switch {
	case errors.IsServiceUnavailable(e):
		e = errors.New(int(e.Code), "服务不可用", "服务不可用")
	case errors.IsGatewayTimeout(e):
		e = errors.New(int(e.Code), "服务超时", "服务超时")
	case errors.IsBadRequest(e):
		e = errors.New(int(e.Code), "参数错误", "参数错误"+e.GetMessage())
	}
	return &HTTPError{
		res.HttpResponse{
			Code:      int64(e.Code),
			Message:   e.GetMessage(),
			Timestamp: time.Now().UnixMilli(),
		},
	}
}

func ErrorEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
	se := FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	//w.WriteHeader(se.Status)
	_, _ = w.Write(body)
}

type HTTPSuccess struct {
	res.HttpResponse
}

func SuccessResponse(v interface{}) *HTTPSuccess {
	var anyMsg *anypb.Any
	if message, ok := v.(proto.Message); ok {
		anyMsg, _ = anypb.New(message)
	}
	return &HTTPSuccess{
		res.HttpResponse{
			Code:      0,
			Message:   "success",
			Timestamp: time.Now().UnixMilli(),
			Data:      anyMsg,
		},
	}
}

type HTTPDefaultSuccess struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func SuccessDefaultResponse(v interface{}) *HTTPDefaultSuccess {
	return &HTTPDefaultSuccess{
		Code:      0,
		Message:   "success",
		Timestamp: time.Now().UnixMilli(),
		Data:      v,
	}
}

// SuccessResponseEncoder encodes the object to the HTTP response.
func SuccessResponseEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(http.Redirector); ok {
		url, code := rd.Redirect()
		stdhttp.Redirect(w, r, url, code)
		return nil
	}
	codec, _ := http.CodecForRequest(r, "Accept")
	var data []byte
	var err error
	if codec.Name() == codeproto.Name {
		sr := SuccessResponse(v)
		data, err = codec.Marshal(sr)
	} else {
		sr := SuccessDefaultResponse(v)
		data, err = codec.Marshal(sr)
	}
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}
