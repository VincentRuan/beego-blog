package producer

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/bitly/go-nsq"
	"github.com/vincentruan/beego-blog/g"
	"gopkg.in/vmihailenco/msgpack.v2"
)

var producer *nsq.Producer

func InitNSQProducer() error {
	if g.NSQAddr == "" {
		return errors.New("Unable to read NSQ address from config file!")
	}

	var err error
	producer, err = nsq.NewProducer(g.NSQAddr, nsq.NewConfig())
	if err != nil {
		return err
	}
	producer.SetLogger(beego.BeeLogger, nsq.LogLevelInfo)

	return nil
}

func PublishMsg(topic string, v ...interface{}) error {
	var body []byte
	var err error
	if len(v) == 1 {
		body, err = encode(v[0])
		if err != nil {
			return err
		}
		return producer.PublishAsync(topic, body, nil)
	}

	bodies := make([][]byte, len(v))
	for i, vv := range v {
		body, err = encode(vv)
		if err != nil {
			beego.Error(err)
			continue
		}
		bodies[i] = body
	}

	return producer.MultiPublishAsync(topic, bodies, nil)
}

func encode(v interface{}) ([]byte, error) {
	if b, ok := v.([]byte); ok {
		return b, nil
	}

	body, err := msgpack.Marshal(v)
	if err != nil {
		return nil, err
	}
	return body, nil
}
