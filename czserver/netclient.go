// netclient
package czserver

import (
	models "cz400/models"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/mahonia"
	_ "github.com/mattn/go-sqlite3"
)

type NetClient struct {
	Con      net.Conn
	parser   ProtoParser
	Device   DevInfo
	uploader *Uploader
}

type Param struct {
	ptype      uint8
	Value      interface{} //当前设备的真实值
	writeValue interface{} //需要写入的值.
	write      bool
	read       bool
}
type DevInfo struct {
	DeviceId  uint16
	Version   string
	GpsReport uint8
	DevReport uint8
	K         float32
	Plate     string
	//OnDateTime string
	timeStamp time.Time
	parats    time.Time
	UnixTime  int64
	IpAddr    string
	//Paras     map[string]Param
}
type DevInfoList []DevInfo

var timeoutS int = 120
var clientList map[string]*NetClient = make(map[string]*NetClient, 100)
var mutex sync.Mutex

func CreateClient(con net.Conn) (client *NetClient) {
	fmt.Println("new Client")

	cli := &NetClient{}
	cli.Con = con
	cli.parser = ProtoParser{}
	cli.parser.Data = make([]byte, 0, 512)
	cli.parser.waitHead = true
	cli.Device.IpAddr = con.RemoteAddr().String()
	cli.uploader = &Uploader{}
	//cli.Device.OnDateTime = time.Now().String()
	cli.Device.timeStamp = time.Now()
	mutex.Lock()
	defer mutex.Unlock()
	clientList[con.RemoteAddr().String()] = cli
	return cli
}
func UploadDevice(id uint16, filename string, action string) error {
	client := GetClientById(id)
	if client == nil {

		return errors.New("can not find id")
	}
	return nil
}
func RemoveClient(con net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(clientList, con.RemoteAddr().String())
}
func handleGps(ipaddr string, gps *GpsDef) {
	mutex.Lock()
	defer mutex.Unlock()

	msg := new(models.Gps)
	msg.Longitude = gps.Longitude
	msg.Latitude = gps.Latitude
	msg.Ns = gps.Ns
	msg.Ew = gps.Ew

	msg.Date = time.Now().Unix()
	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库
	_, err := o.Insert(msg)
	if err != nil {
		fmt.Println(err)
	}

}

//func RequestReadParam(ipaddr string, key string) {
//	mutex.Lock()
//	defer mutex.Unlock()
//	fmt.Println("RequestReadParam")
//	if c, ok := clientList[ipaddr]; ok { //存在}

//		p := clientList[ipaddr].Device.Paras[key]
//		p.read = true
//		clientList[ipaddr].Device.Paras[key] = p
//		readParam(c, p.ptype)
//	}
//}

//处理参数
func handleParam(ipaddr string, msg Message) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := clientList[ipaddr]; ok { //存在}

		//		key := fmt.Sprintf("key%d", msg.Head.Cmd)
		//		p := clientList[ipaddr].Device.Paras[key]
		//		if msg.Head.Oper == OPER_READ {
		//			//更新读取标志为不读取，更新最新读取值.

		//			p.Value = msg.Val
		//			p.read = false

		//		} else if msg.Head.Oper == OPER_WRITE {
		//			//写入成功了，更新读取值为写入值.不用再去读取一次了.
		//			p.Value = p.writeValue

		//			//更新写入标志，为不再写入.
		//			p.write = false

		//		}
		//		clientList[ipaddr].Device.Paras[key] = p
		v, ok := msg.Val.(float32)
		if ok {
			clientList[ipaddr].Device.K = v
		}

	}

}
func handleUpload(ipaddr string, err byte) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := clientList[ipaddr]; ok { //存在}
		fmt.Println("handle online")
		clientList[ipaddr].uploader.ParseAck(err, clientList[ipaddr].Con)
	}
}
func UpdateParam(info *DevInfo) {
	mutex.Lock()
	defer mutex.Unlock()
	if c, ok := clientList[info.IpAddr]; ok { //存在}
		writeParam(c, CMD_K, int(info.K*1000000))
		c.Device.K = info.K
	}
}
func handleOnline(ipaddr string, dev *DevicePara) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := clientList[ipaddr]; ok { //存在}
		fmt.Println("handle online")
		clientList[ipaddr].Device.DeviceId = dev.DeviceId
		clientList[ipaddr].Device.DevReport = dev.DevReport
		clientList[ipaddr].Device.GpsReport = dev.GpsReport
		//		p := Param{}
		//		p.ptype = CMD_K
		//		p.read = false
		//		p.write = false
		//		p.Value = 0
		//		p.writeValue = 0
		//		clientList[ipaddr].Device.Paras = make(map[string]Param)
		//		key := fmt.Sprintf("key%d", CMD_K)
		//		clientList[ipaddr].Device.Paras[key] = p
		enc := mahonia.NewDecoder("GBK")
		src := string(dev.Plate[:])

		clientList[ipaddr].Device.Plate = enc.ConvertString(src)
		clientList[ipaddr].Device.Version = fmt.Sprintf("v%d.%d.%d", (dev.Version>>16)&0xff, (dev.Version>>8)&0xff, dev.Version&0xff)
		fmt.Println("deviceId", clientList[ipaddr].Device.DeviceId)
		readParam(clientList[ipaddr], CMD_K)
	}
}
func resetTimeout(con net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := clientList[con.RemoteAddr().String()]; ok { //存在}
		fmt.Println("reset timeout")
		clientList[con.RemoteAddr().String()].Device.timeStamp = time.Now()
	}
}
func resetParamTimeout(con net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := clientList[con.RemoteAddr().String()]; ok { //存在}
		fmt.Println("reset param timeout")
		clientList[con.RemoteAddr().String()].Device.parats = time.Now()
	}
}
func SendPacket(client *NetClient, cmd uint8, oper uint8, data []byte) bool {
	head := MsgHead{}
	head.Cmd = cmd
	head.DevId = client.Device.DeviceId
	head.Dir = 0x5A
	head.Oper = oper
	head.Len = uint16(len(data))
	buf := MsgHeadToBytes(head)
	buf = append(buf, data...)
	crc16 := Reentrent_CRC16(buf, uint32(len(buf)))
	buf = append(buf, uint8(crc16&0xff))
	buf = append(buf, uint8(crc16>>8&0xff))
	fmt.Println("---------------")
	fmt.Println(head.DevId)
	fmt.Printf(" %02x ", buf)
	_, err := client.Con.Write(buf)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func writeParam(client *NetClient, cmd uint8, v int) {

	data := IntToBytes(v)

	SendPacket(client, cmd, OPER_WRITE, data)
}

func readParam(client *NetClient, cmd uint8) {
	fmt.Println("readPara")
	var data = make([]byte, 0)
	SendPacket(client, cmd, OPER_READ, data)
}

////处理参数更新
//func handleParamUpdate() {
//	mutex.Lock()
//	defer mutex.Unlock()
//	fmt.Println("handleParamUpdate")
//	for _, client := range clientList {

//		diff := time.Now().Sub(client.Device.parats)
//		s := time.Duration(20) * time.Second
//		if diff > s {
//			for _, v := range client.Device.Paras {
//				if v.write {
//					//先处理写指令
//					writeParam(client, v.ptype, v.writeValue)
//					//复位计数器，下一个20秒再处理下一个读写指令
//					client.Device.parats = time.Now()
//				} else if v.read {
//					readParam(client, v.ptype)
//					client.Device.parats = time.Now()
//				}
//			}
//		}

//	}
//}
func handleTimeout() {
	mutex.Lock()
	defer mutex.Unlock()
	for _, value := range clientList {

		diff := time.Now().Sub(value.Device.timeStamp)
		s := time.Duration(timeoutS) * time.Second

		if diff > s {
			fmt.Println(value.Device.IpAddr, " timeout")
			value.Con.Close()
		}

	}

}
func handleMsg(msg Message, con net.Conn) {
	fmt.Println("cmd=", msg.Head.Cmd, "id=", msg.Head.DevId)
	resetTimeout(con)
	switch msg.Head.Cmd {
	case CMD_DEV2HOST_ONE_WEIGHT:
		//var p PointWet

		p, ok := msg.Val.(*PointWet)
		if !ok {
			fmt.Println("convt PointWet failed", p)
			return
		}
		insertOneWeight(msg.Head, p)

	case CMD_DEV2HOST_ALL_WEIGHT:
		fallthrough
	case CMD_DEV2HOST_WATER_WEIGHT:

		p, ok := msg.Val.(*CommWeight)
		if !ok {
			fmt.Println("convt CommWeight failed", p)
			return
		}
		insertCommonWeight(msg.Head, p)

	case CMD_DEV_ONLINE:
		d, ok := msg.Val.(*DevicePara)
		if !ok {
			fmt.Println("convt DevicePara failed", d)
			return
		}
		handleOnline(con.RemoteAddr().String(), d)
	case CMD_DEV2HOST_GPS:
		d, ok := msg.Val.(*GpsDef)
		if !ok {
			fmt.Println("convt gps failed", d)
			return
		}
		handleGps(con.RemoteAddr().String(), d)
	case CMD_DEV2HOST_HEART:
	case CMD_UPDATE:
		d, ok := msg.Val.(byte)
		if !ok {
			fmt.Println("convt gps failed", d)
			return
		}

		handleUpload(con.RemoteAddr().String(), d)
	case CMD_K:

		handleParam(con.RemoteAddr().String(), msg)
	default:
		fmt.Println("unkown cmd")
	}
}
func (cli *NetClient) Handle(data []byte, n int) (err bool) {

	msgList := cli.parser.Prase(data, n)
	fmt.Println("Handle msg", msgList)
	for i, v := range msgList {
		fmt.Println("handle msg", i)
		handleMsg(v, cli.Con)
	}
	return true
}
func GetClient() DevInfoList {
	mutex.Lock()
	defer mutex.Unlock()
	//infos := make([]DevInfo, 0, 30)
	infos := DevInfoList{}
	for _, v := range clientList {
		//fmt.Println(k, v)
		v.Device.UnixTime = v.Device.timeStamp.Unix() * 1000

		infos = append(infos, v.Device)

	}
	//fmt.Println(infos)
	return infos

}
func GetClientById(id uint16) *NetClient {
	mutex.Lock()
	defer mutex.Unlock()
	//infos := make([]DevInfo, 0, 30)

	for _, v := range clientList {
		//fmt.Println(k, v)
		if v.Device.DeviceId == id {
			return v
		}

	}
	//fmt.Println(infos)
	return nil

}
func fmtDate(dt DateDef) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", 2000+int(dt.Year), dt.Month, dt.Day, dt.Hour, dt.Min, dt.Sec)
}
func fmtGps(gps GpsDef) string {
	return fmt.Sprintf("%.6f,%.6f,%c,%c", gps.Latitude, gps.Longitude, gps.Ew, gps.Ns)
}
func insertOneWeight(head MsgHead, pwt *PointWet) {
	msg := new(models.OneWeight)
	enc := mahonia.NewDecoder("GBK")
	src := string(pwt.Plate[:])
	msg.WType = 1
	msg.Weight = pwt.Wet
	msg.LicensePlate = enc.ConvertString(src)
	msg.DevId = int32(head.DevId)
	src = string(pwt.Duty[:])
	msg.Duty = enc.ConvertString(src)

	msg.UpDate = fmtDate(pwt.UpDate)
	msg.WetDate = fmtDate(pwt.Wdate)
	msg.Gps = fmtGps(pwt.Gps)

	o := orm.NewOrm()
	err := o.Read(msg, "WetDate")
	if err == orm.ErrNoRows {
		o.Using("default") // 默认使用 default，你可以指定为其他数据库
		_, err := o.Insert(msg)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		//存在重复的称重时间
		fmt.Println("exist dup wetdate")
	}

}
func insertCommonWeight(head MsgHead, pwt *CommWeight) {
	msg := new(models.OneWeight)
	enc := mahonia.NewDecoder("GBK")
	src := string(pwt.Plate[:])
	msg.Weight = pwt.Wet
	msg.WType = int32(head.Cmd)
	msg.DevId = int32(head.DevId)
	msg.LicensePlate = enc.ConvertString(src)
	msg.Gps = fmtGps(pwt.Gps)
	msg.UpDate = fmtDate(pwt.UpDate)
	msg.WetDate = fmtDate(pwt.UpDate)
	o := orm.NewOrm()

	o.Using("default") // 默认使用 default，你可以指定为其他数据库
	_, err := o.Insert(msg)
	if err != nil {
		fmt.Println(err)
	}

}

//func paramUpdateTimeout(input chan bool) {
//	t1 := time.NewTimer(time.Second * 5)
//	//	t2 := time.NewTimer(time.Second * 10)
//	var msg bool = false
//	for {
//		select {
//		case msg = <-input:
//			println(msg)
//			if msg {
//				fmt.Println("exit param timeout")
//				break
//			}

//		case <-t1.C:
//			handleParamUpdate()

//			t1.Reset(time.Second * 5)
//		}
//	}
//}

func onlineTimeout(input chan bool) {
	t1 := time.NewTimer(time.Second * 5)
	//	t2 := time.NewTimer(time.Second * 10)
	var msg bool = false
	for {
		select {
		case msg = <-input:
			println(msg)
			if msg {
				fmt.Println("exit online timeout")
				break
			}

		case <-t1.C:
			//println("5s timer")
			handleTimeout()

			t1.Reset(time.Second * 5)

			//		case <-t2.C:
			//			println("10s timer")
			//			t2.Reset(time.Second * 10)
		}
	}
}

var quit chan bool
var Cfg = beego.AppConfig

func init() {

	timeoutS, _ = Cfg.Int("timeout")

	fmt.Println("timeout = ", timeoutS)
	go onlineTimeout(quit)
	//go paramUpdateTimeout(quit)
}
