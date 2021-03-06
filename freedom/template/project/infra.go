package project

import "fmt"

func init() {
	content["/infra/response.go"] = jsonResponseTemplate()
	content["/infra/request.go"] = jsonRequestTemplate()
}

func jsonRequestTemplate() string {
	return `
	//Package infra generated by 'freedom new-project {{.PackagePath}}'
	package infra

	import (
		"io/ioutil"
	
		"encoding/json"
		"github.com/8treenet/freedom"
		"gopkg.in/go-playground/validator.v9"
	)
	
	var validate *validator.Validate
	func init() {
		validate = validator.New()
		freedom.Prepare(func(initiator freedom.Initiator) {
			initiator.BindInfra(false, func() *Request {
				return &Request{}
			})
			initiator.InjectController(func(ctx freedom.Context) (com *Request) {
				initiator.FetchInfra(ctx, &com)
				return
			})
		})
	}
	
	// Request .
	type Request struct {
		freedom.Infra
	}
	
	// BeginRequest .
	func (req *Request) BeginRequest(worker freedom.Worker) {
		req.Infra.BeginRequest(worker)
	}
	
	// ReadJSON .
	func (req *Request) ReadJSON(obj interface{}) error {
		rawData, err := ioutil.ReadAll(req.Worker().IrisContext().Request().Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(rawData, obj)
		if err != nil {
			return err
		}
	
		return validate.Struct(obj)
	}
	
	// ReadQuery .
	func (req *Request) ReadQuery(obj interface{}) error {
		if err := req.Worker().IrisContext().ReadQuery(obj); err != nil {
			return err
		}
		return validate.Struct(obj)
	}
	
	// ReadForm .
	func (req *Request) ReadForm(obj interface{}) error {
		if err := req.Worker().IrisContext().ReadForm(obj); err != nil {
			return err
		}
		return validate.Struct(obj)
	}	
	`
}
func jsonResponseTemplate() string {
	result := `
	//Package infra generated by 'freedom new-project {{.PackagePath}}'
	package infra
	import (
		"encoding/json"
		"strconv"
	
		"github.com/8treenet/freedom"
		"github.com/kataras/iris/v12/hero"
	)
	
	// JSONResponse .
	type JSONResponse struct {
		Code             int
		Error            error
		Object           interface{}
		DisableLogOutput bool
	}
	
	// Dispatch .
	func (jrep JSONResponse) Dispatch(ctx freedom.Context) {
		contentType := "application/json"
		var content []byte
	
		var body struct {
			Code  int         %sjson:"code"%s
			Error string      %sjson:"error"%s
			Data  interface{} %sjson:"data,omitempty"%s
		}
		body.Data = jrep.Object
		body.Code = jrep.Code
	
		if jrep.Error != nil {
			body.Error = jrep.Error.Error()
		}
		if jrep.Error != nil && body.Code == 0 {
			body.Code = 400
		}
	
		if content, jrep.Error = json.Marshal(body); jrep.Error != nil {
			content = []byte(jrep.Error.Error())
		}
	
		ctx.Values().Set("code", strconv.Itoa(body.Code))
		if !jrep.DisableLogOutput {
			ctx.Values().Set("response", string(content))
		}
	
		hero.DispatchCommon(ctx, 200, contentType, content, nil, nil, true)
	}	

`
	return fmt.Sprintf(result, "`", "`", "`", "`", "`", "`")
}
