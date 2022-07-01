package assembly

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"http-spy-forward/logger"
	"http-spy-forward/models"
	"http-spy-forward/vars"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/sirupsen/logrus"
)

type httpStreamFactory struct{}

type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

// NewHttpStreamFactory returns a new httpStreamFactory.
func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hStream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hStream.run()
	return &hStream.r
}

// run assembles the http stream.
// It runs in a separate goroutine.
func (h *httpStream) run() {
	// 创建一个bufio.Reader用于读取重组完成的http请求
	buf := bufio.NewReader(&h.r)
	// 读取http请求
	for {
		// 读取gopacket中的http请求
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			return
		} else if err == nil {
			defer func() {
				_ = req.Body.Close()
			}()
			clientIp, dstIp := SplitNet2Ips(h.net)
			srcPort, dstPort := Transport2Ports(h.transport)

			httpReq, _ := models.NewHttpReq(req, clientIp, dstIp, dstPort)
			// 发送到服务器
			go func(req *models.HttpReq) {
				// 获取http流量
				reqInfo := fmt.Sprintf("%v:%v -> %v(%v:%v), %v, %v, %v, %v", httpReq.Client, srcPort, httpReq.Host, httpReq.Ip,
					httpReq.Port, httpReq.Method, httpReq.URL, httpReq.Header, httpReq.ReqParameters)

				if logger.Log.Logger.Level == logrus.DebugLevel {
					logger.Log.Infof("reqInfo: %v", reqInfo)
				}
				// 将http请求包信息传值到 vars.HttpReqChan
				vars.DataChan <- httpReq

			}(httpReq)
		}
	}
}

func PrettyPrint(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(v)
		return
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(v)
		return
	}

	fmt.Println(out.String())
}

// 从gopacket.Flow中提取出客户端IP和服务端IP
func SplitNet2Ips(net gopacket.Flow) (client, host string) {
	ips := strings.Split(net.String(), "->")
	if len(ips) > 1 {
		client = ips[0]
		host = ips[1]
	}
	return client, host
}

// 从gopacket.Flow中提取出客户端和服务端端口
func Transport2Ports(transport gopacket.Flow) (src, dst string) {
	ports := strings.Split(transport.String(), "->")
	if len(ports) > 1 {
		src = ports[0]
		dst = ports[1]
	}
	return src, dst
}

// 判断是是不是本地发送的请求
func CheckSelfHtml(host string, req *models.HttpReq) (ret bool) {
	if host == strings.Split(req.Host, ":")[0] {
		ret = true
	}
	return ret
}

func ProcessPackets(packets chan gopacket.Packet) {
	// 自定义一个streamFactory
	streamFactory := &httpStreamFactory{}
	// 将streamFactory注册到tcpassembly.StreamPool
	// tcpassembly.StreamPool是并发安全的，可以在多个goroutine中使用
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	// 创建一个tcpassembly.Assembler
	assembler := tcpassembly.NewAssembler(streamPool)

	ticker := time.NewTicker(time.Minute)
	// 循环读取packets中的数据
	for {
		select {
		case packet := <-packets:

			if packet == nil {
				return
			}
			// fmt.Println(packet)
			// 网络层+传输层+TCP层
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				continue
			}
			// 读取传输层
			tcp := packet.TransportLayer().(*layers.TCP)
			// 创建一个新的流
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker.C:
			// 定时刷新缓存中的数据
			assembler.FlushOlderThan(time.Now().Add(time.Second * -20))
		}
	}
}
