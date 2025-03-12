package mistdatafmt

import (
	"encoding/json"
	"fmt"
)

/*
 * Maps API call data format (partial)
 * /sites/:site_id/maps
 */
type ApiDataMapEntry struct {
	Name			string		`json:"name"`
	WidthM			json.Number	`json:"width_m"`
	HeightM			json.Number	`json:"height_m"`
	Width			json.Number	`json:"width"`
	Height			json.Number	`json:"height"`
	PPM			json.Number	`json:"ppm"`
	Type			string		`json:"image"`
	Orientation		json.Number	`json:"orientation"`
	OccupancyLimit		json.Number	`json:"occupancy_limit"`
	Locked			bool		`json:"locked"`
	UseAutoOrientation	bool		`json:"use_auto_orientation"`
	UseAutoPlacement	bool		`json:"use_auto_placement"`
	Id			string		`json:"id"`
	SiteId			string		`json:"site_id"`
	OrgId			string		`json:"org_id"`
	CreatedTime		json.Number	`json:"created_time"`
	ModifiedTime		json.Number	`json:"modified_time"`
	Url			string		`json:"url"`
	ThumbnailUrl		string		`json:"thumbnail_url"`

	// wallpath, sitesurvey_path, wayfinding_path are skipped..
}

func (d *ApiDataMapEntry) GetJsonKeyValue(key string) (interface{}, error) {
	switch key {
	case "name":
		return d.Name, nil
	case "type":
		return d.Name, nil
	case "id":
		return d.Id, nil
	case "site_id":
		return d.SiteId, nil
	case "org_id":
		return d.OrgId, nil
	case "url":
		return d.Url, nil
	case "thumbnail_url":
		return d.ThumbnailUrl, nil
	case "locked":
		return d.Locked, nil
	case "use_auto_orientation":
		return d.UseAutoOrientation, nil
	case "use_auto_placement":
		return d.UseAutoPlacement, nil
	case "width_m":
		return d.WidthM, nil
	case "height_m":
		return d.HeightM, nil
	case "width":
		return d.Width, nil
	case "height":
		return d.Height, nil
	case "ppm":
		return d.PPM, nil
	case "orientation":
		return d.Orientation, nil
	case "occupancy_limit":
		return d.OccupancyLimit, nil
	case "created_time":
		return d.CreatedTime, nil
	case "modified_time":
		return d.ModifiedTime, nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}



func (d *ApiDataMapEntry) GetJsonKeyValueAsStr(key string) (string, error) {
	switch key {
	case "name":
		return d.Name, nil
	case "type":
		return d.Name, nil
	case "id":
		return d.Id, nil
	case "site_id":
		return d.SiteId, nil
	case "org_id":
		return d.OrgId, nil
	case "url":
		return d.Url, nil
	case "thumbnail_url":
		return d.ThumbnailUrl, nil
	case "locked":
		return fmt.Sprintf("%v", d.Locked), nil
	case "use_auto_orientation":
		return fmt.Sprintf("%v", d.UseAutoOrientation), nil
	case "use_auto_placement":
		return fmt.Sprintf("%v", d.UseAutoPlacement), nil
	case "width_m":
		return string(d.WidthM), nil
	case "height_m":
		return string(d.HeightM), nil
	case "width":
		return string(d.Width), nil
	case "height":
		return string(d.Height), nil
	case "ppm":
		return string(d.PPM), nil
	case "orientation":
		return string(d.Orientation), nil
	case "occupancy_limit":
		return string(d.OccupancyLimit), nil
	case "created_time":
		return string(d.CreatedTime), nil
	case "modified_time":
		return string(d.ModifiedTime), nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (d *ApiDataMapEntry) GetJsonKeyValueAsFloat64(key string) (float64, error) {
	switch key {
	case "width_m":
		return d.WidthM.Float64()
	case "height_m":
		return d.HeightM.Float64()
	case "width":
		return d.Width.Float64()
	case "height":
		return d.Height.Float64()
	case "ppm":
		return d.PPM.Float64()
	case "orientation":
		return d.Orientation.Float64()
	case "occupancy_limit":
		return d.OccupancyLimit.Float64()
	case "created_time":
		return d.CreatedTime.Float64()
	case "modified_time":
		return d.ModifiedTime.Float64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

func (d *ApiDataMapEntry) GetJsonKeyValueAsInt64(key string) (int64, error) {
	switch key {
	case "width_m":
		return d.WidthM.Int64()
	case "height_m":
		return d.HeightM.Int64()
	case "width":
		return d.Width.Int64()
	case "height":
		return d.Height.Int64()
	case "ppm":
		return d.PPM.Int64()
	case "orientation":
		return d.Orientation.Int64()
	case "occupancy_limit":
		return d.OccupancyLimit.Int64()
	case "created_time":
		return d.CreatedTime.Int64()
	case "modified_time":
		return d.ModifiedTime.Int64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

/*
 * Zones API call data format
 * /sites/:site_id/zones
 */
type ApiDataZoneEntry struct {
	Name			string			`json:"name"`
	OccupancyLimit		json.Number		`json:"occupancy_limit"`
	Id			string			`json:"id"`
	MapId			string			`json:"map_id"`
	SiteId			string			`json:"site_id"`
	OrgId			string			`json:"org_id"`
	CreatedTime		json.Number		`json:"created_time"`
	ModifiedTime		json.Number		`json:"modified_time"`
	
	Vertices		[]ApiDataZoneVertice	`json:"vertices"`
	VerticesM		[]ApiDataZoneVertice	`json:"vertices_m"`
}

type ApiDataZoneVertice struct {
	X			json.Number		`json:"x"`
	Y			json.Number		`json:"y"`
}

func (d *ApiDataZoneEntry) GetJsonKeyValue(key string) (interface{}, error) {
	switch key {
	case "name":
		return d.Name, nil
	case "id":
		return d.Id, nil
	case "map_id":
		return d.MapId, nil
	case "site_id":
		return d.SiteId, nil
	case "org_id":
		return d.OrgId, nil
	case "occupancy_limit":
		return d.OccupancyLimit, nil
	case "created_time":
		return d.CreatedTime, nil
	case "modified_time":
		return d.ModifiedTime, nil
	case "vertices":
		return d.Vertices, nil
	case "vertices_m":
		return d.VerticesM, nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (d *ApiDataZoneEntry) GetJsonKeyValueAsStr(key string) (string, error) {
	switch key {
	case "name":
		return d.Name, nil
	case "id":
		return d.Id, nil
	case "map_id":
		return d.MapId, nil
	case "site_id":
		return d.SiteId, nil
	case "org_id":
		return d.OrgId, nil
	case "occupancy_limit":
		return string(d.OccupancyLimit), nil
	case "created_time":
		return string(d.CreatedTime), nil
	case "modified_time":
		return string(d.ModifiedTime), nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (d *ApiDataZoneEntry) GetJsonKeyValueAsFloat64(key string) (float64, error) {
	switch key {
	case "occupancy_limit":
		return d.OccupancyLimit.Float64()
	case "created_time":
		return d.CreatedTime.Float64()
	case "modified_time":
		return d.ModifiedTime.Float64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

func (d *ApiDataZoneEntry) GetJsonKeyValueAsInt64(key string) (int64, error) {
	switch key {
	case "occupancy_limit":
		return d.OccupancyLimit.Int64()
	case "created_time":
		return d.CreatedTime.Int64()
	case "modified_time":
		return d.ModifiedTime.Int64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

func (d *ApiDataZoneVertice) GetJsonKeyValue(key string) (interface{}, error) {
	switch key {
	case "x":
		return d.X, nil
	case "y":
		return d.Y, nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (d *ApiDataZoneVertice) GetJsonKeyValueAsStr(key string) (string, error) {
	switch key {
	case "x":
		return string(d.X), nil
	case "y":
		return string(d.Y), nil
	default:
		return "", fmt.Errorf("Specified key not found")
	}

	return "", fmt.Errorf("Specified key not found")
}

func (d *ApiDataZoneVertice) GetJsonKeyValueAsFloat64(key string) (float64, error) {
	switch key {
	case "x":
		return d.X.Float64()
	case "y":
		return d.Y.Float64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

func (d *ApiDataZoneVertice) GetJsonKeyValueAsInt64(key string) (int64, error) {
	switch key {
	case "x":
		return d.X.Int64()
	case "y":
		return d.Y.Int64()
	default:
		return 0, fmt.Errorf("Specified key not found")
	}

	return 0, fmt.Errorf("Specified key not found")
}

/*
 * Asset search API data format
 * /api/v1/sites/:site_id/stats/assets/search?mac=fbc721cc3022
 */
type ApiDataAssetSearchResult struct {
	Start			json.Number		`json:"start"`
	End			json.Number		`json:"end"`
	Limit			json.Number		`json:"limit"`
	Total			json.Number		`json:"total"`
	
	Results			[]ApiDataAssetEntry	`json:"results"`
}

type ApiDataAssetEntry struct {
	Mac			string		`json:"mac"`
	MapId			string		`json:"map_id"`
	By			string		`json:"by"`
	Name			string		`json:"name"`

	Manufacture		string		`json:"manufacture"`
	Temperature		json.Number	`json:"temperature"`
	BattVoltage		json.Number	`json:"battery_voltage"`
	EddystoneUIDNamespace	string		`json:"eddystone_uid_namespace"`
	EddystoneUIDInstance	string		`json:"eddystone_uid_instance"`
	EddystoneURL		string		`json:"eddystone_url"`

	Beam			json.Number	`json:"beam"`
	Rssi			json.Number	`json:"rssi"`
	ApMac			string		`json:"ap_mac"`
	Lastseen		json.Number	`json:"last_seen"`
	Timestamp		json.Number	`json:"_timestamp"`
}


