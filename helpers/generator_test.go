package helpers

import (
	"fmt"
	"testing"

	kubereplayv1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	"reflect"
)

func TestArgsFromSpec(t *testing.T) {
	tests := []struct {
		Workers           int32
		Timeout           string
		FileSilo          kubereplayv1alpha1.FileSilo
		TcpSilo           kubereplayv1alpha1.TcpSilo
		StdoutSilo        kubereplayv1alpha1.StdoutSilo
		HttpSilo          kubereplayv1alpha1.HttpSilo
		ElasticsearchSilo kubereplayv1alpha1.ElasticsearchSilo
		KafkaSilo         kubereplayv1alpha1.KafkaSilo
		ExpResult         []string
	}{
		{
			1,
			"10s",
			kubereplayv1alpha1.FileSilo{Enabled: false},
			kubereplayv1alpha1.TcpSilo{Enabled: false},
			kubereplayv1alpha1.StdoutSilo{Enabled: false},
			kubereplayv1alpha1.HttpSilo{Enabled: false},
			kubereplayv1alpha1.ElasticsearchSilo{Enabled: false},
			kubereplayv1alpha1.KafkaSilo{Enabled: false},
			[]string{"--input-tcp", ":28020", "--output-http-workers", "1", "-output-http-timeout", "10s"},
		},
		{
			5,
			"10s",
			kubereplayv1alpha1.FileSilo{
				Enabled:  true,
				Filename: "/tmp/file",
			},
			kubereplayv1alpha1.TcpSilo{Enabled: false},
			kubereplayv1alpha1.StdoutSilo{Enabled: false},
			kubereplayv1alpha1.HttpSilo{Enabled: false},
			kubereplayv1alpha1.ElasticsearchSilo{Enabled: false},
			kubereplayv1alpha1.KafkaSilo{Enabled: false},
			[]string{
				"--input-tcp",
				":28020",
				"--output-file",
				"/tmp/file",
				"--output-http-workers",
				"5",
				"-output-http-timeout",
				"10s",
			},
		},
		{
			1,
			"10s",
			kubereplayv1alpha1.FileSilo{Enabled: false},
			kubereplayv1alpha1.TcpSilo{
				Enabled: true,
				Uri:     "tcp://test",
			},
			kubereplayv1alpha1.StdoutSilo{Enabled: false},
			kubereplayv1alpha1.HttpSilo{Enabled: false},
			kubereplayv1alpha1.ElasticsearchSilo{Enabled: false},
			kubereplayv1alpha1.KafkaSilo{Enabled: false},
			[]string{
				"--input-tcp",
				":28020",
				"--output-tcp",
				"tcp://test",
				"--output-http-workers",
				"1",
				"-output-http-timeout",
				"10s",
			},
		},
		{
			1,
			"10s",
			kubereplayv1alpha1.FileSilo{Enabled: false},
			kubereplayv1alpha1.TcpSilo{Enabled: false},
			kubereplayv1alpha1.StdoutSilo{
				Enabled: true,
			},
			kubereplayv1alpha1.HttpSilo{Enabled: false},
			kubereplayv1alpha1.ElasticsearchSilo{Enabled: false},
			kubereplayv1alpha1.KafkaSilo{Enabled: false},
			[]string{
				"--input-tcp",
				":28020",
				"--output-stdout",
				"--output-http-workers",
				"1",
				"-output-http-timeout",
				"10s",
			},
		},
		{
			1,
			"10s",
			kubereplayv1alpha1.FileSilo{Enabled: false},
			kubereplayv1alpha1.TcpSilo{Enabled: false},
			kubereplayv1alpha1.StdoutSilo{Enabled: false},
			kubereplayv1alpha1.HttpSilo{
				Enabled: true,
				Uri:     "https://localhost:2080",
			},
			kubereplayv1alpha1.ElasticsearchSilo{Enabled: false},
			kubereplayv1alpha1.KafkaSilo{Enabled: false},
			[]string{
				"--input-tcp",
				":28020",
				"--output-http",
				"https://localhost:2080",
				"--output-http-workers",
				"1",
				"-output-http-timeout",
				"10s"},
		},
		{
			1,
			"10s",
			kubereplayv1alpha1.FileSilo{Enabled: false},
			kubereplayv1alpha1.TcpSilo{Enabled: false},
			kubereplayv1alpha1.StdoutSilo{Enabled: false},
			kubereplayv1alpha1.HttpSilo{Enabled: false},
			kubereplayv1alpha1.ElasticsearchSilo{
				Enabled: true,
				Uri:     "http://localhost:9200",
			},
			kubereplayv1alpha1.KafkaSilo{Enabled: false},
			[]string{
				"--input-tcp",
				":28020",
				"--output-http-elasticsearch",
				"http://localhost:9200",
				"--output-http-workers",
				"1",
				"-output-http-timeout",
				"10s"},
		},
		{
			1,
			"10s",
			kubereplayv1alpha1.FileSilo{Enabled: false},
			kubereplayv1alpha1.TcpSilo{Enabled: false},
			kubereplayv1alpha1.StdoutSilo{Enabled: false},
			kubereplayv1alpha1.HttpSilo{Enabled: false},
			kubereplayv1alpha1.ElasticsearchSilo{Enabled: false},
			kubereplayv1alpha1.KafkaSilo{
				Enabled: true,
				Uri:     "192.168.1.1",
			},
			[]string{
				"--input-tcp",
				":28020",
				"--output-kafka-host",
				"192.168.1.1",
				"--output-http-workers",
				"1",
				"-output-http-timeout",
				"10s",
			},
		},
		{
			1,
			"10s",
			kubereplayv1alpha1.FileSilo{
				Enabled:  true,
				Filename: "/tmp/file",
			},
			kubereplayv1alpha1.TcpSilo{
				Enabled: true,
				Uri:     "tcp://test",
			},
			kubereplayv1alpha1.StdoutSilo{
				Enabled: true,
			},
			kubereplayv1alpha1.HttpSilo{
				Enabled: true,
				Uri:     "https://localhost:2080",
			},
			kubereplayv1alpha1.ElasticsearchSilo{
				Enabled: true,
				Uri:     "http://localhost:9200",
			},
			kubereplayv1alpha1.KafkaSilo{
				Enabled: true,
				Uri:     "192.168.1.1",
			},
			[]string{
				"--input-tcp",
				":28020",
				"--output-file",
				"/tmp/file",
				"--output-tcp",
				"tcp://test",
				"--output-stdout",
				"--output-http",
				"https://localhost:2080",
				"--output-http-elasticsearch",
				"http://localhost:9200",
				"--output-kafka-host",
				"192.168.1.1",
				"--output-http-workers",
				"1",
				"-output-http-timeout",
				"10s",
			},
		},
	}

	for i, tt := range tests {
		testCase := fmt.Sprintf("%d/%d", i+1, len(tests))
		spec := kubereplayv1alpha1.RefinerySpec{
			Workers: tt.Workers,
			Timeout: tt.Timeout,
			Storage: &kubereplayv1alpha1.RefineryStorage{
				File:          &tt.FileSilo,
				Tcp:           &tt.TcpSilo,
				Stdout:        &tt.StdoutSilo,
				Http:          &tt.HttpSilo,
				Elasticsearch: &tt.ElasticsearchSilo,
				Kafka:         &tt.KafkaSilo,
			},
		}
		res := argsFromSpec(&spec)
		if !reflect.DeepEqual(res, &tt.ExpResult) {
			t.Fatalf("%s failed, \n expected result %v,\ns got %v", testCase, &tt.ExpResult, res)
		}
	}
}

func TestHttpSiloToArgs(t *testing.T) {
	tests := []struct {
		Uri            string
		Debug          bool
		ResponseBuffer int
		ExpResult      []string
		ExpErr         bool
	}{
		{
			"http://192.168.1.1",
			false,
			0,
			[]string{
				"--output-http",
				"http://192.168.1.1",
			},
			false,
		},
		{
			"https://192.168.1.1",
			true,
			10,
			[]string{
				"--output-http",
				"https://192.168.1.1",
				"--output-http-debug",
				"--output-http-response-buffer",
				"10",
			},
			false,
		},
		{
			"",
			true,
			0,
			[]string{},
			true,
		},
	}
	for i, tt := range tests {
		testCase := fmt.Sprintf("%d/%d", i+1, len(tests))
		spec := kubereplayv1alpha1.HttpSilo{
			Uri:            tt.Uri,
			Debug:          tt.Debug,
			ResponseBuffer: tt.ResponseBuffer,
		}
		res, err := httpSiloToArgs(&spec)
		if err != nil {
			if tt.ExpErr == false {
				t.Fatalf("%s failed, got unexpected err %v", testCase, err)
			}
		} else {
			if tt.ExpErr == true {
				t.Fatalf("%s failed, error was expected but did not happen", testCase)
			} else {
				if !reflect.DeepEqual(res, &tt.ExpResult) {
					t.Fatalf("%s failed, \n expected result %v,\ns got %v", testCase, &tt.ExpResult, res)
				}
			}
		}

	}
}

func TestTcpSiloToArgs(t *testing.T) {
	tests := []struct {
		Uri       string
		ExpResult []string
		ExpErr    bool
	}{
		{
			"tcp://localhost:28020",
			[]string{
				"--output-tcp",
				"tcp://localhost:28020",
			},
			false,
		},
		{
			"",
			[]string{
				"should-be-error-anyways",
			},
			true,
		},
		{
			":28020",
			[]string{
				"--output-tcp",
				":28020",
			},
			false,
		},
	}

	for i, tt := range tests {
		testCase := fmt.Sprintf("%d/%d", i+1, len(tests))
		spec := kubereplayv1alpha1.TcpSilo{
			Uri: tt.Uri,
		}
		res, err := tcpSiloToArgs(&spec)
		if err != nil {
			if tt.ExpErr == false {
				t.Fatalf("%s failed, got unexpected err %v", testCase, err)
			}
		} else {
			if tt.ExpErr == true {
				t.Fatalf("%s failed, error was expected but did not happen", testCase)
			} else {
				if !reflect.DeepEqual(res, &tt.ExpResult) {
					t.Fatalf("%s failed, \n expected result %v,\ns got %v", testCase, &tt.ExpResult, res)
				}
			}
		}

	}

}

func TestFileSiloToArgs(t *testing.T) {
	tests := []struct {
		Filename      string
		Append        bool
		FlushInterval string
		QueueSize     int32
		FileLimit     string
		ExpResult     []string
		ExpErr        bool
	}{
		{
			"/bin/test",
			true,
			"10s",
			1,
			"10",
			[]string{
				"--output-file",
				"/bin/test",
				"--output-file-append",
				"--output-file-flush-interval",
				"10s",
				"--output-file-queue-limit",
				"1",
				"--output-file-size-limit",
				"10",
			},
			false,
		},
		{
			"/bin/test",
			false,
			"10s",
			0,
			"10",
			[]string{
				"--output-file",
				"/bin/test",
				"--output-file-flush-interval",
				"10s",
				"--output-file-size-limit",
				"10",
			},
			false,
		},
		{
			"",
			true,
			"10s",
			1,
			"10",

			[]string{
				"--output-file",
				"/bin/test",
				"--output-file-append",
				"--output-file-flush-interval",
				"10s",
				"--output-file-queue-limit",
				"1",
				"--output-file-size-limit",
				"10",
			},
			true,
		},
	}
	for i, tt := range tests {
		testCase := fmt.Sprintf("%d/%d", i+1, len(tests))
		spec := kubereplayv1alpha1.FileSilo{
			Filename:      tt.Filename,
			Append:        tt.Append,
			FlushInterval: tt.FlushInterval,
			QueueSize:     tt.QueueSize,
			FileLimit:     tt.FileLimit,
		}
		res, err := fileSiloToArgs(&spec)
		if err != nil {
			if tt.ExpErr == false {
				t.Fatalf("%s failed, got unexpected err %v", testCase, err)
			}
		} else {
			if tt.ExpErr == true {
				t.Fatalf("%s failed, error was expected but did not happen", testCase)
			} else {
				if !reflect.DeepEqual(res, &tt.ExpResult) {
					t.Fatalf("%s failed, \n expected result %v,\ns got %v", testCase, &tt.ExpResult, res)
				}
			}
		}

	}
}
