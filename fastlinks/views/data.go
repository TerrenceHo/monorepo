package views

type TemplateType int64

const (
	ROOT_VIEW TemplateType = 0
	NEW_VIEW  TemplateType = 1
)

type TemplateData struct {
	Type TemplateType
	Data interface{}
}
