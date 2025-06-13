package main

import "github.com/gin-gonic/gin"

var Routes = map[string]map[string]gin.HandlerFunc{
	"/api/v1/ticket/create": {},
	"/api/v1/ticket/list":   {},
	"/api/v1/ticket/update": {},
}
