package restfire

type rtdb struct {
	url     string
	tokenid string
}
type rtdberr struct {
	Err string `json:"error"`
}

func NewRealtimeDatabase(url string, tokenid string) RealtimeDatabase {
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	return &rtdb{url, tokenid}
}

func (x *rtdb) authParam() string {
	return "auth=" + x.tokenid
}

func (x *rtdb) Get(path string, out interface{}) error {
	return requestdata("GET", x.url+path+".json?"+x.authParam(), nil, out, &rtdberr{})
}
func (x *rtdb) Write(path string, data interface{}) error {
	return requestdata("PUT", x.url+path+".json?"+x.authParam(), data, nil, &rtdberr{})
}
func (x *rtdb) Update(path string, data interface{}) error {
	return requestdata("PATCH", x.url+path+".json?"+x.authParam(), data, nil, &rtdberr{})
}
func (x *rtdb) Push(path string, data interface{}) (string, error) {
	var out struct {
		Name string `json:"name"`
	}
	err := requestdata("POST", x.url+path+".json?"+x.authParam(), data, &out, &rtdberr{})
	if err != nil {
		return "", err
	}
	return out.Name, nil
}
func (x *rtdb) Delete(path string) error {
	return requestdata("DELETE", x.url+path+".json?"+x.authParam(), nil, nil, &rtdberr{})
}
func (x *rtdberr) Error() string {
	return x.Err
}
