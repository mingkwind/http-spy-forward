package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"http-spy-forward/logger"
	"http-spy-forward/modules/assembly"
	"http-spy-forward/modules/forward"
	"http-spy-forward/vars"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	snapshotLen int32 = 10240
	promiscuous bool  = true
	err         error
	timeout     time.Duration = pcap.BlockForever
	handle      *pcap.Handle

	DebugMode  bool
	DeviceName string
	filter     = ""
)

func Start(ctx *cli.Context) error {
	if ctx.IsSet("device") {
		DeviceName = ctx.String("device")
	} else {
		logger.Log.Fatal("--device value, -i value is required")
	}

	if ctx.IsSet("debug") {
		DebugMode = ctx.Bool("debug")
	}
	if DebugMode {
		logger.Log.Logger.Level = logrus.DebugLevel
	}

	if ctx.IsSet("length") {
		snapshotLen = int32(ctx.Int("len"))
	}

	if ctx.IsSet("url") {
		vars.TargetURL = ctx.String("url")
	} else {
		logger.Log.Fatal("--url value, -u value is required")
	}
	// Open device
	handle, err = pcap.OpenLive(DeviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	if ctx.IsSet("filter") {
		filter = ctx.String("filter")
	}

	err = handle.SetBPFFilter(filter)
	if err != nil {
		return err
	}

	// 开启50个线程进行包转发
	for i := 0; i < 50; i++ {
		go func() {
			for req := range vars.DataChan {
				forward.Forward(req)
			}
		}()
	}

	// 开始抓包
	logger.Log.Info("start sniffing...")
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	assembly.ProcessPackets(packetSource.Packets())

	return err
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
