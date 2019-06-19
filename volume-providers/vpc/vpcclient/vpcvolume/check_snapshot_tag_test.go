/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpcvolume_test

import (
	//"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/riaas/test"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/vpcvolume"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestCheckSnapshotTag(t *testing.T) {
	// Setup new style zap logger
	logger, _ := GetTestContextLogger()
	defer logger.Sync()

	testCases := []struct {
		name string

		// backend url
		url string

		// Response
		status  int
		content string

		// Expected return
		expectErr string
		verify    func(*testing.T, error)
	}{
		{
			name:   "Verify that the correct endpoint is invoked",
			status: http.StatusNoContent,
			url:    vpcvolume.Version + "/volumes/volume1/snapshots/snapshot1/tags/test-tag",
		}, {
			name:      "Verify that a 404 is returned to the caller",
			status:    http.StatusNotFound,
			url:       vpcvolume.Version + "/volumes/volume1/snapshots/snapshot1/tags/test-tag",
			content:   "{\"errors\":[{\"message\":\"testerr\"}]}",
			expectErr: "Trace Code:, testerr Please check ",
		}, {
			name:    "Verify that the snapshot is parsed correctly",
			status:  http.StatusNotFound,
			url:     vpcvolume.Version + "/volumes/volume1/snapshots/snapshot1/tags/test-tag",
			content: "{\"id\":\"snapshot1\",\"name\":\"snapshot1\",\"status\":\"pending\"}",
			verify: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		}, {
			name:    "Verify that the snapshot tag exists",
			status:  http.StatusOK,
			url:     vpcvolume.Version + "/volumes/volume1/snapshots/snapshot1/tags/test-tag",
			content: "{\"id\":\"snapshot1\",\"name\":\"snapshot1\",\"status\":\"pending\", \"tags\":[\"test-tag\"]}",
			verify: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			mux, client, teardown := test.SetupServer(t)
			emptyString := ""
			test.SetupMuxResponse(t, mux, testcase.url, http.MethodGet, &emptyString, testcase.status, testcase.content, nil)

			defer teardown()

			logger.Info("Test case being executed", zap.Reflect("testcase", testcase.name))

			snapshotService := vpcvolume.NewSnapshotManager(client)

			err := snapshotService.CheckSnapshotTag("volume1", "snapshot1", "test-tag", logger)

			if testcase.verify != nil {
				testcase.verify(t, err)
			}
			// vpc snapshot functionality is not yet ready. It would return error for now
		})
	}
}
