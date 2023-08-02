// +build ignore

// addins.go
package plugin

import (
	"strings"
	"sync/atomic"
	"syscall"
	"unsafe"
)

var (
	openHandleSeq  uint32
	freerdpClient2 = syscall.NewLazyDLL("freerdp-client2.dll")
	winpr2         = syscall.NewLazyDLL("winpr2.dll")
	freerdp2       = syscall.NewLazyDLL("freerdp2.dll")
)

var (
	virtualChannelEntry   = freerdpClient2.NewProc("VirtualChannelEntry")
	virtualChannelEntryEx = freerdpClient2.NewProc("VirtualChannelEntryEx")
)

func VirtualChannelEntryEx(ex *ChannelEntryPointsEx, pInitHandle interface{}) (err error) {
	r0, _, ec := virtualChannelEntryEx.Call(uintptr(unsafe.Pointer(ex)),
		uintptr(unsafe.Pointer(&pInitHandle)))
	if r0 == 0 {
		err = error(ec)
	}
	return
}

func ChannelsClientLoadEx(cs *rdpChannels) {
	var client ChannelClientData
	client.entryEx = VIRTUALCHANNELENTRYEX(VirtualChannelEntryEx)
	cs.clientDataList = append(cs.clientDataList, client)
	cs.clientDataCount++

	var init ChannelInitData
	init.channels = cs
	init.openDataMap = make(map[uint32]*ChannelOpenData)
	cs.initDataList = append(cs.initDataList, init)
	cs.initDataCount++

	ex := NewChannelEntryPointsEx()
	ex.PVirtualChannelInitEx = VIRTUALCHANNELINITEX(RdpVirtualChannelInitEx)
	ex.PVirtualChannelOpenEx = VIRTUALCHANNELOPENEX(RdpVirtualChannelOpenEx)
	ex.PVirtualChannelCloseEx = VIRTUALCHANNELCLOSEEX(RdpVirtualChannelCloseEx)
	ex.PVirtualChannelWriteEx = VIRTUALCHANNELWRITEEX(RdpVirtualChannelWriteEx)
	client.entryEx(ex, uintptr(unsafe.Pointer(&init)))
}

type ChannelClientData struct {
	//PVIRTUALCHANNELENTRY entry;
	entryEx VIRTUALCHANNELENTRYEX
	// pChannelInitEventProc *CHANNEL_INIT_EVENT_FN
	pChannelInitEventProcEx CHANNEL_INIT_EVENT_EX_FN
	pInitHandle             interface{}
	lpUserParam             interface{}
}
type ChannelOpenData struct {
	name        string
	OpenHandle  uint32
	options     uint32
	flags       int
	pInterface  interface{}
	channels    *rdpChannels
	lpUserParam interface{}
	//pChannelOpenEventProc *CHANNEL_OPEN_EVENT_FN
	pChannelOpenEventProcEx *CHANNEL_OPEN_EVENT_EX_FN
}
type ChannelInitData struct {
	channels    *rdpChannels
	pInterface  interface{}
	openDataMap map[uint32]*ChannelOpenData
}
type rdpChannels struct {
	clientDataCount int
	clientDataList  []ChannelClientData

	openDataCount int
	openDataList  []ChannelOpenData

	initDataCount int
	initDataList  []ChannelInitData

	/* control for entry into MyVirtualChannelInit */
	can_call_init bool

	/* true once freerdp_channels_post_connect is called */
	connected bool

	/* used for locating the channels for a given instance */
	//freerdp* instance;

	//wMessageQueue* queue;

	//DrdynvcClientContext* drdynvc;
	//CRITICAL_SECTION channelsLock;
}
type ChannelOpenEvent struct {
	Data             interface{}
	DataLength       uint32
	UserData         interface{}
	pChannelOpenData *ChannelOpenData
}

func RdpVirtualChannelInitEx(lpUserParam interface{}, clientContext interface{},
	pInitHandle interface{}, pChannel []ChannelDef,
	channelCount int, versionRequested uint32,
	pChannelInitEventProcEx CHANNEL_INIT_EVENT_EX_FN) uint {
	var (
		//rdpSettings* settings;
		pChannelInitData   *ChannelInitData
		pChannelClientData *ChannelClientData
		channels           *rdpChannels
	)

	if pInitHandle == nil {
		return CHANNEL_RC_BAD_INIT_HANDLE
	}

	if pChannel == nil {
		return CHANNEL_RC_BAD_CHANNEL
	}

	if (channelCount <= 0) || pChannelInitEventProcEx == nil {
		return CHANNEL_RC_INITIALIZATION_ERROR
	}

	pChannelInitData = pInitHandle.(*ChannelInitData)
	//WINPR_ASSERT(pChannelInitData);

	channels = pChannelInitData.channels
	//WINPR_ASSERT(channels);

	if !channels.can_call_init {
		return CHANNEL_RC_NOT_IN_VIRTUALCHANNELENTRY
	}

	if (channels.openDataCount + channelCount) > 30 {
		return CHANNEL_RC_TOO_MANY_CHANNELS
	}

	if channels.connected {
		return CHANNEL_RC_ALREADY_CONNECTED
	}

	if versionRequested != VIRTUAL_CHANNEL_VERSION_WIN2000 {
	}

	for i := range pChannel {
		pChannelDef := &pChannel[i]
		if getChannelOpenDataByName(channels, pChannelDef.Name) == nil {
			return CHANNEL_RC_BAD_CHANNEL
		}
	}

	pChannelClientData = &channels.clientDataList[channels.clientDataCount]
	pChannelClientData.pChannelInitEventProcEx = pChannelInitEventProcEx
	pChannelClientData.pInitHandle = pInitHandle
	pChannelClientData.lpUserParam = lpUserParam
	channels.clientDataCount++

	//WINPR_ASSERT(channels->instance);
	//WINPR_ASSERT(channels->instance->context);
	//settings = channels.instance.context.settings
	//WINPR_ASSERT(settings);

	for i := range pChannel {
		pChannelDef := &pChannel[i]
		var pChannelOpenData ChannelOpenData

		//WINPR_ASSERT(pChannelOpenData)

		pChannelOpenData.OpenHandle = atomic.AddUint32(&openHandleSeq, 1)
		pChannelOpenData.channels = channels
		pChannelOpenData.lpUserParam = lpUserParam
		if _, ok := pChannelInitData.openDataMap[pChannelOpenData.OpenHandle]; ok {
			return CHANNEL_RC_INITIALIZATION_ERROR
		}

		pChannelInitData.pInterface = clientContext

		pChannelOpenData.flags = 1
		pChannelOpenData.name = pChannelDef.Name
		pChannelOpenData.options = pChannelDef.Options
		pChannelInitData.openDataMap[pChannelOpenData.OpenHandle] = &pChannelOpenData
		channels.openDataList = append(channels.openDataList, pChannelOpenData)
		channels.openDataCount++
		/*
			if settings.ChannelCount < 30 {
				channel := freerdp_settings_get_pointer_array_writable(
					settings, FreeRDP_ChannelDefArray, settings.ChannelCount)
				channel.name = pChannelDef.Name
				channel.options = pChannelDef.Options
				settings.ChannelCount++
			}*/

		channels.openDataCount++
	}

	return CHANNEL_RC_OK
}
func getChannelOpenDataByName(channel *rdpChannels, name string) *ChannelOpenData {
	for _, v := range channel.openDataList {
		if strings.EqualFold(name, v.name) {
			return &v
		}
	}
	return nil
}
func RdpVirtualChannelOpenEx(pInitHandle interface{}, pOpenHandle *uint32, pChannelName string,
	pChannelOpenEventProcEx *CHANNEL_OPEN_EVENT_EX_FN) uint {
	pChannelInitData := pInitHandle.(*ChannelInitData)
	channels := pChannelInitData.channels
	pInterface := pChannelInitData.pInterface

	if pOpenHandle == nil {
		return CHANNEL_RC_BAD_CHANNEL_HANDLE
	}
	if pChannelOpenEventProcEx == nil {
		return CHANNEL_RC_BAD_PROC
	}

	if !channels.connected {
		return CHANNEL_RC_NOT_CONNECTED
	}

	pChannelOpenData := getChannelOpenDataByName(channels, pChannelName)

	if pChannelOpenData == nil {
		return CHANNEL_RC_UNKNOWN_CHANNEL_NAME
	}

	if pChannelOpenData.flags == 2 {
		return CHANNEL_RC_ALREADY_OPEN
	}

	pChannelOpenData.flags = 2 /* open */
	pChannelOpenData.pInterface = pInterface
	pChannelOpenData.pChannelOpenEventProcEx = pChannelOpenEventProcEx
	*pOpenHandle = pChannelOpenData.OpenHandle
	return CHANNEL_RC_OK
}
func RdpVirtualChannelCloseEx(pInitHandle interface{}, openHandle uint32) uint {
	if pInitHandle == nil {
		return CHANNEL_RC_BAD_INIT_HANDLE
	}
	pChannelInitData := pInitHandle.(*ChannelInitData)
	pChannelOpenData := pChannelInitData.openDataMap[openHandle]

	if pChannelOpenData == nil {
		return CHANNEL_RC_BAD_CHANNEL_HANDLE
	}

	if pChannelOpenData.flags != 2 {
		return CHANNEL_RC_NOT_OPEN
	}

	pChannelOpenData.flags = 0

	return CHANNEL_RC_OK
}
func RdpVirtualChannelWriteEx(pInitHandle interface{}, openHandle uint32,
	pData interface{}, dataLength uint32,
	pUserData interface{}) uint {

	//wMessage message;

	if pInitHandle == nil {
		return CHANNEL_RC_BAD_INIT_HANDLE
	}

	pChannelInitData := pInitHandle.(*ChannelInitData)
	channels := pChannelInitData.channels

	if channels == nil {
		return CHANNEL_RC_BAD_CHANNEL_HANDLE
	}

	pChannelOpenData := pChannelInitData.openDataMap[openHandle]
	if pChannelOpenData == nil {
		return CHANNEL_RC_BAD_CHANNEL_HANDLE
	}

	if !channels.connected {
		return CHANNEL_RC_NOT_CONNECTED
	}

	if pData == nil {
		return CHANNEL_RC_NULL_DATA
	}

	if dataLength == 0 {
		return CHANNEL_RC_ZERO_LENGTH
	}

	if pChannelOpenData.flags != 2 {
		return CHANNEL_RC_NOT_OPEN
	}

	pChannelOpenEvent := new(ChannelOpenEvent)

	if pChannelOpenEvent == nil {
		return CHANNEL_RC_NO_MEMORY

	}

	pChannelOpenEvent.Data = pData
	pChannelOpenEvent.DataLength = dataLength
	pChannelOpenEvent.UserData = pUserData
	pChannelOpenEvent.pChannelOpenData = pChannelOpenData
	/*message.context = channels;
	message.id = 0;
	message.wParam = pChannelOpenEvent;
	message.lParam = NULL;
	message.Free = channel_queue_message_free;

	if (!MessageQueue_Dispatch(channels->queue, &message))
	{
		free(pChannelOpenEvent);
		return CHANNEL_RC_NO_MEMORY;
	}*/

	return CHANNEL_RC_OK
}
