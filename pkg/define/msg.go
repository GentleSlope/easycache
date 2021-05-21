package define

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "请求参数错误",
	ERROR_EXIST_TAG:                 "已存在该标签名称",
	ERROR_EXIST_TAG_FAIL:            "获取已存在标签失败",
	ERROR_NOT_EXIST_TAG:             "该标签不存在",
	ERROR_GET_TAGS_FAIL:             "获取所有标签失败",
	ERROR_COUNT_TAG_FAIL:            "统计标签失败",
	ERROR_ADD_TAG_FAIL:              "新增标签失败",
	ERROR_EDIT_TAG_FAIL:             "修改标签失败",
	ERROR_DELETE_TAG_FAIL:           "删除标签失败",
	ERROR_EXPORT_TAG_FAIL:           "导出标签失败",
	ERROR_IMPORT_TAG_FAIL:           "导入标签失败",
	ERROR_NOT_EXIST_ARTICLE:         "该文章不存在",
	ERROR_ADD_ARTICLE_FAIL:          "新增文章失败",
	ERROR_DELETE_ARTICLE_FAIL:       "删除文章失败",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:  "检查文章是否存在失败",
	ERROR_EDIT_ARTICLE_FAIL:         "修改文章失败",
	ERROR_COUNT_ARTICLE_FAIL:        "统计文章失败",
	ERROR_GET_ARTICLES_FAIL:         "获取多个文章失败",
	ERROR_GET_ARTICLE_FAIL:          "获取单个文章失败",
	ERROR_GEN_ARTICLE_POSTER_FAIL:   "生成文章海报失败",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",

	ErrorGetTasksFail: "查询任务失败",
	ErrorNotExistTask: "无可下发任务",
	ErrorUpdateTask:   "更新任务失败",
	ErrorGetRedisUser: "无法找到redis用户",
	ErrorDeleteTask:   "删除任务失败",

	ErrorInsertInstance: "插入新节点失败",
	ErrorGetInstance:    "查询实例记录失败",
	ErrorDeleteInstance: "删除实例记录失败",
	ErrorStartInstance:  "启动实例失败",
	ErrorStopInstance:   "停止实例失败",

	ErrorInsertResource: "插入资源信息失败",
	ErrorAllResource:    "获取所有资源信息错误",
	ErrorGetResource:    "获取单个资源信息失败",
	ErrorDeleteResource: "删除单个资源信息失败",

	ErrorCreateCluster: "创建集群失败",

	ErrorInsertInfo:   "插入信息失败",
	ErrorGetConfig:    "获取配置信息失败",
	ErrorUpdateConfig: "更改配置文件失败",

	ErrorInsertAlarm:  "插入警告失败",
	ErrorGetAlarms:    "获取警告失败",
	ErrorDeleteAlarms: "删除警告失败",

	ErrorAllCluster:          "获取集群列表失败",
	ErrorGetClusterInstances: "获取集群实例失败",
	ErrorDeleteCluster:       "删除集群失败",

	ErrorAllUsers:   "获取分组失败",
	ErrorInsertUser: "插入分组失败",
	ErrorDeleteUser: "删除分组失败",

	ErrorInsertLog:  "插入警告失败",
	ErrorGetLogs:    "获取警告失败",
	ErrorDeleteLogs: "删除警告失败",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
