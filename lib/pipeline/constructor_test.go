// Copyright (c) 2018 Ashley Jeffs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package pipeline

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Jeffail/benthos/lib/metrics"
	"github.com/Jeffail/benthos/lib/processor"
	"github.com/Jeffail/benthos/lib/types"
	"github.com/Jeffail/benthos/lib/util/service/log"
)

func TestSanitise(t *testing.T) {
	var actObj interface{}
	var act []byte
	var err error

	exp := `{` +
		`"processors":[],` +
		`"threads":10` +
		`}`

	conf := NewConfig()
	conf.Threads = 10
	conf.Processors = nil

	if actObj, err = SanitiseConfig(conf); err != nil {
		t.Fatal(err)
	}
	if act, err = json.Marshal(actObj); err != nil {
		t.Fatal(err)
	}
	if string(act) != exp {
		t.Errorf("Wrong sanitised output: %s != %v", act, exp)
	}

	exp = `{` +
		`"processors":[` +
		`{` +
		`"type":"combine",` +
		`"combine":{` +
		`"parts":2` +
		`}` +
		`},` +
		`{` +
		`"type":"archive",` +
		`"archive":{` +
		`"format":"binary",` +
		`"path":"nope"` +
		`}` +
		`}` +
		`],` +
		`"threads":10` +
		`}`

	proc := processor.NewConfig()
	proc.Type = "combine"
	conf.Processors = append(conf.Processors, proc)

	proc = processor.NewConfig()
	proc.Type = "archive"
	proc.Archive.Path = "nope"
	conf.Processors = append(conf.Processors, proc)

	if actObj, err = SanitiseConfig(conf); err != nil {
		t.Fatal(err)
	}
	if act, err = json.Marshal(actObj); err != nil {
		t.Fatal(err)
	}
	if string(act) != exp {
		t.Errorf("Wrong sanitised output: %s != %v", act, exp)
	}
}

func TestProcCtor(t *testing.T) {
	firstProc := processor.NewConfig()
	firstProc.Type = "bounds_check"
	firstProc.BoundsCheck.MinPartSize = 5

	secondProc := processor.NewConfig()
	secondProc.Type = "insert_part"
	secondProc.InsertPart.Content = "1"
	secondProc.InsertPart.Index = 0

	conf := NewConfig()
	conf.Processors = append(conf.Processors, firstProc)

	pipe, err := New(
		conf, nil,
		log.NewLogger(os.Stdout, log.LoggerConfig{LogLevel: "NONE"}),
		metrics.DudType{},
		func() (processor.Type, error) {
			return processor.New(
				secondProc, nil,
				log.NewLogger(os.Stdout, log.LoggerConfig{LogLevel: "NONE"}),
				metrics.DudType{},
			)
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	tChan := make(chan types.Transaction)
	resChan := make(chan types.Response)

	if err = pipe.StartReceiving(tChan); err != nil {
		t.Fatal(err)
	}

	select {
	case <-time.After(time.Second):
		t.Fatal("timed out")
	case tChan <- types.NewTransaction(
		types.NewMessage([][]byte{[]byte("foo bar baz")}), resChan,
	):
	}

	var tran types.Transaction
	select {
	case <-time.After(time.Second):
		t.Fatal("timed out")
	case tran = <-pipe.TransactionChan():
	}

	exp := [][]byte{
		[]byte("1"),
		[]byte("foo bar baz"),
	}
	if act := tran.Payload.GetAll(); !reflect.DeepEqual(exp, act) {
		t.Errorf("Wrong contents: %s != %s", act, exp)
	}

	select {
	case <-time.After(time.Second):
		t.Fatal("timed out")
	case tran.ResponseChan <- types.NewSimpleResponse(nil):
	}

	select {
	case <-time.After(time.Second):
		t.Fatal("timed out")
	case res := <-resChan:
		if res.Error() != nil {
			t.Error(res.Error())
		}
	}

	pipe.CloseAsync()
	if err = pipe.WaitForClose(time.Second); err != nil {
		t.Error(err)
	}
}
