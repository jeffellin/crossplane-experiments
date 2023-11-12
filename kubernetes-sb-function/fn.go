package main

import (
	"context"
	"github.com/crossplane/function-sdk-go/errors"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/davecgh/go-spew/spew"

	//	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"

	"github.com/crossplane/function-sdk-go/response"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.log.Info("I got this far 1")
	rsp := response.To(req, response.DefaultTTL)
	observedCompositeResource, err := request.GetObservedCompositeResource(req)
	desiredCompositeResource, err := request.GetDesiredCompositeResource(req)

	connectionSecret, err := observedCompositeResource.Resource.GetString("spec.writeConnectionSecretToRef.name")

	//spew.Dump(observedCompositeResource)
	if err != nil {
		// If the function can't read the XR, the request is malformed. This
		// should never happen. The function returns a fatal result. This tells
		// Crossplane to stop running functions and return an error.
		response.Fatal(rsp, errors.Wrapf(err, "cannot get observed composite resource from %T", req))
		return rsp, nil
	}
	if err := desiredCompositeResource.Resource.SetString("status.binding.name", connectionSecret); err != nil {
		return nil, err
	}
	desiredComposedResources, err := request.GetDesiredComposedResources(req)

	if err = response.SetDesiredComposedResources(rsp, desiredComposedResources); err != nil {
		return nil, err
	}

	if err = response.SetDesiredCompositeResource(rsp, desiredCompositeResource); err != nil {
		return nil, err
	}
	f.log.Info("I got this far 5")
	spew.Dump(rsp)
	return rsp, nil
}
