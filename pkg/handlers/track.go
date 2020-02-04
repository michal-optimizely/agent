/****************************************************************************
 * Copyright 2019-2020, Optimizely, Inc. and contributors                        *
 *                                                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");          *
 * you may not use this file except in compliance with the License.         *
 * You may obtain a copy of the License at                                  *
 *                                                                          *
 *    http://www.apache.org/licenses/LICENSE-2.0                            *
 *                                                                          *
 * Unless required by applicable law or agreed to in writing, software      *
 * distributed under the License is distributed on an "AS IS" BASIS,        *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 * See the License for the specific language governing permissions and      *
 * limitations under the License.                                           *
 ***************************************************************************/

// Package handlers //
package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/optimizely/go-sdk/pkg/entities"

	"github.com/optimizely/agent/pkg/middleware"
)

type trackBody struct {
	userID         string
	userAttributes map[string]interface{}
	eventTags      map[string]interface{}
}

// TrackEvent - track a given event for the current user
func TrackEvent(w http.ResponseWriter, r *http.Request) {
	optlyClient, _, err := parseContext(r)
	if err != nil {
		RenderError(err, http.StatusInternalServerError, w, r)
		return
	}

	var body trackBody
	err = ParseRequestBody(r, &body)
	if err != nil {
		RenderError(err, http.StatusBadRequest, w, r)
		return
	}

	eventKey := r.URL.Query().Get("eventKey")
	if eventKey == "" {
		err = fmt.Errorf("missing required path parameter: eventKey")
		RenderError(err, http.StatusBadRequest, w, r)
		return
	}

	uc := entities.UserContext{
		ID:         body.userID,
		Attributes: body.userAttributes,
	}

	err = optlyClient.Track(eventKey, uc, body.eventTags)

	if err != nil {
		middleware.GetLogger(r).Error().Err(err).Str("eventKey", eventKey).Msg("error tracking event")
		RenderError(err, http.StatusNotFound, w, r)
		return
	}
	middleware.GetLogger(r).Debug().Str("eventKey", eventKey).Msg("tracking event")
	render.NoContent(w, r)
}
