package images

import (
	"net/http"

	"github.com/kaibling/apiforge/envelope"
	"github.com/kaibling/cerodev/api/apierrs"
	"github.com/kaibling/cerodev/bootstrap"
	"github.com/kaibling/cerodev/errs"
	"github.com/kaibling/cerodev/errs/msg"
)

func getImages(w http.ResponseWriter, r *http.Request) {
	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_images")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	cs, err := bootstrap.NewContainerService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.ContainerServiceName, err))
		e.SetError(apierrs.HandleError(err)).Finish(w, r, l)

		return
	}

	images, err := cs.GetImages()
	if err != nil {
		l.Warn(errs.ErrMsg("cannot get images", err))
		e.SetError(apierrs.HandleError(err)).Finish(w, r, l)

		return
	}

	e.SetResponse(images).Finish(w, r, l)
}
