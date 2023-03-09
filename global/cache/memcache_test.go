package cache

import (
	"bytes"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"monaToolBox/bootstrap"
	"reflect"
	"testing"
)

var globalmc *memcache.Client
var mc *cacheMc
var getStruct []byte

func init() {
	bootstrap.InitConfig()

	// 初始化memCached
	globalmc = bootstrap.InitializeMemcached()
	mc = NewMemcachedClient("test")
	mc.c = globalmc

	getStruct, _ = json.Marshal(struct {
		A string  `json:"a"`
		B int     `json:"b"`
		c float64 `json:"c"`
	}{"aaaaa", 2, 3.1415})
}

func TestMain(m *testing.M) {

	m.Run()
	//
	// mc.Delete("string")
	// mc.Delete("byteArr")
	// mc.Delete("timeout")
	// mc.Delete("struct")
}

func TestNewMemcachedClient(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want *cacheMc
	}{
		{
			name: "get a new memcached client",
			args: args{
				prefix: "test",
			},
			want: &cacheMc{prefix: "test_"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemcachedClient(tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemcachedClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cacheMc_Add(t *testing.T) {

	mc.Delete("string")
	mc.Delete("byteArr")
	mc.Delete("timeout")
	mc.Delete("struct")

	type args struct {
		key        string
		item       any
		expiration int32
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		getResult []byte
	}{
		{
			name: "test add",
			args: args{
				key:        "string",
				item:       "hello world!",
				expiration: 3600,
			},
			wantErr:   false,
			getResult: []byte("hello world!"),
		},
		{
			name: "test add",
			args: args{
				key:        "byteArr",
				item:       []byte("byte array"),
				expiration: 3600,
			},
			wantErr:   false,
			getResult: []byte("byte array"),
		},
		{
			name: "test add",
			args: args{
				key:        "timeout",
				item:       []byte("timeout 1s"),
				expiration: 1,
			},
			wantErr:   false,
			getResult: []byte("timeout 1s"),
		},
		{
			name: "test add",
			args: args{
				key: "struct",
				item: struct {
					A string  `json:"a"`
					B int     `json:"b"`
					c float64 `json:"c"`
				}{"aaaaa", 2, 3.1415},
				expiration: 3600,
			},
			wantErr:   false,
			getResult: getStruct,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mc.Set(tt.args.key, tt.args.item, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if get, err := mc.Get(tt.args.key); err != nil || bytes.Compare(get, tt.getResult) != 0 {
				t.Errorf("Set() error = %v, result = %v need = %v", err, get, tt.getResult)
			}
		})
	}
}

func Test_cacheMc_Get(t *testing.T) {

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"get test",
			args{"string"},
			[]byte("hello world!"),
			false,
		},
		{
			"get test",
			args{"byteArr"},
			[]byte("byte array"),
			false,
		},
		{
			"get test",
			args{"timeout"},
			nil,
			true,
		},
		{
			"get test",
			args{"struct"},
			getStruct,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mc.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v, err %v", got, tt.want, err)
			}
		})
	}
}

func Test_cacheMc_GetMulti(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		wantErr bool
	}{
		{
			name: "multi get",
			args: args{
				keys: []string{"string", "byteArr", "timeout", "struct"},
			},
			want: map[string][]byte{
				mc.buildKey("string"):  []byte("hello world!"),
				mc.buildKey("byteArr"): []byte("byte array"),
				mc.buildKey("struct"):  getStruct,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mc.GetMulti(tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMulti() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cacheMc_Delete(t *testing.T) {

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test delete",
			args: args{
				key: "string",
			},
			wantErr: false,
		},
		{
			name: "test delete err",
			args: args{
				key: "stringErr",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mc.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cacheMc_buildKey(t *testing.T) {
	type fields struct {
		c      *memcache.Client
		prefix string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Set test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &cacheMc{
				c:      tt.fields.c,
				prefix: tt.fields.prefix,
			}
			if got := mc.buildKey(tt.args.key); got != tt.want {
				t.Errorf("buildKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
