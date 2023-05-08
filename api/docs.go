// Code generated by swaggo/swag. DO NOT EDIT.

package api

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/albums": {
            "get": {
                "description": "get all the albums in the store",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "albums"
                ],
                "summary": "Get all Albums",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Album"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "add a new album to the store",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "albums"
                ],
                "summary": "Create album",
                "parameters": [
                    {
                        "description": "album",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Album"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Album"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ServerError"
                        }
                    }
                }
            }
        },
        "/albums/{id}": {
            "get": {
                "description": "get as single album by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "albums"
                ],
                "summary": "Get Album by id",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "int valid",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Album"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ServerError"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "description": "get Prometheus metrics for the service",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "albums"
                ],
                "summary": "Prometheus metrics",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Album": {
            "type": "object",
            "required": [
                "artist",
                "price",
                "title"
            ],
            "properties": {
                "artist": {
                    "type": "string",
                    "maxLength": 1000,
                    "minLength": 2
                },
                "id": {
                    "type": "integer",
                    "maximum": 10000,
                    "minimum": 1
                },
                "price": {
                    "type": "number",
                    "maximum": 10000,
                    "minimum": 0
                },
                "title": {
                    "type": "string",
                    "maxLength": 1000,
                    "minLength": 2
                }
            }
        },
        "model.BindingErrorMsg": {
            "type": "object",
            "required": [
                "field",
                "message"
            ],
            "properties": {
                "field": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "model.ServerError": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.BindingErrorMsg"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Album Store API",
	Description:      "Simple golang album store CRUD application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
