package DataManipulation

import (
	"fmt"
	"github.com/weilaihui/fdfs_client"
)

func UploadByBuffer(filebuffer*[]byte,fileExt string) (fileId string,err error) {
	fd,err:=fdfs_client.NewFdfsClient("/home/parallels/goProject/src/uhome/UHomeWeb/conf/client.conf")

	if err!=nil {
		fmt.Println("创建fdfs句炳失败",err)
		return
	}
	fd_rsp,err:=fd.UploadByBuffer(*filebuffer,fileExt)

	if err!=nil {
		fmt.Println("fdfs上传失败",err)
		return
	}
	fileId=fd_rsp.RemoteFileId
	return
}