package unmarshaller

import (
	"net"

	"gitlab.x.lan/yunshan/droplet-libs/app"
	"gitlab.x.lan/yunshan/droplet-libs/datatype"
	"gitlab.x.lan/yunshan/droplet-libs/grpc"
	"gitlab.x.lan/yunshan/droplet-libs/utils"
	"gitlab.x.lan/yunshan/droplet-libs/zerodoc"
	"gitlab.x.lan/yunshan/droplet/roze/msg"
	pf "gitlab.x.lan/yunshan/droplet/roze/platformdata"
)

const (
	EdgeCode    = zerodoc.IPPath | zerodoc.L3EpcIDPath
	MainAddCode = zerodoc.RegionID | zerodoc.HostID | zerodoc.L3Device | zerodoc.SubnetID | zerodoc.PodNodeID | zerodoc.AZID | zerodoc.PodGroupID | zerodoc.PodNSID | zerodoc.PodID | zerodoc.PodClusterID
	EdgeAddCode = zerodoc.RegionIDPath | zerodoc.HostIDPath | zerodoc.L3DevicePath | zerodoc.SubnetIDPath | zerodoc.PodNodeIDPath | zerodoc.AZIDPath | zerodoc.PodGroupIDPath | zerodoc.PodNSIDPath | zerodoc.PodIDPath | zerodoc.PodClusterIDPath
)

func releaseRozeDocument(rd *msg.RozeDocument) {
	rd.Document = nil
	msg.ReleaseRozeDocument(rd)
}

func DocToRozeDocuments(doc *app.Document) *msg.RozeDocument {
	rd := msg.AcquireRozeDocument()
	rd.Document = doc
	t := doc.Tag.(*zerodoc.Tag)
	t.SetID("") // 由于需要修改Tag增删Field，清空ID避免字段脏

	// vtap_acl 分钟级数据不用填充
	if doc.Meter.ID() == zerodoc.PACKET_ID &&
		t.DatabaseSuffixID() == 1 { // 只有acl后缀
		return rd
	}

	var info, info1 *grpc.Info
	myRegionID := uint16(pf.PlatformData.QueryRegionID())
	if t.Code&EdgeCode == EdgeCode {
		t.Code |= EdgeAddCode

		if t.L3EpcID == datatype.EPC_FROM_INTERNET && t.L3EpcID1 == datatype.EPC_FROM_INTERNET {
			return rd
		}

		// 当MAC/MAC1非0时，通过MAC来获取资源信息
		if t.MAC != 0 && t.MAC1 != 0 {
			info, info1 = pf.PlatformData.QueryMacInfosPair(t.MAC, t.MAC1)
		} else if t.MAC != 0 {
			info = pf.PlatformData.QueryMacInfo(t.MAC)
			if t.IsIPv6 != 0 {
				info1 = pf.PlatformData.QueryIPV6Infos(t.L3EpcID1, t.IP61)
			} else {
				info1 = pf.PlatformData.QueryIPV4Infos(t.L3EpcID1, t.IP1)
			}
		} else if t.MAC1 != 0 {
			if t.IsIPv6 != 0 {
				info = pf.PlatformData.QueryIPV6Infos(t.L3EpcID, t.IP6)
			} else {
				info = pf.PlatformData.QueryIPV4Infos(t.L3EpcID, t.IP)
			}
			info1 = pf.PlatformData.QueryMacInfo(t.MAC1)
		} else if t.IsIPv6 != 0 {
			info, info1 = pf.PlatformData.QueryIPV6InfosPair(t.L3EpcID, t.IP6, t.L3EpcID1, t.IP61)
		} else {
			info, info1 = pf.PlatformData.QueryIPV4InfosPair(t.L3EpcID, t.IP, t.L3EpcID1, t.IP1)
		}
		if info1 != nil {
			t.RegionID1 = uint16(info1.RegionID)
			t.HostID1 = uint16(info1.HostID)
			t.L3DeviceID1 = uint16(info1.DeviceID)
			t.L3DeviceType1 = zerodoc.DeviceType(info1.DeviceType)
			t.SubnetID1 = uint16(info1.SubnetID)
			t.PodNodeID1 = uint16(info1.PodNodeID)
			t.PodNSID1 = uint16(info1.PodNSID)
			t.AZID1 = uint16(info1.AZID)
			t.PodGroupID1 = int16(info1.PodGroupID)
			t.PodID1 = uint16(info1.PodID)
			t.PodClusterID1 = uint16(info1.PodClusterID)
			if info == nil {
				var ip0 net.IP
				if t.IsIPv6 != 0 {
					ip0 = t.IP6
				} else {
					ip0 = utils.IpFromUint32(t.IP)
				}
				// 当0侧是组播ip时，使用1侧的region_id,subnet_id,az_id来填充
				if ip0.IsMulticast() {
					t.RegionID = t.RegionID1
					t.SubnetID = t.SubnetID1
					t.AZID = t.AZID1
				}
			}
			if myRegionID != 0 && t.RegionID1 != 0 {
				if t.TAPSide == zerodoc.Server && t.RegionID1 != myRegionID { // 对于双端 的统计值，需要去掉 tap_side 对应的一侧与自身region_id 不匹配的内容。
					releaseRozeDocument(rd)
					return nil
				}
			}
		}
	} else {
		t.Code |= MainAddCode
		if t.L3EpcID == datatype.EPC_FROM_INTERNET {
			return rd
		}

		if t.MAC != 0 {
			info = pf.PlatformData.QueryMacInfo(t.MAC)
		} else if t.IsIPv6 != 0 {
			info = pf.PlatformData.QueryIPV6Infos(t.L3EpcID, t.IP6)
		} else {
			info = pf.PlatformData.QueryIPV4Infos(t.L3EpcID, t.IP)
		}
	}

	if info != nil {
		t.RegionID = uint16(info.RegionID)
		t.HostID = uint16(info.HostID)
		t.L3DeviceID = uint16(info.DeviceID)
		t.L3DeviceType = zerodoc.DeviceType(info.DeviceType)
		t.SubnetID = uint16(info.SubnetID)
		t.PodNodeID = uint16(info.PodNodeID)
		t.PodNSID = uint16(info.PodNSID)
		t.AZID = uint16(info.AZID)
		t.PodGroupID = int16(info.PodGroupID)
		t.PodID = uint16(info.PodID)
		t.PodClusterID = uint16(info.PodClusterID)
		if info1 == nil && (t.Code&EdgeCode == EdgeCode) {
			var ip1 net.IP
			if t.IsIPv6 != 0 {
				ip1 = t.IP61
			} else {
				ip1 = utils.IpFromUint32(t.IP1)
			}
			// 当1侧是组播ip时，使用0侧的region_id,subnet_id,az_id来填充
			if ip1.IsMulticast() {
				t.RegionID1 = t.RegionID
				t.SubnetID1 = t.SubnetID
				t.AZID1 = t.AZID
			}
		}

		if myRegionID != 0 && t.RegionID != 0 {
			if t.Code&EdgeCode == EdgeCode { // 对于双端 的统计值，需要去掉 tap_side 对应的一侧与自身region_id 不匹配的内容。
				if t.TAPSide == zerodoc.Client && t.RegionID != myRegionID {
					releaseRozeDocument(rd)
					return nil
				}
			} else { // 对于单端的统计值，需要去掉与自身region_id不匹配的内容
				if t.RegionID != myRegionID {
					releaseRozeDocument(rd)
					return nil
				}
			}
		}
	}

	return rd
}
