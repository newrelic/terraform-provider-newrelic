package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// NewInsertClient makes a new client for the user to send data with
func NewInsertClient(insertKey string, accountID string) *InsertClient {
	client := &InsertClient{}
	client.URL = createInsertURL(accountID)
	client.InsertKey = insertKey
	client.Logger = log.New()
	client.Compression = None

	// Defaults
	client.RequestTimeout = DefaultInsertRequestTimeout
	client.RetryCount = DefaultRetries
	client.RetryWait = DefaultRetryWaitTime

	// Defaults for buffered client.
	// These are here so they can be overwritten before calling start().
	client.WorkerCount = DefaultWorkerCount
	client.BatchTime = DefaultBatchTimeout
	client.BatchSize = DefaultBatchEventCount

	return client
}

func createInsertURL(accountID string) *url.URL {
	insightsURL, _ := url.Parse(insightsInsertURL)
	insightsURL.Path = fmt.Sprintf("%s/%s/events", insightsURL.Path, accountID)
	return insightsURL
}

// Start runs the insert client in batch mode.
func (c *InsertClient) Start() error {
	if c.eventQueue != nil {
		return errors.New("the Insights client is already in daemon mode")
	}

	c.eventQueue = make(chan []byte, c.BatchSize)
	c.eventTimer = time.NewTimer(c.BatchTime)
	c.flushQueue = make(chan bool, c.WorkerCount)

	// TODO: errors returned from the call to watchdog()
	// and batchWorker() are simply dropped on the floor.
	go func() {
		err := c.watchdog()
		if err != nil {
			log.Errorf("watchdog returned error: %v", err)
		}
	}()

	go func() {
		err := c.batchWorker()
		if err != nil {
			log.Errorf("batch worker returned error: %v", err)
		}
	}()

	c.Logger.Infof("the Insights client has launched in daemon mode with endpoint %s", c.URL)

	return nil
}

// StartListener creates a goroutine that consumes from a channel and
// Enqueues events as to not block the writing of events to the channel
//
func (c *InsertClient) StartListener(inputChannel chan interface{}) (err error) {
	// Allow this to be called instead of Start()
	if c.eventQueue == nil {
		if err = c.Start(); err != nil {
			return err
		}
	}
	if inputChannel == nil {
		return errors.New("channel to listen is nil")
	}

	go func() {
		err := c.queueWorker(inputChannel)
		if err != nil {
			log.Errorf("queue worker returned error: %v", err)
		}
	}()

	c.Logger.Info("the Insights client started channel listener")

	return nil
}

// Validate makes sure the InsertClient is configured correctly for use
func (c *InsertClient) Validate() error {
	if correct, _ := regexp.MatchString("collector.newrelic.com/v1/accounts/[0-9]+/events", c.URL.String()); !correct {
		return fmt.Errorf("invalid insert endpoint %s", c.URL)
	}

	if len(c.InsertKey) < 1 {
		return fmt.Errorf("not a valid license key: %s", c.InsertKey)
	}
	return nil
}

// EnqueueEvent handles the queueing. Only works in batch mode.
func (c *InsertClient) EnqueueEvent(data interface{}) (err error) {
	if c.eventQueue == nil {
		return errors.New("queueing not enabled for this client")
	}

	var jsonData []byte
	atomic.AddInt64(&c.Statistics.EventCount, 1)

	if jsonData, err = json.Marshal(data); err != nil {
		return err
	}

	c.eventQueue <- jsonData

	return err
}

// PostEvent allows sending a single event directly.
func (c *InsertClient) PostEvent(data interface{}) error {
	var jsonData []byte

	switch data := data.(type) {
	case []byte:
		jsonData = data
	case string:
		jsonData = []byte(data)
	default:
		var jsonErr error
		jsonData, jsonErr = json.Marshal(data)
		if jsonErr != nil {
			return fmt.Errorf("error marshaling event data: %s", jsonErr.Error())
		}
	}

	// Needs to handle array of events. maybe pull into separate validation func
	if !strings.Contains(string(jsonData), "eventType") {
		return fmt.Errorf("event data must contain eventType field. %s", jsonData)
	}

	c.Logger.Debugf("Posting to insights: %s", jsonData)

	if requestErr := c.jsonPostRequest(jsonData); requestErr != nil {
		return requestErr
	}

	return nil
}

// Flush gives the user a way to manually flush the queue in the foreground.
// This is also used by watchdog when the timer expires.
func (c *InsertClient) Flush() error {
	if c.flushQueue == nil {
		return errors.New("queueing not enabled for this client")
	}
	c.Logger.Debug("Flushing insights client")
	atomic.AddInt64(&c.Statistics.FlushCount, 1)

	c.flushQueue <- true

	return nil
}

//
// queueWorker watches a channel and Enqueues items as they appear so
// we don't block on EnqueueEvent
//
func (c *InsertClient) queueWorker(inputChannel chan interface{}) (err error) {
	for { //nolint:gosimple
		select {
		case msg := <-inputChannel:
			err = c.EnqueueEvent(msg)
			if err != nil {
				return err
			}
		}
	}
}

//
// watchdog has a Timer that will send the results once the
// it has expired.
//
func (c *InsertClient) watchdog() (err error) {
	if c.eventTimer == nil {
		return errors.New("invalid timer for watchdog()")
	}

	for { //nolint:gosimple
		select {
		case <-c.eventTimer.C:
			// Timer expired, and we have data, send it
			atomic.AddInt64(&c.Statistics.TimerExpiredCount, 1)
			c.Logger.Debug("Timeout expired, flushing queued events")
			if err = c.Flush(); err != nil {
				return
			}
			c.eventTimer.Reset(c.BatchTime)
		}
	}
}

//
// batchWorker reads []byte from the queue until a threshold is passed,
// then copies the []byte it has read and sends that batch along to Insights
// in its own goroutine.
//
func (c *InsertClient) batchWorker() (err error) {
	eventBuf := make([][]byte, c.BatchSize)
	count := 0
	for {
		select {
		case item := <-c.eventQueue:
			eventBuf[count] = item
			count++
			if count >= c.BatchSize {
				c.grabAndConsumeEvents(count, eventBuf)
				count = 0
			}
		case <-c.flushQueue:
			if count > 0 {
				c.grabAndConsumeEvents(count, eventBuf)
				count = 0
			}
		}
	}
}

// grabAndConsumeEvents makes a copy of the event handles,
// and asynchronously writes those events in its own goroutine.
// The write is attempted up to c.RetryCount times.
//
// TODO: Any errors encountered doing the write are dropped on the floor.
// Even the last error (in the event of trying c.RetryCount times)
// is dropped.
//
func (c *InsertClient) grabAndConsumeEvents(count int, eventBuf [][]byte) {
	if count < c.BatchSize-20 {
		atomic.AddInt64(&c.Statistics.PartialFlushCount, 1) // Allow for some fuzz, although there should be none
	} else {
		atomic.AddInt64(&c.Statistics.FullFlushCount, 1)
	}

	saved := make([][]byte, count)
	for i := 0; i < count; i++ {
		saved[i] = eventBuf[i]
		eventBuf[i] = nil
	}

	go func(count int, saved [][]byte) {
		// only send the slice that we pulled into the buffer
		for tries := 0; tries < c.RetryCount; tries++ {
			if sendErr := c.sendEvents(saved[0:count]); sendErr != nil {
				if tries+1 >= c.RetryCount {
					//failed last retry
					c.Logger.Errorf("Failed to send insights events [%d/%d] times. Retry limit reached -- Abandoning data. Error: %v",
						tries+1, c.RetryCount, sendErr)
				} else {
					c.Logger.Errorf("Failed to send insights events [%d/%d]. Will retry. Error: %v", tries+1, c.RetryCount, sendErr)
					atomic.AddInt64(&c.Statistics.InsightsRetryCount, 1)
					time.Sleep(c.RetryWait)
				}
			} else {
				break
			}
		}
		atomic.AddInt64(&c.Statistics.ProcessedEventCount, int64(count))
	}(count, saved)
}

// sendEvents accepts a slice of marshalled JSON and sends it to Insights
//
func (c *InsertClient) sendEvents(events [][]byte) error {
	var buf bytes.Buffer

	// Since we already marshalled all of the data into JSON, let's make a
	// hand-crafted, artisanal JSON array
	buf.WriteString("[")
	eventCount := len(events) - 1
	for e := range events {
		buf.Write(events[e])
		if e < eventCount {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]")
	atomic.AddInt64(&c.Statistics.ByteCount, int64(buf.Len()))

	return c.jsonPostRequest(buf.Bytes())
}

// SetCompression allows modification of the compression type used in communication
//
func (c *InsertClient) SetCompression(compression Compression) {
	c.Compression = Gzip
	// use gzip only for now
	// c.Compression = compression
	log.Debugf("Compression set: %d", c.Compression)
}

func (c *InsertClient) jsonPostRequest(body []byte) (err error) {
	const prependText = "Insights Post: "

	req, reqErr := c.generateJSONPostRequest(body)
	if reqErr != nil {
		return fmt.Errorf("%s: %v", prependText, reqErr)
	}

	ctx, cancel := context.WithTimeout(req.Context(), c.RequestTimeout)
	defer cancel()
	resp, respErr := http.DefaultClient.Do(req.WithContext(ctx))
	if respErr != nil {
		return fmt.Errorf("%s: %v", prependText, respErr)
	}
	defer func() {
		respErr = resp.Body.Close()
		if respErr != nil && err == nil {
			err = respErr // Don't mask previous errors
		}
	}()

	if parseErr := c.parseResponse(resp); parseErr != nil {
		return fmt.Errorf("%s: %v", prependText, parseErr)
	}

	return nil
}

func (c *InsertClient) generateJSONPostRequest(body []byte) (request *http.Request, err error) {
	var readBuffer io.Reader
	var encoding string

	switch c.Compression {
	case None:
		c.Logger.Debug("Compression: None")
		readBuffer = bytes.NewBuffer(body)
	case Deflate:
		c.Logger.Debug("Compression: Deflate")
		readBuffer = nil
	case Gzip:
		c.Logger.Debug("Compression: Gzip")
		readBuffer, err = gZipBuffer(body)
		encoding = "gzip"
	case Zlib:
		c.Logger.Debug("Compression: Zlib")
		readBuffer = nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	request, err = http.NewRequest("POST", c.URL.String(), readBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to construct request for: %s", body)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Insert-Key", c.InsertKey)
	if encoding != "" {
		request.Header.Add("Content-Encoding", encoding)
	}

	return request, nil
}

// parseResponse checks the Insert response for errors and reports the message
// if an error happened
func (c *InsertClient) parseResponse(response *http.Response) error {
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return fmt.Errorf("failed to read response body: %s", readErr.Error())
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("bad response from Insights: %d \n\t%s", response.StatusCode, string(body))
	}

	c.Logger.Debugf("Response %d body: %s", response.StatusCode, body)

	respJSON := insertResponse{}
	if err := json.Unmarshal(body, &respJSON); err != nil {
		return fmt.Errorf("failed to unmarshal insights response: %v", err)
	}

	// Success
	if response.StatusCode == 200 && respJSON.Success {
		return nil
	}

	// Non 200 response (or 200 not success, if such a thing)
	if respJSON.Error == "" {
		respJSON.Error = "Error unknown"
	}

	return fmt.Errorf("%d: %s", response.StatusCode, respJSON.Error)
}
