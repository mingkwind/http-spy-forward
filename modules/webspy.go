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
	"github.com/spf13/cobra"
)

var (
	snapshotLen int32 = 10240
	promiscuous bool  = true
	err         error
	timeout     time.Duration = pcap.BlockForever
	handle      *pcap.Handle

	DebugMode bool
	Iface     string
	filter    = ""
)

func Start(cmd *cobra.Command, args []string) {
	Iface, _ = cmd.Flags().GetString("iface")
	DebugMode, _ = cmd.Flags().GetBool("debug")
	if DebugMode {
		logger.Log.Logger.Level = logrus.DebugLevel
	}
	snapshotLen, _ = cmd.Flags().GetInt32("length")
	vars.TargetURL, _ = cmd.Flags().GetString("url")
	filter, _ = cmd.Flags().GetString("filter")

	// Open device
	handle, err = pcap.OpenLive(Iface, snapshotLen, promiscuous, timeout)
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer handle.Close()

	err = handle.SetBPFFilter(filter)
	if err != nil {
		logger.Log.Fatal(err)
		return
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
