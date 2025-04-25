package container

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

func getContainers(w http.ResponseWriter, r *http.Request) {
	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_container")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	cs, err := bootstrap.NewContainerService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.ContainerServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	containers, err := cs.GetAll()
	if err != nil {
		l.Warn(errs.ErrMsg("cannot get all containers", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetResponse(containers).Finish(w, r, l)
}

func createContainer(w http.ResponseWriter, r *http.Request) {
	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_container")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	var requestContainer model.Container
	if err := route.ReadPostData(r, &requestContainer); err != nil {
		l.Warn(errs.ErrMsg(msg.RequestParse, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	requestContainer.ID = ""

	cs, err := bootstrap.NewContainerService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.ContainerServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	newContainer, err := cs.Create(&requestContainer)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot create container", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetResponse(newContainer).Finish(w, r, l)
}

func startContainer(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	containerID := route.ReadURLParam("id", r)

	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_container")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	cs, err := bootstrap.NewContainerService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.ContainerServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	err = cs.StartContainer(containerID)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot start container", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetSuccess().Finish(w, r, l)
}

func stopContainer(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	containerID := route.ReadURLParam("id", r)

	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_container")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	cs, err := bootstrap.NewContainerService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.ContainerServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	err = cs.StopContainer(containerID)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot stop container", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetSuccess().Finish(w, r, l)
}

func deleteContainer(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	containerID := route.ReadURLParam("id", r)

	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_container")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	cs, err := bootstrap.NewContainerService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.ContainerServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	err = cs.DeleteContainer(containerID)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot delete container", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetSuccess().Finish(w, r, l)
}
