package template

import (
	"net/http"

	"github.com/kaibling/apiforge/apierror"
	"github.com/kaibling/apiforge/envelope"
	"github.com/kaibling/apiforge/route"
	"github.com/kaibling/cerodev/bootstrap"
	"github.com/kaibling/cerodev/errs"
	"github.com/kaibling/cerodev/errs/msg"
	"github.com/kaibling/cerodev/model"
)

func getTemplates(w http.ResponseWriter, r *http.Request) {
	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_template")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	cs, err := bootstrap.NewTemplateService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.TemplateServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	templates, err := cs.GetAll()
	if err != nil {
		l.Warn(errs.ErrMsg("cannot get all tokens", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetResponse(templates).Finish(w, r, l)
}

func createTemplate(w http.ResponseWriter, r *http.Request) {
	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_template")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	var requestTemplate model.Template
	if err := route.ReadPostData(r, &requestTemplate); err != nil {
		l.Warn(errs.ErrMsg(msg.RequestParse, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	requestTemplate.ID = ""

	cs, err := bootstrap.NewTemplateService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.TemplateServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	newTemplate, err := cs.Create(&requestTemplate)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot create template", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetResponse(newTemplate).Finish(w, r, l)
}

func deleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := route.ReadURLParam("id", r)

	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_template")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	cs, err := bootstrap.NewTemplateService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.TemplateServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	err = cs.Delete(templateID)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot delete template", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetSuccess().Finish(w, r, l)
}

func buildImage(w http.ResponseWriter, r *http.Request) {
	templateID := route.ReadURLParam("id", r)

	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_template")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	var buildParams model.BuildParams
	if err := route.ReadPostData(r, &buildParams); err != nil {
		l.Warn(errs.ErrMsg(msg.RequestParse, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	buildParams.Validate()
	buildParams.TemplateID = templateID

	ctrs, err := bootstrap.NewContainerService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.TemplateServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	err = ctrs.BuildTemplate(buildParams.TemplateID, buildParams.Tag, buildParams.BuildArgs)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot build template", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetSuccess().Finish(w, r, l)
}

func updateTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := route.ReadURLParam("id", r)

	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_template")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	var template model.Template
	if err := route.ReadPostData(r, &template); err != nil {
		l.Warn(errs.ErrMsg(msg.RequestParse, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	template.ID = templateID

	ts, err := bootstrap.NewTemplateService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.TemplateServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	newTemplate, err := ts.Update(&template)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot update template", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetResponse(newTemplate).Finish(w, r, l)
}
