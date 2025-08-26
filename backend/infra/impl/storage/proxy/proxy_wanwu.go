package proxy

import (
	"context"
	"net/url"
	"os"
	"path"

	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

// CheckIfNeedReplaceHostByWanwu 参考CheckIfNeedReplaceHost
// 将http://minio-wanwu:9000/abc/def => http://${WANWU_MINIO_PROXY_ENDPOINT}/${WANWU_MINIO_PROXY_PREFIX}/abc/def，需要配合nginx代理
func CheckIfNeedReplaceHostByWanwu(ctx context.Context, originURLStr string) (ok bool, proxyURL string) {
	// url parse
	originURL, err := url.Parse(originURLStr)
	if err != nil {
		logs.CtxWarnf(ctx, "[CheckIfNeedReplaceHost] url parse failed, err: %v", err)
		return false, ""
	}

	proxyEndpoint := os.Getenv("WANWU_MINIO_PROXY_ENDPOINT")
	if proxyEndpoint == "" {
		return false, ""
	}

	proxyPrefix := os.Getenv("WANWU_MINIO_PROXY_PREFIX")

	currentScheme, ok := ctxcache.Get[string](ctx, consts.RequestSchemeKeyInCtx)
	if !ok {
		return false, ""
	}

	originURL.Host = proxyEndpoint
	originURL.Path = path.Join(proxyPrefix, originURL.Path)
	originURL.Scheme = currentScheme
	logs.CtxDebugf(ctx, "[CheckIfNeedReplaceHost] reset originURL.String = %s", originURL.String())
	return true, originURL.String()
}
