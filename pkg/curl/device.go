package curl

import "strings"

func Mobile(useragent string) bool {
	// the list below is taken from
	// https://github.com/bcit-ci/CodeIgniter/blob/develop/system/libraries/User_agent.php

	mobiles := []string{"Mobile Explorer", "Palm", "Motorola", "Nokia", "Palm", "Apple iPhone", "iPad", "Apple iPod Touch", "Sony Ericsson", "Sony Ericsson", "BlackBerry", "O2 Cocoon", "Treo", "LG", "Amoi", "XDA", "MDA", "Vario", "HTC", "Samsung",
		"Sharp", "Siemens", "Alcatel", "BenQ", "HP iPaq", "Motorola", "PlayStation Portable", "PlayStation 3", "PlayStation Vita", "Danger Hiptop", "NEC", "Panasonic", "Philips", "Sagem", "Sanyo", "SPV", "ZTE", "Sendo", "Nintendo DSi", "Nintendo DS", "Nintendo 3DS", "Nintendo Wii", "Open Web", "OpenWeb", "Android", "Symbian", "SymbianOS", "Palm", "Symbian S60", "Windows CE", "Obigo", "Netfront Browser", "Openwave Browser", "Mobile Explorer", "Opera Mini", "Opera Mobile", "Firefox Mobile", "Digital Paths", "AvantGo", "Xiino", "Novarra Transcoder", "Vodafone", "NTT DoCoMo", "O2", "mobile", "wireless", "j2me", "midp", "cldc", "up.link", "up.browser", "smartphone", "cellphone", "Generic Mobile"}

	for _, device := range mobiles {
		if strings.Index(useragent, device) > -1 {
			return true
		}
	}
	return false
}

// GetDeviceType returns device type
// return 1 of 2 values:
//   - mobile
//   - web
func GetDeviceType(useragent string) string {
	if Mobile(useragent) {
		return "mobile"
	}
	return "web"
}
