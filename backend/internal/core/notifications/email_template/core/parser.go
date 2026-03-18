package core

import (
	"bytes"
	htmltemplate "html/template"
	texttemplate "text/template"

	"backend/core/notifications/email_template/port"
	apperrors "backend/port/errors"
	"github.com/samber/oops"
)

func (s service) ParseTemplate(tmpl port.EmailTemplate, data any) (port.ParsedTemplate, error) {
	subjectTmpl, err := texttemplate.New(tmpl.Name + "_subject").Parse(tmpl.Subject)
	if err != nil {
		return port.ParsedTemplate{}, oops.In(apperrors.LayerService).Wrapf(err, "failed to parse subject template %q", tmpl.Name)
	}

	var subjectBuf bytes.Buffer
	if err = subjectTmpl.Execute(&subjectBuf, data); err != nil {
		return port.ParsedTemplate{}, oops.In(apperrors.LayerService).Wrapf(err, "failed to execute subject template %q", tmpl.Name)
	}

	contentTmpl, err := htmltemplate.New(tmpl.Name + "_content").Parse(tmpl.Content)
	if err != nil {
		return port.ParsedTemplate{}, oops.In(apperrors.LayerService).Wrapf(err, "failed to parse content template %q", tmpl.Name)
	}

	var contentBuf bytes.Buffer
	if err = contentTmpl.Execute(&contentBuf, data); err != nil {
		return port.ParsedTemplate{}, oops.In(apperrors.LayerService).Wrapf(err, "failed to execute content template %q", tmpl.Name)
	}

	return port.ParsedTemplate{
		Subject: subjectBuf.String(),
		Content: contentBuf.String(),
	}, nil
}
