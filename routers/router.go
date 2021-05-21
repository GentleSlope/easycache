package routers

import (
	"easycache/middleware/jwt"
	"easycache/routers/api/v2"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "easycache/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"easycache/pkg/export"
	"easycache/pkg/qrcode"
	"easycache/pkg/upload"
	"easycache/routers/api"
	"easycache/routers/api/v1"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.POST("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/auth", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		r.POST("/tags/export", v1.ExportTag)
		//导入标签
		r.POST("/tags/import", v1.ImportTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		//生成文章海报
		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}

	apiV2 := r.Group("/api/v2")
	//apiV2.Use(jwt.JWT())
	{
		// 任务相关接口
		apiV2.GET("/tasks/:replica_id", v2.GetTaskByReplicaId)
		apiV2.POST("/tasks/between", v2.GetTaskInfo)
		apiV2.POST("/tasks/update", v2.UpdateDeployTask)
		apiV2.POST("/tasks/update2", v2.UpdateNotDeployTask)
		apiV2.GET("/nodes/get", v2.GetNodeTaskByReplicaId)
		apiV2.GET("/deploy/get/:uuid/:ip", v2.GetDeployTask)
		apiV2.DELETE("/tasks/delete", v2.DeleteTask)

		//实例接口
		apiV2.PUT("/instance/insert/:redis_user_id", v2.InsertInstance)
		apiV2.GET("/instance/all", v2.GetInstances)
		apiV2.GET("/instance/start", v2.StartInstance)
		apiV2.GET("/instance/stop", v2.StopInstance)
		apiV2.DELETE("/instance/delete", v2.DeleteInstance)

		// 资源接口
		apiV2.POST("/resource/insert", v2.InsertResource)
		apiV2.GET("/resource/all", v2.GetResources)
		apiV2.DELETE("/resource/delete", v2.DeleteResource)

		// 集群接口
		apiV2.POST("/cluster/create", v2.CreateCluster)
		apiV2.GET("/cluster/all", v2.GetClusters)
		apiV2.DELETE("/cluster/delete", v2.DeleteCluster)
		apiV2.GET("/cluster-instance/get", v2.GetClusterInstanceById)
		apiV2.GET("/cluster/scale", v2.ScaleCluster)
		apiV2.DELETE("/cluster-node/delete", v2.DeleteNode)

		// 监控接口
		apiV2.POST("/monitor/insert", v2.InsertMonitorInfo)
		apiV2.GET("/monitor/normal/get", v2.GetMonitorInfo)

		// 配置接口
		apiV2.GET("/config/get", v2.GetNodeConfigByReplicaId)
		apiV2.POST("/config/update", v2.UpdateNodeConfig)

		// 告警接口
		apiV2.POST("/alarm/insert", v2.InsertAlarmInfo)
		apiV2.GET("/alarm/all", v2.GetAlarmInfos)
		apiV2.DELETE("/alarm/delete", v2.DeleteAlarmInfos)

		// 分类接口
		apiV2.GET("/user/get", v2.GerRedisUsers)
		apiV2.POST("/user/insert", v2.InsertRedisUser)
		apiV2.DELETE("/user/delete", v2.DeleteRedisUser)

		// 日志接口
		apiV2.POST("/log/insert", v2.InsertLogs)
		apiV2.GET("/log/all", v2.GetLogs)
		apiV2.DELETE("/log/delete", v2.DeleteLogs)
	}

	return r
}
