// Code generated by statik. DO NOT EDIT.

package embed

import (
	"github.com/rakyll/statik/fs"
)


func init() {
	data := "PK\x03\x04\x14\x00\x08\x00\x00\x00Mr)Q\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00client.ts.tmplUT\x05\x00\x013\xe4X_{{define \"client\"}}\n{{- if .Services}}\n//\n// Client\n//\n\n{{- range .Services}}\nexport class {{.Name}} implements {{.Name | serviceInterfaceName}} {\n  private hostname: string\n  private fetch: Fetch\n  private path = '/rpc/{{.Name}}/'\n\n  constructor(hostname: string, fetch: Fetch) {\n    this.hostname = hostname\n    this.fetch = fetch\n  }\n\n  private url(name: string): string {\n    return this.hostname + this.path + name\n  }\n  {{range .Methods}}\n  {{.Name | methodName}} = ({{. | methodInputs}}): {{. | methodOutputs}} => {\n    return this.fetch(\n      this.url('{{.Name}}'),\n      {{- if .Inputs | len}}\n      createHTTPRequest(args, headers, signal),\n      {{- else}}\n      createHTTPRequest({}, headers, signal),\n      {{end}}\n    ).then((res) => {\n      return buildResponse(res).then(_data => {\n        return {\n        {{- $outputsCount := .Outputs|len -}}\n        {{- range $i, $output := .Outputs}}\n          {{$output | newOutputArgResponse}}{{listComma $i $outputsCount}}\n        {{- end}}\n        }\n      })\n    })\n  }\n  {{end}}\n}\n{{end -}}\n{{end -}}\n{{end}}\nPK\x07\x08M\x05wl-\x04\x00\x00-\x04\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xee\xa0MQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x16\x00	\x00client_helpers.ts.tmplUT\x05\x00\x01\x00	\x86_{{define \"client_helpers\"}}\nexport interface WebRPCError extends Error {\n  code: string\n  msg: string\n	status: number\n}\n\nconst createHTTPRequest = (body: object = {}, headers: object = {}, signal: AbortSignal = null): object => {\n  return {\n    method: 'POST',\n    headers: { ...headers, 'Content-Type': 'application/json' },\n    body: JSON.stringify(body || {}),\n    signal\n  }\n}\n\nconst buildResponse = (res: Response): Promise<any> => {\n  const { status } = res\n  return res.text().then(text => {\n    let data\n    try {\n      if (text == \"\") {\n        data = { status }\n      } else {\n        data = { ...JSON.parse(text), status }\n      }\n    } catch(err) {\n      throw { code: 'unknown', msg: `expecting JSON, got: ${text}`, status } as WebRPCError\n    }\n    if (!res.ok) {\n      throw data // webrpc error response\n    }\n    return data\n  })\n}\n\nexport type Fetch = (input: RequestInfo, init?: RequestInit, abortSignal?: AbortSignal) => Promise<Response>\n{{end}}\nPK\x07\x08+\xc5\x01\x01\xc7\x03\x00\x00\xc7\x03\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00;\xa1MQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x11\x00	\x00proto.gen.ts.tmplUT\x05\x00\x01\x93	\x86_{{- define \"proto\" -}}\n/* tslint:disable */\n// {{.Name}} {{.SchemaVersion}}\n// --\n// This file has been generated by https://github.com/webrpc/webrpc using gen/typescript\n// Do not edit by hand. Update your webrpc schema and re-generate.\n\n// WebRPC description and code-gen version\nexport const WebRPCVersion = \"{{.WebRPCVersion}}\"\n\n// Schema version of your RIDL schema\nexport const WebRPCSchemaVersion = \"{{.SchemaVersion}}\"\n\n{{template \"types\" .}}\n\n{{- if .TargetOpts.Client}}\n  {{template \"client\" .}}\n  {{template \"client_helpers\" .}}\n{{- end}}\n\n{{- if .TargetOpts.Server}}\n  {{template \"server\" .}}\n  {{template \"server_helpers\" .}}\n{{- end}}\n\n{{- end}}\nPK\x07\x086\xe0@\x83\x94\x02\x00\x00\x94\x02\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00server.ts.tmplUT\x05\x00\x01\x0e\x85\xc5^{{define \"server\"}}\n\n{{- if .Services}}\n//\n// Server\n//\nexport class WebRPCError extends Error {\n    statusCode?: number\n\n    constructor(msg: string = \"error\", statusCode?: number) {\n        super(\"webrpc error: \" + msg);\n\n        Object.setPrototypeOf(this, WebRPCError.prototype);\n\n        this.statusCode = statusCode;\n    }\n}\n\nimport express from 'express'\n\n    {{- range .Services}}\n        {{$name := .Name}}\n        {{$serviceName := .Name | serviceInterfaceName}}\n\n        export type {{$serviceName}}Service = {\n            {{range .Methods}}\n                {{.Name}}: (args: {{.Name}}Args) => {{.Name}}Return | Promise<{{.Name}}Return>\n            {{end}}\n        }\n\n        export const create{{$serviceName}}App = (serviceImplementation: {{$serviceName}}Service) => {\n            const app = express();\n\n            app.use(express.json())\n\n            app.post('/*', async (req, res) => {\n                const requestPath = req.baseUrl + req.path\n\n                if (!req.body) {\n                    res.status(400).send(\"webrpc error: missing body\");\n\n                    return\n                }\n\n                switch(requestPath) {\n                    {{range .Methods}}\n\n                    case \"/rpc/{{$name}}/{{.Name}}\": {                        \n                        try {\n                            {{ range .Inputs }}\n                                {{- if not .Optional}}\n                                    if (!(\"{{ .Name }}\" in req.body)) {\n                                        throw new WebRPCError(\"Missing Argument `{{ .Name }}`\")\n                                    }\n                                {{end -}}\n\n                                if (\"{{ .Name }}\" in req.body && !validateType(req.body[\"{{ .Name }}\"], \"{{ .Type | jsFieldType }}\")) {\n                                    throw new WebRPCError(\"Invalid Argument: {{ .Name }}\")\n                                }\n                            {{end}}\n\n                            const response = await serviceImplementation[\"{{.Name}}\"](req.body);\n\n                            {{ range .Outputs}}\n                                if (!(\"{{ .Name }}\" in response)) {\n                                    throw new WebRPCError(\"internal\", 500);\n                                }\n                            {{end}}\n\n                            res.status(200).json(response);\n                        } catch (err) {\n                            if (err instanceof WebRPCError) {\n                                const statusCode = err.statusCode || 400\n                                const message = err.message\n\n                                res.status(statusCode).json({\n                                    msg: message,\n                                    status: statusCode,\n                                    code: \"\"\n                                });\n\n                                return\n                            }\n\n                            if (err.message) {\n                                res.status(400).send(err.message);\n\n                                return;\n                            }\n\n                            res.status(400).end();\n                        }\n                    }\n\n                    return;\n                    {{end}}\n\n                    default: {\n                        res.status(404).end()\n                    }\n                }\n            });\n\n            return app;\n        };\n    {{- end}}\n{{end -}}\n{{end}}\nPK\x07\x08XT\xde\xa9\x92\x0d\x00\x00\x92\x0d\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x16\x00	\x00server_helpers.ts.tmplUT\x05\x00\x01\x0e\x85\xc5^{{ define \"server_helpers\" }}\n\nconst JS_TYPES = [\n    \"bigint\",\n    \"boolean\",\n    \"function\",\n    \"number\",\n    \"object\",\n    \"string\",\n    \"symbol\",\n    \"undefined\"\n]\n\n{{ range .Messages }}\n    const validate{{ .Name }} = (value: any) => {\n        {{ range .Fields }}\n            {{ if .Optional }}\n                if (\"{{ . | exportedJSONField }}\" in value && !validateType(value[\"{{ . | exportedJSONField }}\"], \"{{ .Type | jsFieldType }}\")) {\n                    return false\n                }\n            {{ else }}\n                if (!(\"{{ . | exportedJSONField }}\" in value) || !validateType(value[\"{{ . | exportedJSONField }}\"], \"{{ .Type | jsFieldType }}\")) {\n                    return false\n                }\n            {{ end }}\n        {{ end }}\n\n        return true\n    }\n{{ end }}\n\nconst TYPE_VALIDATORS: { [type: string]: (value: any) => boolean } = {\n    {{ range .Messages }}\n        {{ .Name }}: validate{{ .Name }},\n    {{ end }}\n}\n\nconst validateType = (value: any, type: string) => {\n    if (JS_TYPES.indexOf(type) > -1) {\n        return typeof value === type;\n    }\n\n    const validator = TYPE_VALIDATORS[type];\n\n    if (!validator) {\n        return false;\n    }\n\n    return validator(value);\n}\n\n{{ end }}PK\x07\x08\x93\xb2\xd6w\xce\x04\x00\x00\xce\x04\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0d\x00	\x00types.ts.tmplUT\x05\x00\x01\x0e\x85\xc5^{{define \"types\"}}\n//\n// Types\n//\n\n{{- if .Messages -}}\n{{range .Messages -}}\n\n{{if .Type | isEnum -}}\n{{$enumName := .Name}}\nexport enum {{$enumName}} {\n{{- range $i, $field := .Fields}}\n  {{- if $i}},{{end}}\n  {{$field.Name}} = '{{$field.Name}}'\n{{- end}}\n}\n{{end -}}\n\n{{- if .Type | isStruct  }}\nexport interface {{.Name | interfaceName}} {\n  {{- range .Fields}}\n  {{if . | exportableField -}}{{. | exportedJSONField}}{{if .Optional}}?{{end}}: {{.Type | fieldType}}{{- end -}}\n  {{- end}}\n}\n{{end -}}\n{{end -}}\n{{end -}}\n\n{{if .Services}}\n{{- range .Services}}\nexport interface {{.Name | serviceInterfaceName}} {\n{{- range .Methods}}\n  {{.Name | methodName}}({{. | methodInputs}}): {{. | methodOutputs}}\n{{- end}}\n}\n\n{{range .Methods -}}\nexport interface {{. | methodArgumentInputInterfaceName}} {\n{{- range .Inputs}}\n  {{.Name}}{{if .Optional}}?{{end}}: {{.Type | fieldType}}\n{{- end}}\n}\n\nexport interface {{. | methodArgumentOutputInterfaceName}} {\n{{- range .Outputs}}\n  {{.Name}}{{if .Optional}}?{{end}}: {{.Type | fieldType}}\n{{- end}}  \n}\n{{end}}\n\n{{- end}}\n{{end -}}\n{{end}}\nPK\x07\x08=A*\xed=\x04\x00\x00=\x04\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00Mr)QM\x05wl-\x04\x00\x00-\x04\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x00client.ts.tmplUT\x05\x00\x013\xe4X_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xee\xa0MQ+\xc5\x01\x01\xc7\x03\x00\x00\xc7\x03\x00\x00\x16\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81r\x04\x00\x00client_helpers.ts.tmplUT\x05\x00\x01\x00	\x86_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00;\xa1MQ6\xe0@\x83\x94\x02\x00\x00\x94\x02\x00\x00\x11\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x86\x08\x00\x00proto.gen.ts.tmplUT\x05\x00\x01\x93	\x86_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4PXT\xde\xa9\x92\x0d\x00\x00\x92\x0d\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81b\x0b\x00\x00server.ts.tmplUT\x05\x00\x01\x0e\x85\xc5^PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P\x93\xb2\xd6w\xce\x04\x00\x00\xce\x04\x00\x00\x16\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x819\x19\x00\x00server_helpers.ts.tmplUT\x05\x00\x01\x0e\x85\xc5^PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xa9\x9b\xb4P=A*\xed=\x04\x00\x00=\x04\x00\x00\x0d\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81T\x1e\x00\x00types.ts.tmplUT\x05\x00\x01\x0e\x85\xc5^PK\x05\x06\x00\x00\x00\x00\x06\x00\x06\x00\xb0\x01\x00\x00\xd5\"\x00\x00\x00\x00"
		fs.Register(data)
	}
	