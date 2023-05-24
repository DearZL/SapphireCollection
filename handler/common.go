package handler

import (
	"P/model"
	"P/resp"
	"P/service"
	"github.com/gin-gonic/gin"
	"log"
)

type CommonHandler struct {
	CommonSrvI service.CommonServiceInterface
}

func (h *CommonHandler) Upload(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "上传错误！",
		Data: nil,
	}
	fileType := c.Request.Header.Get("FileType")
	if fileType != "Icon" && fileType != "Commodity" {
		c.AbortWithStatus(500)
		log.Println("已拒绝未携带文件类型请求")
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err.Error())
		entity.SetCodeAndMsg(500, "参数错误")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	files := form.File["files"]
	if files == nil {
		log.Println("文件集为空")
		entity.SetCodeAndMsg(500, "参数错误")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	location := ""
	if fileType == "Icon" {
		location = "./file/icon/"
	} else {
		location = "./file/commodity/"
	}
	fs := &model.Files{}
	for _, file := range files {
		f := model.File{}
		f.Set(file.Filename, location, fileType)
		fs.Files = append(fs.Files, f)
		err = c.SaveUploadedFile(file, location+file.Filename)
		if err != nil {
			log.Println(err.Error())
			c.JSON(200, gin.H{"entity": entity})
			return
		}
	}
	err = h.CommonSrvI.Upload(fs.Files)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	entity.SetEntityAndHeaderToken(c)
	entity.SetCodeAndMsg(200, "上传成功！")
	c.JSON(200, gin.H{"entity": entity})
	return
}

func (h *CommonHandler) Download(c *gin.Context) {
	entity := resp.EntityA{
		Code: 500,
		Msg:  "下载错误！",
		Data: nil,
	}
	filename := c.Param("filename")
	if filename == "" {
		entity.SetCodeAndMsg(500, "参数错误！")
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	file, err := h.CommonSrvI.Download(filename)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	//inline;attachment
	way := c.Request.Header.Get("Content-Disposition")
	c.Header("Content-Disposition", way)
	c.Header("Content-Type", "image")
	c.File(file.Location + filename)
	return
}
