package SdkApi

func GetCacheByPerfix(config map[string]string, prefix string) (*Response, error) {
	return doGet(config["sdkRpcServer"] + "/sdk/consensus/get_cache_by_prefix?bcag_name=" + config["bcagName"] + "&prefix=" + prefix)
}
