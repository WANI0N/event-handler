// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Marek Beck",
            "email": "marek.beck2@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/event": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Event"
                ],
                "summary": "Creates event to database",
                "parameters": [
                    {
                        "description": "Event Data",
                        "name": "event",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.EventData"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.EventResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/weberrors.AppError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/weberrors.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/weberrors.AppError"
                        }
                    }
                }
            }
        },
        "/event/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Event"
                ],
                "summary": "Retrieves event from database",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Event ID (uuid)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.EventResponseData"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/weberrors.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/weberrors.AppError"
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "Event"
                ],
                "summary": "Delete event from database",
                "parameters": [
                    {
                        "type": "string",
                        "description": "\u003ctoken string goes here\u003e",
                        "name": "API-AUTHENTICATION",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Event ID (uuid)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/weberrors.AppError"
                        }
                    }
                }
            }
        },
        "/healthcheck": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health check"
                ],
                "summary": "Checks health of this service",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.JsonHealthCheckStatus"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.EventData": {
            "description": "If not provided, ` + "`" + `videoQuality` + "`" + ` \u0026 ` + "`" + `audioQuality` + "`" + ` default to ` + "`" + `[\"720p\"]` + "`" + ` \u0026 ` + "`" + `[\"Low\"]` + "`" + `, respectively. If provided, first item in the list is event default.",
            "type": "object",
            "required": [
                "date",
                "invitees",
                "languages",
                "name"
            ],
            "properties": {
                "audioQuality": {
                    "type": "array",
                    "uniqueItems": true,
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Low",
                        "Mid",
                        "High"
                    ]
                },
                "date": {
                    "description": "YYYY-MM-DDTHH:MM:SSZ",
                    "type": "string",
                    "example": "2006-01-02T15:04:05Z"
                },
                "description": {
                    "type": "string",
                    "maxLength": 512
                },
                "invitees": {
                    "type": "array",
                    "maxItems": 100,
                    "minItems": 1,
                    "uniqueItems": true,
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "example@mail.com"
                    ]
                },
                "languages": {
                    "type": "array",
                    "minItems": 1,
                    "uniqueItems": true,
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "English",
                        "French"
                    ]
                },
                "name": {
                    "description": "allowed chars: A-Za-z0-9 _-",
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 1,
                    "example": "A event-Name3_x"
                },
                "videoQuality": {
                    "type": "array",
                    "uniqueItems": true,
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "720p",
                        "1080p",
                        "1440p",
                        "2160p"
                    ]
                }
            }
        },
        "models.EventResponseData": {
            "type": "object",
            "required": [
                "date",
                "invitees",
                "languages",
                "name"
            ],
            "properties": {
                "audioQuality": {
                    "type": "array",
                    "uniqueItems": true,
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Low",
                        "Mid",
                        "High"
                    ]
                },
                "date": {
                    "description": "YYYY-MM-DDTHH:MM:SSZ",
                    "type": "string",
                    "example": "2006-01-02T15:04:05Z"
                },
                "description": {
                    "type": "string",
                    "maxLength": 512
                },
                "id": {
                    "type": "string",
                    "example": "db6bed50-7172-4051-86ab-d1e90705c692"
                },
                "invitees": {
                    "type": "array",
                    "maxItems": 100,
                    "minItems": 1,
                    "uniqueItems": true,
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "example@mail.com"
                    ]
                },
                "languages": {
                    "type": "array",
                    "minItems": 1,
                    "uniqueItems": true,
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "English",
                        "French"
                    ]
                },
                "name": {
                    "description": "allowed chars: A-Za-z0-9 _-",
                    "type": "string",
                    "maxLength": 255,
                    "minLength": 1,
                    "example": "A event-Name3_x"
                },
                "videoQuality": {
                    "type": "array",
                    "uniqueItems": true,
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "720p",
                        "1080p",
                        "1440p",
                        "2160p"
                    ]
                }
            }
        },
        "models.JsonHealthCheckStatus": {
            "type": "object",
            "properties": {
                "deployDate": {
                    "type": "string"
                },
                "result": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "weberrors.AppError": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:3000",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "EventHandler API",
	Description:      "An event management service API in Go using Gin framework.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
