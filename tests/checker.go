package tests

import (
	"github.com/bozaro/tech-db-forum/client"
	"github.com/go-openapi/runtime"
	http_transport "github.com/go-openapi/runtime/client"
	"log"
	"net/http"
)

type Checker struct {
	// Имя текущей проверки.
	Name string
	// Описание текущей проверки.
	Description string
	// Функция для текущей проверки.
	FnCheck func(c *client.Forum)
	// Тесты, без которых проверка не имеет смысл.
	Deps []string
}

type CheckerTransport struct {
	t      runtime.ClientTransport
	report *Report
}

func (self *CheckerTransport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	tracker := NewValidator(operation.Context, self.report)
	operation.Client = &http.Client{Transport: tracker}
	return self.t.Submit(operation)
}

func Checkpoint(c *client.Forum, message string) {
	c.Transport.(*CheckerTransport).report.Checkpoint(message)
}

var checks []Checker

func Register(checker Checker) {
	checks = append(checks, checker)
}

func RunCheck(check Checker, report *Report) {
	report.result = REPORT_SUCCESS
	cfg := client.DefaultTransportConfig().WithHost("localhost:5000").WithSchemes([]string{"http"})
	transport := http_transport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	defer func() {
		if r := recover(); r != nil {
			report.AddError(r)
		}
	}()
	check.FnCheck(client.New(&CheckerTransport{transport, report}, nil))
}

func Run() {
	for _, check := range checks {
		log.Printf("=== RUN:  %s", check.Name)
		report := Report{}
		RunCheck(check, &report)
		if report.result != REPORT_SUCCESS {
			report.Show()
		}
		log.Printf("--- DONE: %s (%d)", check.Name, report.result)
	}
}