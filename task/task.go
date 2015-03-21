package task

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
	"io/ioutil"
	"time"
)

type Task struct {
	Id          string
	Expressions []string
}

type TaskSlice struct {
	Tasks []Task
}

type Job interface {
	doInterval(a ...interface{}) bool
}

func InitTasks() {
	filePath := fmt.Sprintf("%s/conf/%s", beego.AppPath, "task.json")
	beego.BeeLogger.Debug("Init scheduler with configuration file %s...", filePath)

	b, err := ioutil.ReadFile(filePath)
	if nil != err {
		beego.BeeLogger.Error("Unable to open/read file ---- [%s]", filePath)
		return
	}

	var taskSlice TaskSlice
	json.Unmarshal(b, &taskSlice)

	var c *cron.Cron
	for _, task := range taskSlice.Tasks {
		fmt.Println(task)
		if len(task.Expressions) == 0 {
			continue
		}
		c = cron.New()
		for _, expression := range task.Expressions {
			c.AddFunc(expression, p)
		}
		c.Start()
	}
}

func p() {
	fmt.Println(fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")))
}
