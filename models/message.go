// message
package models

type OneWeight struct {
	Id           int
	DevId        int32
	WType        int32
	Weight       int32  // 单点重量、
	WetDate      string `orm:"size(100)"` // 单点重量的获取日期时间、
	Gps          string `orm:"size(100)"` // GPS信息、
	LicensePlate string `orm:"size(100)"` // 车辆号牌信息（或本机信息）、
	Duty         string `orm:"size(100)"` // 值班员（或司机信息）、
	UpDate       string `orm:"size(100)"` // 发送的实时日期时间.
}

type Gps struct {
	Id        int
	DevId     int32   //设备编号
	Longitude float64 // 经度
	Latitude  float64 // 纬度
	Ns        uint8   // 南北值为,'n'或's'
	Ew        uint8   // 东西,'e'或'w'
	Date      int64
}
