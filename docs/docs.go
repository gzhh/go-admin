// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/user": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "用户-列表接口",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "页号",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "maximum": 1000,
                        "minimum": 1,
                        "type": "integer",
                        "description": "每页大小",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handlers.resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/admin.UserListResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "返回结构中有token，则说明授权成功；若无token，则说明账号在审核。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "用户-登录",
                "parameters": [
                    {
                        "description": "数据参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/validators.UserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handlers.resp"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/admin.UserResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "admin.UserListResponse": {
            "type": "object",
            "properties": {
                "list": {
                    "description": "数据列表",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/admin.UserResponse"
                    }
                },
                "page_info": {
                    "description": "分页信息",
                    "allOf": [
                        {
                            "$ref": "#/definitions/pagination.PageInfo"
                        }
                    ]
                }
            }
        },
        "admin.UserResponse": {
            "type": "object",
            "properties": {
                "create_time": {
                    "description": "创建时间",
                    "type": "string"
                },
                "id": {
                    "description": "ID",
                    "type": "integer"
                },
                "is_delete": {
                    "description": "是否删除：0否 1是",
                    "type": "integer"
                },
                "remark": {
                    "description": "备注",
                    "type": "string"
                },
                "role_type": {
                    "description": "角色：0普通用户 1管理员",
                    "type": "integer"
                },
                "status": {
                    "description": "角色：0停用 1启用",
                    "type": "integer"
                },
                "token": {
                    "description": "Token",
                    "type": "string"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "handlers.resp": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "pagination.PageInfo": {
            "type": "object",
            "properties": {
                "page": {
                    "description": "页号",
                    "type": "integer"
                },
                "page_size": {
                    "description": "每页大小",
                    "type": "integer"
                },
                "total_number": {
                    "description": "数据总条数",
                    "type": "integer"
                },
                "total_page": {
                    "description": "数据总页数",
                    "type": "integer"
                }
            }
        },
        "validators.UserLogin": {
            "type": "object",
            "required": [
                "login_type"
            ],
            "properties": {
                "code": {
                    "description": "第三方登录临时授权码code",
                    "type": "string"
                },
                "login_type": {
                    "description": "注册类型：0-账号密码 1-钉钉",
                    "type": "integer"
                },
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
