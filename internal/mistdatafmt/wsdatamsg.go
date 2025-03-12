package mistdatafmt

import (
	"encoding/json"
	"fmt"
)

/*
 * Base WebSocket message data format
 */
type WsMsgData struct {
	Event	string			`json:"event"`
	Channel	string			`json:"channel"`
	Data	string			`json:"data"`
	Detail	string			`json:"detail"`
}

/*
 * Client to Mist message data format 
 */
type WsMsgSubscribe struct {
	Subscribe	string		`json:"subscribe"`
}

/*
 * Client statistics message data format
 * For receiving data for:
 * 1. /sites/:site_id/stats/client
 */
type WsMsgClientStat struct {
	Mac		string		`json:"mac"`
	SiteId		string		`json:"site_id"`
	AssocTime	json.Number	`json:"assoc_time"`
	Family		string		`json:"family"`
	Model		string		`json:"model"`
	Os		string		`json:"os"`
	Manufacture	string		`json:"manufacture"`
	Bssid		string		`json:"bssid"`
	Username	string		`json:"username"`
	Hostname	string		`json:"hostname"`
	Ip		string		`json:"ip"`
	Ip6		string		`json:"ip6"`
	ApMac		string		`json:"ap_mac"`
	ApId		string		`json:"ap_id"`
	LastSeen	json.Number	`json:"last_seen"`
	Uptime		json.Number	`json:"uptime"`
	Ssid		string		`json:"ssid"`
	WlanId		string		`json:"wlan_id"`
	PskId		string		`json:"psk_id"`
	DualBand	bool		`json:"dual_band"`
	IsGuest		bool		`json:"is_guest"`
	KeyMgmt		string		`json:"key_mgmt"`
	Group		string		`json:"group"`
	Band		string		`json:"band"`
	Channel		json.Number	`json:"channel"`
	VlanId		string		`json:"vlan_id"`
	Proto		string		`json:"proto"`
	Rssi		json.Number	`json:"rssi"`
	Snr		json.Number	`json:"snr"`
	IdleTime	json.Number	`json:"idle_time"`
	TxRate		json.Number	`json:"tx_rate"`
	RxRate		json.Number	`json:"rx_rate"`
	TxPkts		json.Number	`json:"tx_pkts"`
	RxPkts		json.Number	`json:"rx_pkts"`
	TxBytes		json.Number	`json:"tx_bytes"`
	RxBytes		json.Number	`json:"rx_bytes"`
	TxRetries	json.Number	`json:"tx_retries"`
	RxRetries	json.Number	`json:"rx_retries"`
	TxBps		json.Number	`json:"tx_bps"`
	RxBps		json.Number	`json:"rx_bps"`
	MapId		string		`json:"map_id"`
	MapX		json.Number	`json:"x"`
	MapY		json.Number	`json:"y"`
	MapXM		json.Number	`json:"x_m"`
	MapYM		json.Number	`json:"y_m"`
	NumLocatingAps	json.Number	`json:"num_locating_aps"`
	Ttl		json.Number	`json:"_ttl"`
}

func (m *WsMsgClientStat) GetJsonKeyValue(key string) (interface{}, error) {
	switch key {
	case "mac":
		return m.Mac, nil
	case "site_id":
		return m.SiteId, nil
	case "assoc_time":
		return m.AssocTime, nil
	case "family":
		return m.Family, nil
	case "model":
		return m.Model, nil
	case "os":
		return m.Os, nil
	case "manufacture":
		return m.Manufacture, nil
	case "bssid":
		return m.Bssid, nil
	case "username":
		return m.Username, nil
	case "hostname":
		return m.Hostname, nil
	case "ip":
		return m.Ip, nil
	case "ip6":
		return m.Ip6, nil
	case "ap_mac":
		return m.ApMac, nil
	case "ap_id":
		return m.ApId, nil
	case "last_seen":
		return m.LastSeen, nil
	case "uptime":
		return m.Uptime, nil
	case "ssid":
		return m.Ssid, nil
	case "wlan_id":
		return m.WlanId, nil
	case "psk_id":
		return m.PskId, nil
	case "dual_band":
		return m.DualBand, nil
	case "is_guest":
		return m.IsGuest, nil
	case "key_mgmt":
		return m.KeyMgmt, nil
	case "group":
		return m.Group, nil
	case "band":
		return m.Band, nil
	case "channel":
		return m.Channel, nil
	case "vlan_id":
		return m.VlanId, nil
	case "proto":
		return m.Proto, nil
	case "rssi":
		return m.Rssi, nil
	case "snr":
		return m.Snr, nil
	case "idle_time":
		return m.IdleTime, nil
	case "tx_rate":
		return m.TxRate, nil
	case "rx_rate":
		return m.RxRate, nil
	case "tx_pkts":
		return m.TxPkts, nil
	case "rx_pkts":
		return m.RxPkts, nil
	case "tx_bytes":
		return m.TxBytes, nil
	case "rx_bytes":
		return m.RxBytes, nil
	case "tx_retries":
		return m.TxRetries, nil
	case "rx_retries":
		return m.RxRetries, nil
	case "tx_bps":
		return m.TxBps, nil
	case "rx_bps":
		return m.RxBps, nil
	case "map_id":
		return m.MapId, nil
	case "x":
		return m.MapX, nil
	case "y":
		return m.MapY, nil
	case "x_m":
		return m.MapXM, nil
	case "y_m":
		return m.MapYM, nil
	case "num_locating_aps":
		return m.NumLocatingAps, nil
	case "_ttl":
		return m.Ttl, nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (m *WsMsgClientStat) GetJsonKeyValueAsStr(key string) (string, error) {
	switch key {
	case "mac":
		return m.Mac, nil
	case "site_id":
		return m.SiteId, nil
	case "assoc_time":
		return string(m.AssocTime), nil
	case "family":
		return m.Family, nil
	case "model":
		return m.Model, nil
	case "os":
		return m.Os, nil
	case "manufacture":
		return m.Manufacture, nil
	case "bssid":
		return m.Bssid, nil
	case "username":
		return m.Username, nil
	case "hostname":
		return m.Hostname, nil
	case "ip":
		return m.Ip, nil
	case "ip6":
		return m.Ip6, nil
	case "ap_mac":
		return m.ApMac, nil
	case "ap_id":
		return m.ApId, nil
	case "last_seen":
		return string(m.LastSeen), nil
	case "uptime":
		return string(m.Uptime), nil
	case "ssid":
		return m.Ssid, nil
	case "wlan_id":
		return m.WlanId, nil
	case "psk_id":
		return m.PskId, nil
	case "dual_band":
		return fmt.Sprintf("%v", m.DualBand), nil
	case "is_guest":
		return fmt.Sprintf("%v", m.IsGuest), nil
	case "key_mgmt":
		return m.KeyMgmt, nil
	case "group":
		return m.Group, nil
	case "band":
		return m.Band, nil
	case "channel":
		return string(m.Channel), nil
	case "vlan_id":
		return m.VlanId, nil
	case "proto":
		return m.Proto, nil
	case "rssi":
		return string(m.Rssi), nil
	case "snr":
		return string(m.Snr), nil
	case "idle_time":
		return string(m.IdleTime), nil
	case "tx_rate":
		return string(m.TxRate), nil
	case "rx_rate":
		return string(m.RxRate), nil
	case "tx_pkts":
		return string(m.TxPkts), nil
	case "rx_pkts":
		return string(m.RxPkts), nil
	case "tx_bytes":
		return string(m.TxBytes), nil
	case "rx_bytes":
		return string(m.RxBytes), nil
	case "tx_retries":
		return string(m.TxRetries), nil
	case "rx_retries":
		return string(m.RxRetries), nil
	case "tx_bps":
		return string(m.TxBps), nil
	case "rx_bps":
		return string(m.RxBps), nil
	case "map_id":
		return m.MapId, nil
	case "x":
		return string(m.MapX), nil
	case "y":
		return string(m.MapY), nil
	case "x_m":
		return string(m.MapXM), nil
	case "y_m":
		return string(m.MapYM), nil
	case "num_locating_aps":
		return string(m.NumLocatingAps), nil
	case "_ttl":
		return string(m.Ttl), nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (m *WsMsgClientStat) GetJsonKeyValueAsFloat64(key string) (float64, error) {
	switch key {
	case "assoc_time":
		return m.AssocTime.Float64()
	case "last_seen":
		return m.LastSeen.Float64()
	case "uptime":
		return m.Uptime.Float64()
	case "channel":
		return m.Channel.Float64()
	case "rssi":
		return m.Rssi.Float64()
	case "snr":
		return m.Snr.Float64()
	case "idle_time":
		return m.IdleTime.Float64()
	case "tx_rate":
		return m.TxRate.Float64()
	case "rx_rate":
		return m.RxRate.Float64()
	case "tx_pkts":
		return m.TxPkts.Float64()
	case "rx_pkts":
		return m.RxPkts.Float64()
	case "tx_bytes":
		return m.TxBytes.Float64()
	case "rx_bytes":
		return m.RxBytes.Float64()
	case "tx_retries":
		return m.TxRetries.Float64()
	case "rx_retries":
		return m.RxRetries.Float64()
	case "tx_bps":
		return m.TxBps.Float64()
	case "rx_bps":
		return m.RxBps.Float64()
	case "x":
		return m.MapX.Float64()
	case "y":
		return m.MapY.Float64()
	case "x_m":
		return m.MapXM.Float64()
	case "y_m":
		return m.MapYM.Float64()
	case "num_locating_aps":
		return m.NumLocatingAps.Float64()
	case "_ttl":
		return m.Ttl.Float64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

func (m *WsMsgClientStat) GetJsonKeyValueAsInt64(key string) (int64, error) {
	switch key {
	case "assoc_time":
		return m.AssocTime.Int64()
	case "last_seen":
		return m.LastSeen.Int64()
	case "uptime":
		return m.Uptime.Int64()
	case "channel":
		return m.Channel.Int64()
	case "rssi":
		return m.Rssi.Int64()
	case "snr":
		return m.Snr.Int64()
	case "idle_time":
		return m.IdleTime.Int64()
	case "tx_rate":
		return m.TxRate.Int64()
	case "rx_rate":
		return m.RxRate.Int64()
	case "tx_pkts":
		return m.TxPkts.Int64()
	case "rx_pkts":
		return m.RxPkts.Int64()
	case "tx_bytes":
		return m.TxBytes.Int64()
	case "rx_bytes":
		return m.RxBytes.Int64()
	case "tx_retries":
		return m.TxRetries.Int64()
	case "rx_retries":
		return m.RxRetries.Int64()
	case "tx_bps":
		return m.TxBps.Int64()
	case "rx_bps":
		return m.RxBps.Int64()
	case "x":
		return m.MapX.Int64()
	case "y":
		return m.MapY.Int64()
	case "x_m":
		return m.MapXM.Int64()
	case "y_m":
		return m.MapYM.Int64()
	case "num_locating_aps":
		return m.NumLocatingAps.Int64()
	case "_ttl":
		return m.Ttl.Int64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}



/*
 * Wi-Fi client location message data format
 * For receiving data for:
 * 1. /sites/:site_id/stats/maps/:map_id/unconnected_clients
 * 2. /sites/:site_id/stats/maps/:map_id/clients
 */
type WsMsgMapClient struct {
	Mac		string		`json:"mac"`
	MapId		string		`json:"map_id"`
	MapX		json.Number	`json:"x"`
	MapY		json.Number	`json:"y"`
	MapXM		json.Number	`json:"x_m"`
	MapYM		json.Number	`json:"y_m"`
	NumLocatingAps	json.Number	`json:"num_locating_aps"`
	Rssi		json.Number	`json:"rssi"`
	Ttl		json.Number	`json:"_ttl"`

	// connected client only
	ConnectedAp	string		`json:"connected_ap"`

	// unconnected client only
	ApMac		string		`json:"ap_mac"` 
	Lastseen	json.Number	`json:"last_seen"`
}

func (m *WsMsgMapClient) GetJsonKeyValue(key string) (interface{}, error) {
	switch key {
	case "mac":
		return m.Mac, nil
	case "rssi":
		return m.Rssi, nil
	case "map_id":
		return m.MapId, nil
	case "x":
		return m.MapX, nil
	case "y":
		return m.MapY, nil
	case "x_m":
		return m.MapXM, nil
	case "y_m":
		return m.MapYM, nil
	case "num_locating_aps":
		return m.NumLocatingAps, nil
	case "connected_ap":
		return m.ConnectedAp, nil
	case "ap_mac":
		return m.ApMac, nil
	case "last_seen":
		return m.Lastseen, nil
	case "_ttl":
		return m.Ttl, nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (m *WsMsgMapClient) GetJsonKeyValueAsStr(key string) (string, error) {
	switch key {
	case "mac":
		return m.Mac, nil
	case "rssi":
		return string(m.Rssi), nil
	case "map_id":
		return m.MapId, nil
	case "x":
		return string(m.MapX), nil
	case "y":
		return string(m.MapY), nil
	case "x_m":
		return string(m.MapXM), nil
	case "y_m":
		return string(m.MapYM), nil
	case "num_locating_aps":
		return string(m.NumLocatingAps), nil
	case "connected_ap":
		return m.ConnectedAp, nil
	case "ap_mac":
		return m.ApMac, nil
	case "last_seen":
		return string(m.Lastseen), nil
	case "_ttl":
		return string(m.Ttl), nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (m *WsMsgMapClient) GetJsonKeyValueAsFloat64(key string) (float64, error) {
	switch key {
	case "rssi":
		return m.Rssi.Float64()
	case "x":
		return m.MapX.Float64()
	case "y":
		return m.MapY.Float64()
	case "x_m":
		return m.MapXM.Float64()
	case "y_m":
		return m.MapYM.Float64()
	case "num_locating_aps":
		return m.NumLocatingAps.Float64()
	case "last_seen":
		return m.Lastseen.Float64()
	case "_ttl":
		return m.Ttl.Float64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

func (m *WsMsgMapClient) GetJsonKeyValueAsInt64(key string) (int64, error) {
	switch key {
	case "rssi":
		return m.Rssi.Int64()
	case "x":
		return m.MapX.Int64()
	case "y":
		return m.MapY.Int64()
	case "x_m":
		return m.MapXM.Int64()
	case "y_m":
		return m.MapYM.Int64()
	case "num_locating_aps":
		return m.NumLocatingAps.Int64()
	case "last_seen":
		return m.Lastseen.Int64()
	case "_ttl":
		return m.Ttl.Int64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}


/*
 * BLE location on Map WebSocket Message
 * For receiving data for:
 * 1. /sites/:site_id/stats/maps/:map_id/discovered_assets
 * 2. /sites/:site_id/stats/maps/:map_id/assets
 */
type WsMsgMapBleAsset struct {
	Mac			string		`json:"mac"`
	CurrSite		string		`json:"curr_site"`
	MapId			string		`json:"map_id"`
	AssetId			string		`json:"id"`
	DevId			string		`json:"device_id"`
	By			string		`json:"by"`
	Name			string		`json:"name"`
	DevName			string		`json:"device_name"`

	Id			string		`json:"_id"`
	Manufacture		string		`json:"manufacture"`
	Temperature		json.Number	`json:"temperature"`
	BattVoltage		json.Number	`json:"battery_voltage"`
	ServiceData		WsMsgBleSvcData	`json:"service_data"`
	IbeaconUUID		string		`json:"ibeacon_uuid"`
	IbeaconMajor		json.Number	`json:"ibeacon_major"`
	IbeaconMinor		json.Number	`json:"ibeacon_minor"`
	EddystoneUIDNamespace	string		`json:"eddystone_uid_namespace"`
	EddystoneUIDInstance	string		`json:"eddystone_uid_instance"`
	EddystoneURL		string		`json:"eddystone_url"`

	MapX			json.Number	`json:"x"`
	MapY			json.Number	`json:"y"`
	MapXM			json.Number	`json:"x_m"`
	MapYM			json.Number	`json:"y_m"`
	Beam			json.Number	`json:"beam"`
	Rssi			json.Number	`json:"rssi"`
	ApMac			string		`json:"ap_mac"`
	Lastseen		json.Number	`json:"last_seen"`
	Timestamp		json.Number	`json:"_timestamp"`
	Ttl			json.Number	`json:"_ttl"`
}

func (m *WsMsgMapBleAsset) GetJsonKeyValue(key string) (interface{}, error) {
	switch key {
	case "mac":
		return m.Mac, nil
	case "curr_site":
		return m.CurrSite, nil
	case "map_id":
		return m.MapId, nil
	case "id":
		return m.AssetId, nil
	case "device_id":
		return m.DevId, nil
	case "by":
		return m.By, nil
	case "name":
		return m.Name, nil
	case "device_name":
		return m.DevName, nil
	case "_id":
		return m.Id, nil
	case "manufacture":
		return m.Manufacture, nil
	case "temperature":
		return m.Temperature, nil
	case "battery_voltage":
		return m.BattVoltage, nil
	case "service_data":
		return m.ServiceData, nil
	case "ibeacon_uuid":
		return m.IbeaconUUID, nil
	case "ibeacon_major":
		return m.IbeaconMajor, nil
	case "ibeacon_minor":
		return m.IbeaconMinor, nil
	case "eddystone_uid_namespace":
		return m.EddystoneUIDNamespace, nil
	case "eddystone_uid_instance":
		return m.EddystoneUIDInstance, nil
	case "eddystone_url":
		return m.EddystoneURL, nil
	case "x":
		return m.MapX, nil
	case "y":
		return m.MapY, nil
	case "x_m":
		return m.MapXM, nil
	case "y_m":
		return m.MapYM, nil
	case "beam":
		return m.Beam, nil
	case "rssi":
		return m.Rssi, nil
	case "ap_mac":
		return m.ApMac, nil
	case "last_seen":
		return m.Lastseen, nil
	case "_timestamp":
		return m.Timestamp, nil
	case "_ttl":
		return m.Ttl, nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}


func (m *WsMsgMapBleAsset) GetJsonKeyValueAsStr(key string) (string, error) {
	switch key {
	case "mac":
		return m.Mac, nil
	case "curr_site":
		return m.CurrSite, nil
	case "map_id":
		return m.MapId, nil
	case "id":
		return m.AssetId, nil
	case "device_id":
		return m.DevId, nil
	case "by":
		return m.By, nil
	case "name":
		return m.Name, nil
	case "device_name":
		return m.DevName, nil
	case "_id":
		return m.Id, nil
	case "manufacture":
		return m.Manufacture, nil
	case "temperature":
		return string(m.Temperature), nil
	case "battery_voltage":
		return string(m.BattVoltage), nil
	case "ibeacon_uuid":
		return m.IbeaconUUID, nil
	case "ibeacon_major":
		return string(m.IbeaconMajor), nil
	case "ibeacon_minor":
		return string(m.IbeaconMinor), nil
	case "eddystone_uid_namespace":
		return m.EddystoneUIDNamespace, nil
	case "eddystone_uid_instance":
		return m.EddystoneUIDInstance, nil
	case "eddystone_url":
		return m.EddystoneURL, nil
	case "x":
		return string(m.MapX), nil
	case "y":
		return string(m.MapY), nil
	case "x_m":
		return string(m.MapXM), nil
	case "y_m":
		return string(m.MapYM), nil
	case "beam":
		return string(m.Beam), nil
	case "rssi":
		return string(m.Rssi), nil
	case "ap_mac":
		return m.ApMac, nil
	case "last_seen":
		return string(m.Lastseen), nil
	case "_timestamp":
		return string(m.Timestamp), nil
	case "_ttl":
		return string(m.Ttl), nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (m *WsMsgMapBleAsset) GetJsonKeyValueAsFloat64(key string) (float64, error) {
	switch key {
	case "temperature":
		return m.Temperature.Float64()
	case "battery_voltage":
		return m.BattVoltage.Float64()
	case "ibeacon_major":
		return m.IbeaconMajor.Float64()
	case "ibeacon_minor":
		return m.IbeaconMinor.Float64()
	case "x":
		return m.MapX.Float64()
	case "y":
		return m.MapY.Float64()
	case "x_m":
		return m.MapXM.Float64()
	case "y_m":
		return m.MapYM.Float64()
	case "beam":
		return m.Beam.Float64()
	case "rssi":
		return m.Rssi.Float64()
	case "last_seen":
		return m.Lastseen.Float64()
	case "_timestamp":
		return m.Timestamp.Float64()
	case "_ttl":
		return m.Ttl.Float64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

func (m *WsMsgMapBleAsset) GetJsonKeyValueAsInt64(key string) (int64, error) {
	switch key {
	case "temperature":
		return m.Temperature.Int64()
	case "battery_voltage":
		return m.BattVoltage.Int64()
	case "ibeacon_major":
		return m.IbeaconMajor.Int64()
	case "ibeacon_minor":
		return m.IbeaconMinor.Int64()
	case "x":
		return m.MapX.Int64()
	case "y":
		return m.MapY.Int64()
	case "x_m":
		return m.MapXM.Int64()
	case "y_m":
		return m.MapYM.Int64()
	case "beam":
		return m.Beam.Int64()
	case "rssi":
		return m.Rssi.Int64()
	case "last_seen":
		return m.Lastseen.Int64()
	case "_timestamp":
		return m.Timestamp.Int64()
	case "_ttl":
		return m.Ttl.Int64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

type WsMsgBleSvcData struct {
	UUID			string		`json:"uuid"`
	Data			string		`json:"data"`
	RxCnt			json.Number	`json:"rx_cnt"`
	LastRxTime		json.Number	`json:"last_rx_time"`
}

func (m *WsMsgBleSvcData) GetJsonKeyValue(key string) (interface{}, error) {
	switch key {
	case "uuid":
		return m.UUID, nil
	case "data":
		return m.Data, nil
	case "rx_cnt":
		return m.RxCnt, nil
	case "last_rx_time":
		return m.LastRxTime, nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}


func (m *WsMsgBleSvcData) GetJsonKeyValueAsStr(key string) (string, error) {
	switch key {
	case "uuid":
		return m.UUID, nil
	case "data":
		return m.Data, nil
	case "rx_cnt":
		return string(m.RxCnt), nil
	case "last_rx_time":
		return string(m.LastRxTime), nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (m *WsMsgBleSvcData) GetJsonKeyValueAsFloat64(key string) (float64, error) {
	switch key {
	case "rx_cnt":
		return m.RxCnt.Float64()
	case "last_rx_time":
		return m.LastRxTime.Float64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

func (m *WsMsgBleSvcData) GetJsonKeyValueAsInt64(key string) (int64, error) {
	switch key {
	case "rx_cnt":
		return m.RxCnt.Int64()
	case "last_rx_time":
		return m.LastRxTime.Int64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

