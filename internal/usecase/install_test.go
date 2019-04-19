package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
	"github.com/suzuki-shunsuke/akoi/internal/test"
)

func Test_logicInstall(t *testing.T) {
	params := domain.InstallParams{
		ConfigFilePath: "/etc/akoi/akoi.yml", Format: "ansible"}
	lgc := newLogicMock(t)
	result := lgc.Install(context.Background(), params)
	require.False(t, result.Failed())
	lgc.CfgReader = test.NewConfigReader(t, gomic.DoNothing).
		SetReturnRead(domain.Config{}, fmt.Errorf("failed to read config"))
	result = lgc.Install(context.Background(), params)
	require.True(t, result.Failed())
}
