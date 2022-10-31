// Package swagger GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package swagger

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Api service support",
            "email": "support@africarealty.io"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/activation": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "activates a user by token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user id",
                        "name": "userId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "activation token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.User"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "logins user by email/password",
                "parameters": [
                    {
                        "description": "auth request",
                        "name": "loginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "logouts user",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        },
        "/auth/password": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "sets a new password for the user",
                "parameters": [
                    {
                        "description": "set password request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.SetPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        },
        "/auth/registration": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "registers a new client",
                "parameters": [
                    {
                        "description": "registration request",
                        "name": "regRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.RegistrationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.User"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        },
        "/auth/token/refresh": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "refreshes auth token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.SessionToken"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        },
        "/ready": {
            "get": {
                "tags": [
                    "system"
                ],
                "summary": "check system is ready",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.AgentProfile": {
            "type": "object",
            "properties": {
                "avatar": {
                    "description": "Avatar avatar",
                    "type": "string"
                }
            }
        },
        "auth.LoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "Email - login",
                    "type": "string"
                },
                "password": {
                    "description": "Password - password",
                    "type": "string"
                }
            }
        },
        "auth.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "description": "Token - auth token must be passed as  \"Authorization Bearer\" header for all the requests (except ones which don't require authorization)",
                    "$ref": "#/definitions/auth.SessionToken"
                },
                "userId": {
                    "description": "UserId - ID of account",
                    "type": "string"
                }
            }
        },
        "auth.OwnerProfile": {
            "type": "object",
            "properties": {
                "avatar": {
                    "description": "Avatar avatar",
                    "type": "string"
                }
            }
        },
        "auth.RegistrationRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "Email - user email",
                    "type": "string"
                },
                "firstName": {
                    "description": "FirstName - user first name",
                    "type": "string"
                },
                "lastName": {
                    "description": "LastName - user last name",
                    "type": "string"
                },
                "password": {
                    "description": "Password - password",
                    "type": "string"
                },
                "userType": {
                    "description": "UserType - user type",
                    "type": "string"
                }
            }
        },
        "auth.SessionToken": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "description": "AccessToken",
                    "type": "string"
                },
                "accessTokenExpiresAt": {
                    "description": "AccessTokenExpiresAt - when access token expires",
                    "type": "string"
                },
                "refreshToken": {
                    "description": "RefreshToken",
                    "type": "string"
                },
                "refreshTokenExpiresAt": {
                    "description": "RefreshToken - when refresh token expires",
                    "type": "string"
                },
                "sessionId": {
                    "description": "SessionId - session ID",
                    "type": "string"
                }
            }
        },
        "auth.SetPasswordRequest": {
            "type": "object",
            "properties": {
                "newPassword": {
                    "description": "NewPassword - new password",
                    "type": "string"
                },
                "prevPassword": {
                    "description": "PrevPassword - current password",
                    "type": "string"
                }
            }
        },
        "auth.User": {
            "type": "object",
            "properties": {
                "agent": {
                    "description": "Agent - agent profile",
                    "$ref": "#/definitions/auth.AgentProfile"
                },
                "email": {
                    "description": "Email - email",
                    "type": "string"
                },
                "firstName": {
                    "description": "FirstName - user's first name",
                    "type": "string"
                },
                "id": {
                    "description": "Id - user ID",
                    "type": "string"
                },
                "lastName": {
                    "description": "LastName - user's last name",
                    "type": "string"
                },
                "owner": {
                    "description": "Owner - owner profile",
                    "$ref": "#/definitions/auth.OwnerProfile"
                }
            }
        },
        "http.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Code is error code provided by error producer",
                    "type": "string"
                },
                "details": {
                    "description": "Details is additional info provided by error producer",
                    "type": "object",
                    "additionalProperties": true
                },
                "message": {
                    "description": "Message is error description",
                    "type": "string"
                },
                "translationKey": {
                    "description": "TranslationKey is error code translation key",
                    "type": "string"
                },
                "type": {
                    "description": "Type is error type (panic, system, business)",
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "AfricaRealty API",
	Description: "AfricaRealty is an advanced realty service",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
