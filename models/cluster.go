package models

import (
	"easycache/pkg/logger"
	"github.com/jinzhu/gorm"
	"time"
)

type ClusterForm struct {
	Instances               []ClusterInstanceForm   `json:"instances"`
	Password                string                  `json:"password"`
	InstanceConfigExtension InstanceConfigExtension `json:"instance_config_extension"`
	RedisUserId             int64                   `json:"redis_user_id"`
	RedisUserName           string                  `json:"redis_user_name"`
	ClusterName             string                  `json:"cluster_name"`
	SlotMode                string                  `json:"slot_mode"`
}

type ClusterInstanceForm struct {
	ID         int64  `json:"id"`
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	Role       string `json:"role"`
	Slots      string `json:"slots"`
	MasterIp   string `json:"master_ip"`
	MasterPort int    `json:"master_port"`
	Version    string `json:"version"`
}
type InstanceConfigExtension struct {
	Config  NodeConfig `json:"config"`
	Restart bool       `json:"restart"`
}
type Cluster struct {
	Id          int64     `json:"id"`
	RedisUserId int64     `json:"redis_user_id"`
	Name        string    `json:"name"`
	SlotMode    string    `json:"slot_mode"`
	NodesNum    int64     `json:"nodes_num"`
	MasterNum   int64     `json:"master_num"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CLusterDeploy struct {
	Cluster    Cluster           `json:"cluster"`
	ClusterPwd string            `json:"cluster_pwd"`
	Instances  []ClusterInstance `json:"instances"`
}

func UpdateClusterWithLastId(cluster *Cluster) (err error) {
	err = db.Where(&cluster).Order("id desc").First(&cluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

func InsertCluster(cluster *Cluster) (err error) {
	if err := db.Create(&cluster).Error; err != nil {
		return err
	}
	err = UpdateClusterWithLastId(cluster)
	if err != nil {
		logger.Error("models UpdateClusterWithLastId err :", err)
		return err
	}
	return nil
}

func GetAllClusters() (clusters []Cluster, err error) {
	err = db.Find(&clusters).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return clusters, err
	}

	return clusters, nil
}

func DeleteCluster(id int64) error {
	if err := db.Where("id = ?", id).Delete(&Cluster{}).Error; err != nil {
		return err
	}
	return nil
}
