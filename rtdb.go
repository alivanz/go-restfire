package restfire

type rtdb struct {
	url       string
	refresher AuthRefresher
}
type rtdberr struct {
	Err string `json:"error"`
}

func NewRealtimeDatabase(url string, refresher AuthRefresher) RealtimeDatabase {
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	return &rtdb{url, refresher}
}

func (x *rtdb) authParam() string {
	return "auth=" + x.refresher.Token().IDToken()
}

func (x *rtdb) rtdb_request(method string, url string, data interface{}, out interface{}, errframe error) error {
	if x.refresher == nil {
		return requestdata(method, url, data, out, errframe)
	}
	for {
		err := requestdata(method, url+"?"+x.authParam(), data, out, &rtdberr{})
		if err != nil {
			if err.Error() != "Auth token is expired" {
				if err = x.refresher.AuthRefresh(); err != nil {
					return err
				}
			}
			return err
		}
		return nil
	}
}
func (x *rtdb) Get(path string, out interface{}) error {
	return x.rtdb_request("GET", x.url+path+".json", nil, out, &rtdberr{})
}
func (x *rtdb) Write(path string, data interface{}) error {
	return x.rtdb_request("PUT", x.url+path+".json", data, nil, &rtdberr{})
}
func (x *rtdb) Update(path string, data interface{}) error {
	return x.rtdb_request("PATCH", x.url+path+".json", data, nil, &rtdberr{})
}
func (x *rtdb) Push(path string, data interface{}) (string, error) {
	var out struct {
		Name string `json:"name"`
	}
	err := x.rtdb_request("POST", x.url+path+".json", data, &out, &rtdberr{})
	if err != nil {
		return "", err
	}
	return out.Name, nil
}
func (x *rtdb) Delete(path string) error {
	return x.rtdb_request("DELETE", x.url+path+".json", nil, nil, &rtdberr{})
}
func (x *rtdberr) Error() string {
	return x.Err
}
