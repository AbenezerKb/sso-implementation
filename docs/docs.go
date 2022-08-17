// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "2F Capital Support Email",
            "url": "http://www.2fcapital.com",
            "email": "info@1f-capital.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/approval": {
            "get": {
                "description": "is used to approve consent.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth2"
                ],
                "summary": "Approval.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "consentId",
                        "name": "consentId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "access",
                        "name": "access",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "headers": {
                            "Location": {
                                "type": "string",
                                "description": "redirect_uri"
                            }
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        },
                        "headers": {
                            "Location": {
                                "type": "string",
                                "description": "redirect_uri"
                            }
                        }
                    }
                }
            }
        },
        "/authorize": {
            "get": {
                "description": "is used to obtain authorization code.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth2"
                ],
                "summary": "Authorize.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "headers": {
                            "Location": {
                                "type": "string",
                                "description": "redirect_uri"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        },
                        "headers": {
                            "Location": {
                                "type": "string",
                                "description": "redirect_uri"
                            }
                        }
                    }
                }
            }
        },
        "/clients": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create a new client",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "Create a client",
                "parameters": [
                    {
                        "description": "client",
                        "name": "client",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Client"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.Client"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/consent/{id}": {
            "get": {
                "description": "is used to get consent by id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth2"
                ],
                "summary": "GetConsentByID.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ConsentData"
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Login a user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login a user.",
                "parameters": [
                    {
                        "description": "login_credential",
                        "name": "login_credential",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.LoginCredential"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.TokenResponse"
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/otp": {
            "get": {
                "description": "is used to request otp for login and signup",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Request otp.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "phone",
                        "name": "phone",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            "login",
                            "signup"
                        ],
                        "type": "string",
                        "description": "type can be login or signup",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Register a new user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user.",
                "parameters": [
                    {
                        "description": "user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RegisterUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "create a new user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "create a new user.",
                "parameters": [
                    {
                        "description": "user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.Client": {
            "type": "object",
            "properties": {
                "client_type": {
                    "description": "ClientType is the type of the client.\nIt can be either confidential or public.",
                    "type": "string"
                },
                "id": {
                    "description": "ID is the unique identifier for the client.\nIt is automatically generated when the client is registered.",
                    "type": "string"
                },
                "logo_url": {
                    "description": "LogoURL is the URL of the client's logo.\nIt must be a valid URL.",
                    "type": "string"
                },
                "name": {
                    "description": "Name is the name of the client that will be displayed to the user.",
                    "type": "string"
                },
                "redirect_uris": {
                    "description": "RedirectURIs is the list of redirect URIs of the client.\nEach redirect URI must be a valid URL and must use HTTPS.",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "scopes": {
                    "description": "Scopes is the list of default scopes of the client if one is not provided.",
                    "type": "string"
                },
                "secret": {
                    "description": "Secret is the secret the client uses to authenticate itself.\nIt is automatically generated when the client is registered.",
                    "type": "string"
                },
                "status": {
                    "description": "Status is the current status of the client.\nIt is set to active by default.",
                    "type": "string"
                }
            }
        },
        "dto.ConsentData": {
            "type": "object",
            "properties": {
                "approved": {
                    "description": "The consent status.",
                    "type": "boolean"
                },
                "client": {
                    "description": "The client data",
                    "$ref": "#/definitions/dto.Client"
                },
                "client_id": {
                    "description": "The client identifier.",
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "redirect_uri": {
                    "description": "The redirection URI used in the initial authorization request.",
                    "type": "string"
                },
                "response_type": {
                    "description": "The redirection URI used in the initial authorization request.",
                    "type": "string"
                },
                "roles": {
                    "description": "Roles of the user.",
                    "type": "string"
                },
                "scope": {
                    "description": "The scope of the access request expressed as a list of space-delimited,",
                    "type": "string"
                },
                "scopes": {
                    "description": "The scope data",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.Scope"
                    }
                },
                "state": {
                    "description": "the state parameter passed in the initial authorization request.",
                    "type": "string"
                },
                "user": {
                    "description": "The user data",
                    "$ref": "#/definitions/dto.User"
                },
                "userID": {
                    "description": "Users Id",
                    "type": "string"
                }
            }
        },
        "dto.CreateUser": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "CreatedAt is the time when the user is created.\nIt is automatically set when the user is created.",
                    "type": "string"
                },
                "email": {
                    "description": "Email is the email of the user.",
                    "type": "string"
                },
                "first_name": {
                    "description": "FirstName is the first name of the user.",
                    "type": "string"
                },
                "gender": {
                    "description": "Gender is the gender of the user.",
                    "type": "string"
                },
                "id": {
                    "description": "ID is the unique identifier of the user.\nIt is automatically generated when the user is created.",
                    "type": "string"
                },
                "last_name": {
                    "description": "LastName is the last name of the user.",
                    "type": "string"
                },
                "middle_name": {
                    "description": "MiddleName is the middle name of the user.",
                    "type": "string"
                },
                "password": {
                    "description": "Password is the password of the user.\nIt is only used for logging in with email",
                    "type": "string"
                },
                "phone": {
                    "description": "Phone is the phone of the user.",
                    "type": "string"
                },
                "profile_picture": {
                    "description": "ProfilePicture is the profile picture of the user.\nIt is set on a separate setProfilePicture endpoint.",
                    "type": "string"
                },
                "role": {
                    "description": "Role is the role given to the user being created.",
                    "type": "string"
                },
                "user_name": {
                    "description": "UserName is the username of the user.\nIt is currently of no use",
                    "type": "string"
                }
            }
        },
        "dto.LoginCredential": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "Email of the user if for login with password",
                    "type": "string"
                },
                "otp": {
                    "description": "OTP generated from phone number",
                    "type": "string"
                },
                "password": {
                    "description": "Password of the user if for login with password",
                    "type": "string"
                },
                "phone": {
                    "description": "Phone number of the user if for login with otp",
                    "type": "string"
                }
            }
        },
        "dto.RegisterUser": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "CreatedAt is the time when the user is created.\nIt is automatically set when the user is created.",
                    "type": "string"
                },
                "email": {
                    "description": "Email is the email of the user.",
                    "type": "string"
                },
                "first_name": {
                    "description": "FirstName is the first name of the user.",
                    "type": "string"
                },
                "gender": {
                    "description": "Gender is the gender of the user.",
                    "type": "string"
                },
                "id": {
                    "description": "ID is the unique identifier of the user.\nIt is automatically generated when the user is created.",
                    "type": "string"
                },
                "last_name": {
                    "description": "LastName is the last name of the user.",
                    "type": "string"
                },
                "middle_name": {
                    "description": "MiddleName is the middle name of the user.",
                    "type": "string"
                },
                "otp": {
                    "description": "OTP is the one time password of the user.",
                    "type": "string"
                },
                "password": {
                    "description": "Password is the password of the user.\nIt is only used for logging in with email",
                    "type": "string"
                },
                "phone": {
                    "description": "Phone is the phone of the user.",
                    "type": "string"
                },
                "profile_picture": {
                    "description": "ProfilePicture is the profile picture of the user.\nIt is set on a separate setProfilePicture endpoint.",
                    "type": "string"
                },
                "user_name": {
                    "description": "UserName is the username of the user.\nIt is currently of no use",
                    "type": "string"
                }
            }
        },
        "dto.Scope": {
            "type": "object",
            "properties": {
                "description": {
                    "description": "The scope description.",
                    "type": "string"
                },
                "name": {
                    "description": "The scope name.",
                    "type": "string"
                }
            }
        },
        "dto.TokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "description": "AccessToken is the access token for the current login",
                    "type": "string"
                },
                "id_token": {
                    "description": "IDToken is the OpenID specific JWT token",
                    "type": "string"
                },
                "refresh_token": {
                    "description": "RefreshToken is the refresh token for the access token",
                    "type": "string"
                },
                "token_type": {
                    "description": "TokenType is the type of token",
                    "type": "string"
                }
            }
        },
        "dto.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "CreatedAt is the time when the user is created.\nIt is automatically set when the user is created.",
                    "type": "string"
                },
                "email": {
                    "description": "Email is the email of the user.",
                    "type": "string"
                },
                "first_name": {
                    "description": "FirstName is the first name of the user.",
                    "type": "string"
                },
                "gender": {
                    "description": "Gender is the gender of the user.",
                    "type": "string"
                },
                "id": {
                    "description": "ID is the unique identifier of the user.\nIt is automatically generated when the user is created.",
                    "type": "string"
                },
                "last_name": {
                    "description": "LastName is the last name of the user.",
                    "type": "string"
                },
                "middle_name": {
                    "description": "MiddleName is the middle name of the user.",
                    "type": "string"
                },
                "password": {
                    "description": "Password is the password of the user.\nIt is only used for logging in with email",
                    "type": "string"
                },
                "phone": {
                    "description": "Phone is the phone of the user.",
                    "type": "string"
                },
                "profile_picture": {
                    "description": "ProfilePicture is the profile picture of the user.\nIt is set on a separate setProfilePicture endpoint.",
                    "type": "string"
                },
                "user_name": {
                    "description": "UserName is the username of the user.\nIt is currently of no use",
                    "type": "string"
                }
            }
        },
        "model.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Code is the error code. It is not status code",
                    "type": "integer"
                },
                "description": {
                    "description": "Description is the error description.",
                    "type": "string"
                },
                "field_error": {
                    "description": "FieldError is the error detail for each field, if available that is.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.FieldError"
                    }
                },
                "message": {
                    "description": "Message is the error message.",
                    "type": "string"
                },
                "stack_trace": {
                    "description": "StackTrace is the stack trace of the error.\nIt is only returned for debugging",
                    "type": "string"
                }
            }
        },
        "model.FieldError": {
            "type": "object",
            "properties": {
                "description": {
                    "description": "Description is the error description for this field.",
                    "type": "string"
                },
                "name": {
                    "description": "Name is the name of the field that caused the error.",
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "206.189.54.235:8000",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "RidePLUS SSO API",
	Description:      "This is the RidePLUS sso api.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
