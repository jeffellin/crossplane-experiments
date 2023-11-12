package main

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	"os"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/resource"
)

func TestRunFunction(t *testing.T) {
	type args struct {
		ctx context.Context
		req *fnv1beta1.RunFunctionRequest
	}
	type want struct {
		rsp *fnv1beta1.RunFunctionResponse
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"Add status binding name": {
			reason: "The Function should the status binding name",
			args: args{
				req: &fnv1beta1.RunFunctionRequest{
					Observed: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							// MustStructJSON is a handy way to provide mock
							// resources.
							Resource: resource.MustStructJSON(`{
								"apiVersion": "ellin.net/v1alpha1",
								"kind": "XMyHelmishDataStore",
								"metadata": {
									"name": "test"
								},
								"spec": {
									"enablePersistence": "True",
									"databasename": "petclinic",
									"writeConnectionSecretToRef": {
										"name": "brigade-helmish-jellin",
                                         "namespace": "other-namespace"
									}
								}
							}`),
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Meta: &fnv1beta1.ResponseMeta{Ttl: durationpb.New(60 * time.Second)},
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"status": {
                                  "binding": {
								     "name":"brigade-helmish-jellin"
								   }
								}
							}`),
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := &Function{log: logging.NewNopLogger()}
			rsp, err := f.RunFunction(tc.args.ctx, tc.args.req)
			if err == nil {

			}
			fmt.Fprintln(os.Stdout, spew.Sdump(tc.want.rsp))
			fmt.Fprintln(os.Stdout, "---")
			fmt.Fprintln(os.Stdout, spew.Sdump(rsp))

			//		spew.Dump(tc.want.rsp)
			//t.Log("in a test")
			//t.Log(spew.Sdump(rsp))
			if diff := cmp.Diff(tc.want.rsp, rsp, protocmp.Transform()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want rsp, +got rsp:\n%s", tc.reason, diff)
			}

			if diff := cmp.Diff(tc.want.err, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want err, +got err:\n%s", tc.reason, diff)
			}
		})
	}
}
