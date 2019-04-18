package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func Test_logicInstall(t *testing.T) {
	params := domain.InstallParams{
		ConfigFilePath: "/etc/akoi/akoi.yml", Format: "ansible"}
	lgc := newLogicMock(t)
	result := lgc.Install(context.Background(), params)
	if result.Failed() {
		t.Fatal(result.String("ansible"))
	}
	lgc.CfgReader = test.NewConfigReader(t, gomic.DoNothing).
		SetReturnRead(domain.Config{}, fmt.Errorf("failed to read config"))
	result = lgc.Install(context.Background(), params)
	if !result.Failed() {
		t.Fatal("it should be failed to read config")
	}
}
