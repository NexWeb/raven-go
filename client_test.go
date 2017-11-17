The Go PlaygroundRun  Format   Imports  Share About
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
41
42
43
44
45
46
47
48
49
50
51
52
53
54
55
56
57
58
59
60
61
62
63
64
65
66
67
68
69
70
71
72
73
74
75
76
77
78
79
80
81
82
83
84
85
86
87
88
89
90
91
92
93
94
95
96
97
98
99
100
101
102
103
104
105
106
107
108
109
110
111
112
113
114
115
116
117
118
119
120
121
122
123
124
125
126
127
128
129
130
131
132
133
134
135
136
137
138
139
140
141
142
143
144
145
146
147
148
149
150
151
152
153
154
155
156
157
158
159
160
161
162
163
164
165
166
167
168
169
170
171
172
173
174
175
176
177
178
179
180
181
182
183
184
185
186
187
188
189
190
191
192
193
194
195
196
197
198
199
200
201
202
203
204
205
206
207
208
209
210
211
212
213
214
215
216
217
218
219
220
221
222
223
224
225
226
227
228
229
230
231
232
233
234
235
236
237
238
239
240
241
242
243
244
245
246
247
248
249
250
251
252
253
254
255
256
257
258

package raven

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type testInterface struct{}

func (t *testInterface) Class() string   { return "sentry.interfaces.Test" }
func (t *testInterface) Culprit() string { return "codez" }

func TestShouldExcludeErr(t *testing.T) {
	regexpStrs := []string{"ERR_TIMEOUT", "should.exclude", "(?i)^big$"}

	client := &Client{
		Transport: newTransport(),
		Tags:      nil,
		context:   &context{},
		queue:     make(chan *outgoingPacket, MaxQueueBuffer),
	}

	if err := client.SetIgnoreErrors(regexpStrs); err != nil {
		t.Fatalf("invalid regexps %v: %v", regexpStrs, err)
	}

	testCases := []string{
		"there was a ERR_TIMEOUT in handlers.go",
		"do not log should.exclude at all",
		"BIG",
	}

	for _, tc := range testCases {
		if !client.shouldExcludeErr(tc) {
			t.Fatalf("failed to exclude err %q with regexps %v", tc, regexpStrs)
		}
	}
}

func TestPacketJSON(t *testing.T) {
	packet := &Packet{
		Project:     "1",
		EventID:     "2",
		Platform:    "linux",
		Culprit:     "caused_by",
		ServerName:  "host1",
		Release:     "721e41770371db95eee98ca2707686226b993eda",
		Environment: "production",
		Message:     "test",
		Timestamp:   Timestamp(time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC)),
		Level:       ERROR,
		Logger:      "com.getsentry.raven-go.logger-test-packet-json",
		Tags:        []Tag{Tag{"foo", "bar"}},
		Modules:     map[string]string{"foo": "bar"},
		Fingerprint: []string{"{{ default }}", "a-custom-fingerprint"},
		Interfaces:  []Interface{&Message{Message: "foo"}},
	}

	packet.AddTags(map[string]string{"foo": "foo"})
	packet.AddTags(map[string]string{"baz": "buzz"})

	expected := `{"message":"test","event_id":"2","project":"1","timestamp":"2000-01-01T00:00:00.00","level":"error","logger":"com.getsentry.raven-go.logger-test-packet-json","platform":"linux","culprit":"caused_by","server_name":"host1","release":"721e41770371db95eee98ca2707686226b993eda","environment":"production","tags":[["foo","bar"],["foo","foo"],["baz","buzz"]],"modules":{"foo":"bar"},"fingerprint":["{{ default }}","a-custom-fingerprint"],"logentry":{"message":"foo"}}`
	j, err := packet.JSON()
	if err != nil {
		t.Fatalf("JSON marshalling should not fail: %v", err)
	}
	actual := string(j)

	if actual != expected {
		t.Errorf("incorrect json; got %s, want %s", actual, expected)
	}
}

func TestPacketJSONNilInterface(t *testing.T) {
	packet := &Packet{
		Project:     "1",
		EventID:     "2",
		Platform:    "linux",
		Culprit:     "caused_by",
		ServerName:  "host1",
		Release:     "721e41770371db95eee98ca2707686226b993eda",
		Environment: "production",
		Message:     "test",
		Timestamp:   Timestamp(time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC)),
		Level:       ERROR,
		Logger:      "com.getsentry.raven-go.logger-test-packet-json",
		Tags:        []Tag{Tag{"foo", "bar"}},
		Modules:     map[string]string{"foo": "bar"},
		Fingerprint: []string{"{{ default }}", "a-custom-fingerprint"},
		Interfaces:  []Interface{&Message{Message: "foo"}, nil},
	}

	expected := `{"message":"test","event_id":"2","project":"1","timestamp":"2000-01-01T00:00:00.00","level":"error","logger":"com.getsentry.raven-go.logger-test-packet-json","platform":"linux","culprit":"caused_by","server_name":"host1","release":"721e41770371db95eee98ca2707686226b993eda","environment":"production","tags":[["foo","bar"]],"modules":{"foo":"bar"},"fingerprint":["{{ default }}","a-custom-fingerprint"],"logentry":{"message":"foo"}}`
	j, err := packet.JSON()
	if err != nil {
		t.Fatalf("JSON marshalling should not fail: %v", err)
	}
	actual := string(j)

	if actual != expected {
		t.Errorf("incorrect json; got %s, want %s", actual, expected)
	}
}

func TestPacketInit(t *testing.T) {
	packet := &Packet{Message: "a", Interfaces: []Interface{&testInterface{}}}
	packet.Init("foo")

	if packet.Project != "foo" {
		t.Error("incorrect Project:", packet.Project)
	}
	if packet.Culprit != "codez" {
		t.Error("incorrect Culprit:", packet.Culprit)
	}
	if packet.ServerName == "" {
		t.Errorf("ServerName should not be empty")
	}
	if packet.Level != ERROR {
		t.Errorf("incorrect Level: got %v, want %v", packet.Level, ERROR)
	}
	if packet.Logger != "root" {
		t.Errorf("incorrect Logger: got %s, want %s", packet.Logger, "root")
	}
	if time.Time(packet.Timestamp).IsZero() {
		t.Error("Timestamp is zero")
	}
	if len(packet.EventID) != 32 {
		t.Error("incorrect EventID:", packet.EventID)
	}
}

func TestSetDSN(t *testing.T) {
	client := &Client{}
	client.SetDSN("https://u:p@example.com/sentry/1")

	if client.url != "https://example.com/sentry/api/1/store/" {
		t.Error("incorrect url:", client.url)
	}
	if client.projectID != "1" {
		t.Error("incorrect projectID:", client.projectID)
	}
	if client.authHeader != "Sentry sentry_version=4, sentry_key=u, sentry_secret=p" {
		t.Error("incorrect authHeader:", client.authHeader)
	}
}

func TestNewClient(t *testing.T) {
	client := newClient(nil)
	if client.sampleRate != 1.0 {
		t.Error("invalid default sample rate")
	}
}

func TestSetSampleRate(t *testing.T) {
	client := &Client{}
	err := client.SetSampleRate(0.2)

	if err != nil {
		t.Error("invalid sample rate")
	}

	if client.sampleRate != 0.2 {
		t.Error("incorrect sample rate: ", client.sampleRate)
	}
}

func TestSetSampleRateInvalid(t *testing.T) {
	client := &Client{}
	err := client.SetSampleRate(-1.0)

	if err != ErrInvalidSampleRate {
		t.Error("invalid sample rate should return ErrInvalidSampleRate")
	}
}

func TestUnmarshalTag(t *testing.T) {
	actual := new(Tag)
	if err := json.Unmarshal([]byte(`["foo","bar"]`), actual); err != nil {
		t.Fatal("unable to decode JSON:", err)
	}

	expected := &Tag{Key: "foo", Value: "bar"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("incorrect Tag: wanted '%+v' and got '%+v'", expected, actual)
	}
}

func TestUnmarshalTags(t *testing.T) {
	tests := []struct {
		Input    string
		Expected Tags
	}{
		{
			`{"foo":"bar"}`,
			Tags{Tag{Key: "foo", Value: "bar"}},
		},
		{
			`[["foo","bar"],["bar","baz"]]`,
			Tags{Tag{Key: "foo", Value: "bar"}, Tag{Key: "bar", Value: "baz"}},
		},
	}

	for _, test := range tests {
		var actual Tags
		if err := json.Unmarshal([]byte(test.Input), &actual); err != nil {
			t.Fatal("unable to decode JSON:", err)
		}

		if !reflect.DeepEqual(actual, test.Expected) {
			t.Errorf("incorrect Tags: wanted '%+v' and got '%+v'", test.Expected, actual)
		}
	}
}

func TestMarshalTimestamp(t *testing.T) {
	timestamp := Timestamp(time.Date(2000, 01, 02, 03, 04, 05, 0, time.UTC))
	expected := `"2000-01-02T03:04:05.00"`

	actual, err := json.Marshal(timestamp)
	if err != nil {
		t.Error(err)
	}

	if string(actual) != expected {
		t.Errorf("incorrect string; got %s, want %s", actual, expected)
	}
}

func TestUnmarshalTimestamp(t *testing.T) {
	timestamp := `"2000-01-02T03:04:05.00"`
	expected := Timestamp(time.Date(2000, 01, 02, 03, 04, 05, 0, time.UTC))

	var actual Timestamp
	err := json.Unmarshal([]byte(timestamp), &actual)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Errorf("incorrect string; got %v, want %v", actual, expected)
	}
}

func TestNilClient(t *testing.T) {
	var client *Client = nil
	eventID, ch := client.Capture(nil, nil)
	if eventID != "" {
		t.Error("expected empty eventID:", eventID)
	}
	// wait on ch: no send should succeed immediately
	err := <-ch
	if err != nil {
		t.Error("expected nil err:", err)
	}
}
