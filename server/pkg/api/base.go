package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Data any  `json:"data,omitempty"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	// Status int64  `json:"status"`
	Msg   string `json:"msg"`
	Error string `json:"error,omitempty"`
}

// Err 通用错误处理
func Err(msg string, err error) Response {
	res := Response{
		Meta: Meta{
			Msg:   msg,
			Error: "",
		},
	}
	if err != nil {
		res.Meta.Error = err.Error()
	}
	return res
}

// DBErr 数据库操作失败
func DBErr(msg string, err error) Response {
	if msg == "" {
		msg = "数据库操作失败"
	}
	return Err(msg, err)
}

// ErrorResponse 返回错误消息
func ErrorResponse(err error) Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		var errorMsgs []string
		for _, e := range ve {
			field := e.Field()
			tag := e.Tag()
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s%s", field, tag))
			return Err(
				strings.Join(errorMsgs, ", "),
				err,
			)
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return Err("JSON类型不匹配", err)
	}

	return Err("参数错误", err)
}

// PageInfo 分页请求
type PageInfo struct {
	Page     int `json:"page" form:"page"`           // 页码
	PageSize int `json:"page_size" form:"page_size"` // 每页大小
}

// PageResult 带分页的返回
type PageResult struct {
	Data     any   `json:"data,omitempty"`
	Total    int64 `json:"total,omitempty"`
	Page     int   `json:"page,omitempty"`
	PageSize int   `json:"page_size,omitempty"`
	Meta     Meta  `json:"meta"`
}

// IdsReq ID多选请求格式
type IdsReq struct {
	Ids []uint `json:"ids" form:"ids" binding:"required"`
}

// IdReq ID单选请求格式
type IdReq struct {
	Id uint `json:"id" form:"id" binding:"required"`
}

type GetPagingByIdReq struct {
	Id uint `json:"id" form:"id"`
	PageInfo
}

type GetPagingMustByIdReq struct {
	Id uint `json:"id" form:"id" binding:"required"`
	PageInfo
}

type GetPagingMustByIdsReq struct {
	Ids []uint `json:"ids" form:"ids" binding:"required"`
	PageInfo
}

type GetPagingByIdsReq struct {
	Ids []uint `json:"ids" form:"ids"`
	PageInfo
}

// type GetPagingByIdsReq struct {
type SearchIdsStringReq struct {
	Ids    []uint `json:"ids" form:"ids"`
	String string `json:"string" form:"string"`
	PageInfo
}

// IdReq ID+String请求格式
type SearchIdStringReq struct {
	Id     uint   `json:"id" form:"id"`
	String string `json:"string" form:"string"`
	PageInfo
}
