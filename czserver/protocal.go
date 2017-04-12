// protocal
package czserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	//"utils"

	"cz400/models"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type MsgHead struct {
	DevId uint16
	Dir   uint8
	Cmd   uint8
	Oper  uint8
	Len   uint16
}
type GpsDef struct {
	Longitude float64 // 经度
	Latitude  float64 // 纬度
	Ns        uint8   // 南北值为,'n'或's'
	Ew        uint8   // 东西,'e'或'w'
}

type DateDef struct {
	Year  uint8 // 当前年减去2000,如2016年，year实际保存16。
	Month uint8
	Day   uint8
	Hour  uint8
	Min   uint8
	Sec   uint8
}

const (
	LICENSE_LEN = 10
	DUTY_LEN    = 16
)

// B组内容-总重量
type CommWeight struct {
	Wet    int32             // 单点重量、
	Plate  [LICENSE_LEN]byte // 车辆号牌信息（或本机信息）、
	Gps    GpsDef            // GPS信息、
	UpDate DateDef           // 发送的实时日期时间.

}

// A组内容-单点重量
type PointWet struct {
	Wet    int32             // 单点重量、
	Wdate  DateDef           // 单点重量的获取日期时间、
	Gps    GpsDef            // GPS信息、
	Plate  [LICENSE_LEN]byte // 车辆号牌信息（或本机信息）、
	Duty   [DUTY_LEN]byte    // 值班员（或司机信息）、
	UpDate DateDef           // 发送的实时日期时间.
}

type DevicePara struct {
	DeviceId  uint16
	Version   uint32
	GpsReport uint8
	DevReport uint8
	Plate     [LICENSE_LEN]byte
}

func (h *MsgHead) Init(d []byte) {

	err := binary.Read(bytes.NewReader(d), binary.LittleEndian, h)
	if err != nil {

	}
	fmt.Printf("%v\n", *h) // {1 2 3 4 5 6}

}

type ProtoParser struct {
	Data     []byte
	Header   MsgHead
	waitHead bool
}
type Message struct {
	Head MsgHead
	Val  interface{}
}

const (
	CMD_DEV2HOST_ONE_WEIGHT   = 1
	CMD_DEV2HOST_ALL_WEIGHT   = 2
	CMD_DEV2HOST_WATER_WEIGHT = 3
	CMD_DEV2HOST_GPS          = 4
	CMD_DEV2HOST_DEVINFO      = 5
	CMD_DEV2HOST_HEART        = 6
	CMD_UPDATE                = 7
	CMD_RESET                 = 8
	CMD_VER                   = 9
	CMD_REALTIME_WEIGHT       = 10
	CMD_GPS                   = 11
	CMD_GPS_REPORT_TIME       = 12 //轨迹上发时间间隔
	CMD_DEV_REPORT_TIME       = 13 //设备运行情况上发时间间隔
	CMD_DEV_ONLINE            = 14 //上报设备信息.
	CMD_DEV2HOST_MV           = 15 //传感器电压值(只读)
	CMD_K                     = 16 //总K系数(可读写)
)
const (
	OPER_REPORT = 1
	OPER_READ   = 2

	OPER_WRITE = 4
)

func unSerial(data []byte, n uint16, msg interface{}) bool {

	r := bytes.NewReader(data[:n])
	err := binary.Read(r, binary.LittleEndian, msg)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Print("%v", msg)
	return true
}
func cvtFloat(d []byte) float32 {
	var v int32
	binary.Read(bytes.NewReader(d), binary.LittleEndian, &v)
	fmt.Println(v)

	return float32(v) / 1000000
}
func parseMsg(head MsgHead, d []byte, n uint16) Message {
	fmt.Println("handle msgtype ", head.Cmd)
	msg := Message{Head: head, Val: nil}
	switch head.Cmd {
	case CMD_DEV2HOST_ONE_WEIGHT:

		pwt := &PointWet{}

		if unSerial(d, n, pwt) {
			msg.Val = pwt
		}

	case CMD_DEV2HOST_ALL_WEIGHT:
		fallthrough
	case CMD_DEV2HOST_WATER_WEIGHT:
		pwt := &CommWeight{}
		if unSerial(d, n, pwt) {
			msg.Val = pwt
		}
	case CMD_DEV_ONLINE:
		dev := &DevicePara{}
		if unSerial(d, n, dev) {
			msg.Val = dev
		}
	case CMD_DEV2HOST_GPS:
		gps := &GpsDef{}
		if unSerial(d, n, gps) {
			msg.Val = gps
		}
	case CMD_DEV2HOST_HEART:
	case CMD_UPDATE:
		msg.Val = d[0]
	case CMD_K:
		if msg.Head.Oper == OPER_READ {
			msg.Val = cvtFloat(d)
		} else if msg.Head.Oper == OPER_WRITE {
			msg.Val = d[0]
		}

	}
	return msg
}

//分析数据协议.
func (p *ProtoParser) Prase(data []byte, n int) []Message {
	fmt.Printf("recv %02x %d %d\n", data, len(p.Data), n)
	fmt.Println(data)
	for i := 0; i < n; i++ {
		p.Data = append(p.Data, data[i])
	}

	var msgList = []Message{}
	for {
		fmt.Println("len", len(p.Data))
		if len(p.Data) <= 0 {
			break
		}
		if p.waitHead == true {
			//fmt.Println("find head")
			if len(p.Data) >= 7 {
				fmt.Println("find head ok")
				p.Header.Init(p.Data[:7])
				//p.Data = p.Data[7:]
				p.waitHead = false

			} else {
				fmt.Println("< 7")
				break
			}
		} else {
			//fmt.Println("find data len", p.Header.Len)
			var size int = int(p.Header.Len + 9)
			if len(p.Data) < size {
				break
			}
			var crc16 uint16 = Reentrent_CRC16(p.Data, uint32(p.Header.Len+7))
			var crc16_data uint16 = BytesToUint16(p.Data[p.Header.Len+7:])

			p.waitHead = true
			if crc16 != crc16_data {
				p.Data = p.Data[p.Header.Len+9:]
				fmt.Printf("crc error : crc1=%d,crc2=%d\n", crc16, crc16_data)

			}
			p.Data = p.Data[7 : 7+p.Header.Len]
			msg := parseMsg(p.Header, p.Data, p.Header.Len)

			msgList = append(msgList, msg)

			p.Data = p.Data[p.Header.Len:]
		}
	}
	return msgList
}

func init() {
	// 需要在init中注册定义的model
	fmt.Println("init sqlite3")
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "database/orm_test.db")
	orm.RegisterModel(new(models.OneWeight), new(models.Gps))
	orm.RunSyncdb("default", false, true)

}
