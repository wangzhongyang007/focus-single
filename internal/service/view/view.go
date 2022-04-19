package view

import (
	"context"
	"fmt"

	"focus-single/internal/service/category"
	"focus-single/internal/service/menu"
	"focus-single/internal/service/session"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmode"

	"focus-single/internal/model"
)

// 前台系统-获取面包屑列表
func GetBreadCrumb(ctx context.Context, in *model.ViewGetBreadCrumbInput) []model.ViewBreadCrumb {
	breadcrumb := []model.ViewBreadCrumb{
		{Name: "首页", Url: "/"},
	}
	var uriPrefix string
	if in.ContentType != "" {
		uriPrefix = "/" + in.ContentType
		topMenuItem, _ := menu.GetTopMenuByUrl(ctx, uriPrefix)
		if topMenuItem != nil {
			breadcrumb = append(breadcrumb, model.ViewBreadCrumb{
				Name: topMenuItem.Name,
				Url:  topMenuItem.Url,
			})
		}
	}
	if uriPrefix != "" && in.CategoryId > 0 {
		categoryEntity, _ := category.GetItem(ctx, in.CategoryId)
		if categoryEntity != nil {
			breadcrumb = append(breadcrumb, model.ViewBreadCrumb{
				Name: categoryEntity.Name,
				Url:  fmt.Sprintf(`%s?cate=%d`, uriPrefix, categoryEntity.Id),
			})
		}
	}
	if in.ContentId > 0 {
		breadcrumb = append(breadcrumb, model.ViewBreadCrumb{
			Name: "内容详情",
		})
	}
	return breadcrumb
}

// 前台系统-获取标题
func GetTitle(ctx context.Context, in *model.ViewGetTitleInput) string {
	var (
		titleArray []string
		uriPrefix  string
	)
	if in.CurrentName != "" {
		titleArray = append(titleArray, in.CurrentName)
	}
	if in.CategoryId > 0 {
		categoryEntity, _ := category.GetItem(ctx, in.CategoryId)
		if categoryEntity != nil {
			titleArray = append(titleArray, categoryEntity.Name)
		}
	}
	if in.ContentType != "" {
		uriPrefix = "/" + in.ContentType
		topMenuItem, _ := menu.GetTopMenuByUrl(ctx, uriPrefix)
		if topMenuItem != nil {
			titleArray = append(titleArray, topMenuItem.Name)
		}
	}
	return gstr.Join(titleArray, " - ")
}

// 渲染指定模板页面
func RenderTpl(ctx context.Context, tpl string, data ...model.View) {
	var (
		viewObj  = model.View{}
		viewData = make(g.Map)
		request  = g.RequestFromCtx(ctx)
	)
	if len(data) > 0 {
		viewObj = data[0]
	}
	if viewObj.Title == "" {
		viewObj.Title = g.Cfg().MustGet(ctx, `setting.title`).String()
	} else {
		viewObj.Title = viewObj.Title + ` - ` + g.Cfg().MustGet(ctx, `setting.title`).String()
	}
	if viewObj.Keywords == "" {
		viewObj.Keywords = g.Cfg().MustGet(ctx, `setting.keywords`).String()
	}
	if viewObj.Description == "" {
		viewObj.Description = g.Cfg().MustGet(ctx, `setting.description`).String()
	}
	// 去掉空数据
	viewData = gconv.Map(viewObj)
	for k, v := range viewData {
		if g.IsEmpty(v) {
			delete(viewData, k)
		}
	}
	// 内置对象
	viewData["BuildIn"] = &viewBuildIn{httpRequest: request}
	// 内容模板
	if viewData["MainTpl"] == nil {
		viewData["MainTpl"] = getDefaultMainTpl(ctx)
	}
	// 提示信息
	if notice, _ := session.GetNotice(ctx); notice != nil {
		_ = session.RemoveNotice(ctx)
		viewData["Notice"] = notice
	}
	// 渲染模板
	_ = request.Response.WriteTpl(tpl, viewData)
	// 开发模式下，在页面最下面打印所有的模板变量
	if gmode.IsDevelop() {
		_ = request.Response.WriteTplContent(`{{dump .}}`, viewData)
	}
}

// 渲染默认模板页面
func Render(ctx context.Context, data ...model.View) {
	RenderTpl(ctx, g.Cfg().MustGet(ctx, "viewer.indexLayout").String(), data...)
}

// 跳转中间页面
func Render302(ctx context.Context, data ...model.View) {
	view := model.View{}
	if len(data) > 0 {
		view = data[0]
	}
	if view.Title == "" {
		view.Title = "页面跳转中"
	}
	view.MainTpl = getViewFolderName(ctx) + "/pages/302.html"
	Render(ctx, view)
}

// 401页面
func Render401(ctx context.Context, data ...model.View) {
	view := model.View{}
	if len(data) > 0 {
		view = data[0]
	}
	if view.Title == "" {
		view.Title = "无访问权限"
	}
	view.MainTpl = getViewFolderName(ctx) + "/pages/401.html"
	Render(ctx, view)
}

// 403页面
func Render403(ctx context.Context, data ...model.View) {
	view := model.View{}
	if len(data) > 0 {
		view = data[0]
	}
	if view.Title == "" {
		view.Title = "无访问权限"
	}
	view.MainTpl = getViewFolderName(ctx) + "/pages/403.html"
	Render(ctx, view)
}

// 404页面
func Render404(ctx context.Context, data ...model.View) {
	view := model.View{}
	if len(data) > 0 {
		view = data[0]
	}
	if view.Title == "" {
		view.Title = "资源不存在"
	}
	view.MainTpl = getViewFolderName(ctx) + "/pages/404.html"
	Render(ctx, view)
}

// 500页面
func Render500(ctx context.Context, data ...model.View) {
	view := model.View{}
	if len(data) > 0 {
		view = data[0]
	}
	if view.Title == "" {
		view.Title = "请求执行错误"
	}
	view.MainTpl = getViewFolderName(ctx) + "/pages/500.html"
	Render(ctx, view)
}

// 获取视图存储目录
func getViewFolderName(ctx context.Context) string {
	return gstr.Split(g.Cfg().MustGet(ctx, "viewer.indexLayout").String(), "/")[0]
}

// 获取自动设置的MainTpl
func getDefaultMainTpl(ctx context.Context) string {
	var (
		viewFolderPrefix = getViewFolderName(ctx)
		urlPathArray     = gstr.SplitAndTrim(g.RequestFromCtx(ctx).URL.Path, "/")
		mainTpl          string
	)
	if len(urlPathArray) > 0 && urlPathArray[0] == viewFolderPrefix {
		urlPathArray = urlPathArray[1:]
	}
	switch {
	case len(urlPathArray) == 2:
		// 如果2级路由为数字，那么为模块的详情页面，那么路由固定为/xxx/detail。
		// 如果需要定制化内容模板，请在具体路由方法中设置MainTpl。
		if gstr.IsNumeric(urlPathArray[1]) {
			urlPathArray[1] = "detail"
		}
		mainTpl = viewFolderPrefix + "/" + gfile.Join(urlPathArray[0], urlPathArray[1]) + ".html"
	case len(urlPathArray) == 1:
		mainTpl = viewFolderPrefix + "/" + urlPathArray[0] + "/index.html"
	default:
		// 默认首页内容
		mainTpl = viewFolderPrefix + "/index/index.html"
	}
	return gstr.TrimLeft(mainTpl, "/")
}
