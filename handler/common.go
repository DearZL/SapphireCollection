package handler

import (
	"P/model"
	"P/resp"
	"P/service"
	"errors"
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
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	files := form.File["files"]
	if files == nil {
		err := errors.New("文件集为空")
		log.Println(err.Error())
		c.JSON(200, gin.H{"entity": entity})
		return
	}
	fs := &model.Files{}
	for _, file := range files {
		f := model.File{}
		f.Set(file.Filename, "./file/")
		fs.Files = append(fs.Files, f)
		err := c.SaveUploadedFile(file, "./file/"+file.Filename)
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
	token := c.Request.Header.Get("Authorization")
	entity.SetMsg("上传成功！")
	entity.SetCode(200)
	entity.Data = token
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
		entity.SetMsg("参数错误！")
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
