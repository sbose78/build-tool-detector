/*

Package controllers is autogenerated
and containing scaffold outputs
as well as manually created sub-packages
and files.

*/
package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/goadesign/goa"
	"github.com/tinakurian/build-tool-detector/app"
	"github.com/tinakurian/build-tool-detector/controllers/buildtype"
	errs "github.com/tinakurian/build-tool-detector/controllers/error"
	"github.com/tinakurian/build-tool-detector/controllers/git"
	"github.com/tinakurian/build-tool-detector/domain/system"
	logorus "github.com/tinakurian/build-tool-detector/log"
)

var (
	// ErrInternalServerErrorFailedJSONMarshal unable to marshal json.
	ErrInternalServerErrorFailedJSONMarshal = errors.New("unable to marshal json")

	// ErrInternalServerErrorFailedPropagate unable to propagate error
	ErrInternalServerErrorFailedPropagate = errors.New("unable to propagate error")
)

const (
	errorz                      = "error"
	contentType                 = "Content-Type"
	applicationJSON             = "application/json"
	buildToolDetectorController = "BuildToolDetectorController"
)

// BuildToolDetectorController implements the build-tool-detector resource.
type BuildToolDetectorController struct {
	*goa.Controller
	ghClientID     string
	ghClientSecret string
}

// NewBuildToolDetectorController creates a build-tool-detector controller.
func NewBuildToolDetectorController(service *goa.Service, ghClientID string, ghClientSecret string) *BuildToolDetectorController {
	return &BuildToolDetectorController{Controller: service.NewController(buildToolDetectorController), ghClientID: ghClientID, ghClientSecret: ghClientSecret}
}

// Show runs the show action.
func (c *BuildToolDetectorController) Show(ctx *app.ShowBuildToolDetectorContext) error {
	rawURL := ctx.URL
	_, err := git.GetGitServiceType(rawURL)
	if err != nil {
		return handleRequest(ctx, err, nil)
	}

	gitService := system.System{}.GetGitService()
	buildToolType, err := gitService.GetGitHubService(c.ghClientID, c.ghClientSecret).GetContents(ctx.Context, rawURL, ctx.Branch)
	if err != nil {
		if err.StatusCode == http.StatusBadRequest {
			return handleRequest(ctx, err, nil)
		}
		return handleRequest(ctx, err, buildToolType)
	}

	return handleRequest(ctx, nil, buildToolType)
}

// handleRequest handles returning the correct goa context as well as the GoaBuildToolDetector response
func handleRequest(ctx *app.ShowBuildToolDetectorContext, httpTypeError *errs.HTTPTypeError, buildToolType *string) error {
	ctx.ResponseWriter.Header().Set(contentType, applicationJSON)
	if (httpTypeError == nil || httpTypeError.StatusCode == http.StatusInternalServerError) && buildToolType != nil {
		buildTool := buildtype.Unknown()
		if buildtype.MAVEN == *buildToolType {
			buildTool = buildtype.Maven()
		}
		return ctx.OK(buildTool)
	}

	ctx.WriteHeader(httpTypeError.StatusCode)
	jsonHTTPTypeError, err := json.Marshal(httpTypeError)
	if err != nil {
		logorus.Logger().WithError(err).WithField(errorz, httpTypeError).Errorf(ErrInternalServerErrorFailedJSONMarshal.Error())
		return ctx.InternalServerError()
	}

	if _, err := fmt.Fprint(ctx.ResponseWriter, string(jsonHTTPTypeError)); err != nil {
		logorus.Logger().WithError(err).WithField(errorz, jsonHTTPTypeError).Errorf(ErrInternalServerErrorFailedPropagate.Error())
		return ctx.InternalServerError()
	}

	return getErrResponse(ctx, httpTypeError)
}

// getErrResponse will determine the correct goa error response
func getErrResponse(ctx *app.ShowBuildToolDetectorContext, httpTypeError *errs.HTTPTypeError) error {
	var response error
	switch httpTypeError.StatusCode {
	case http.StatusBadRequest:
		response = ctx.BadRequest()
	case http.StatusInternalServerError:
		response = ctx.InternalServerError()
	}

	return response
}
