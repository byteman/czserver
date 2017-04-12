// update.go
package czserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const (
	UPDATE_IDEL = iota //0
	UPDATE_START
	UPDATE_DATA
	UPDATE_END
)
const (
	UPDATE_PACKET_SIZE  = 512
	OPER_HOST_WRITE_DEV = 0x4
	MAX_RETRY_COUNT     = 3
)

type fileStartDef struct {
	mode      int      // -1表示开始文件传输
	file_type uint32   // 0普通文件，1升级文件
	file_size uint32   // 传输的文件大小
	file_name [20]byte // 使用8.3文件名规则,字符串保存，空余填零
}

type Uploader struct {
	m_filebuffer    []byte
	m_reader        *os.File
	m_file_type     int //1 upload file
	m_packet_index  int
	m_total_packet  int
	m_state         int
	m_file_size     int
	m_file_name     [20]byte
	m_dev_id        uint16
	con             net.Conn
	m_timeout_retry int
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

//整形转换成字节
func MsgHeadToBytes(h MsgHead) []byte {

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, h)
	return bytesBuffer.Bytes()
}

func (this *Uploader) SendPacket(data []byte) bool {
	head := MsgHead{}
	head.Cmd = CMD_UPDATE
	head.DevId = this.m_dev_id
	head.Dir = 0x5A
	head.Oper = OPER_HOST_WRITE_DEV
	buf := MsgHeadToBytes(head)
	buf = append(buf, data...)
	crc16 := Reentrent_CRC16(buf, uint32(len(buf)))
	buf = append(buf, uint8(crc16&0xff))
	buf = append(buf, uint8(crc16>>8&0xff))
	_, err := this.con.Write(buf)

	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//拷贝字符串到数组.
func memcpy(dst, src []byte, size int) {
	for i := 0; i < size; i++ {
		dst[i] = src[i]
	}
	return
}
func (this *Uploader) SendStartPacket() bool {
	this.m_state = UPDATE_START

	var fileHead fileStartDef
	c := []byte("release.bin")
	memcpy(fileHead.file_name[:], c, len(c))
	fileHead.file_size = uint32(this.m_file_size) //文件大小
	fileHead.file_type = uint32(this.m_file_type) //升级文件类型
	fileHead.mode = -1

	this.m_packet_index = 0
	return true
}
func (this *Uploader) Upload(id uint16, filepath string, con net.Conn) bool {
	if !checkFileIsExist(filepath) {
		return false
	}
	size, err := FileSize(filepath)
	if err != nil {
		fmt.Println("failed to getsize")
		return false
	}
	this.m_file_size = int(size)
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("failed to open")
		return false
	}

	this.m_reader = file
	this.m_file_type = 1
	this.m_total_packet = (this.m_file_size + UPDATE_PACKET_SIZE - 1) / UPDATE_PACKET_SIZE
	this.m_packet_index = 0
	this.m_dev_id = id
	this.m_file_name[0] = 'u'
	this.m_timeout_retry = 0
	this.con = con
	this.SendStartPacket()
	return true
}
func (this *Uploader) isUpdateTimeout() bool {
	if this.m_timeout_retry > MAX_RETRY_COUNT {

		this.m_state = UPDATE_IDEL
		return true
	}
	return false
}
func (this *Uploader) updateComplete() {
	this.m_state = UPDATE_IDEL
	this.m_packet_index = 0
}
func (this *Uploader) sendUpdateStopRequest() {
	// qDebug() << "sendUpdateStopRequest";
	if this.isUpdateTimeout() {
		//qDebug() << "sendUpdateStartRequest timeout";
		return
	}
	var data = make([]byte, 4)
	this.m_state = UPDATE_END
	var oper = -2
	o := IntToBytes(oper)
	data = append(data, o...)

	this.SendPacket(data)
}
func (this *Uploader) sendUpdateData(index int) {
	pos := index * UPDATE_PACKET_SIZE
	//qDebug() << "sendUpdateData index=" << index << "total = " << m_total_packet;
	if this.isUpdateTimeout() {
		// qDebug() << "sendUpdateData timeout packet" << m_packet_index;
		return
	}
	if pos >= this.m_file_size {
		//qDebug() << "sendUpdateData complete";
		//已经读到最后长度了.
		this.sendUpdateStopRequest()
		return
	}
	this.m_packet_index = index

}
func (this *Uploader) ParseAck(err byte, con net.Conn) {
	switch this.m_state {
	case UPDATE_START:
		this.sendUpdateData(0)
	case UPDATE_DATA:
		this.m_packet_index++
		this.sendUpdateData(this.m_packet_index)
	case UPDATE_END:
		this.updateComplete()

	}
}
