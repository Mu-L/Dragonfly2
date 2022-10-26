/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import "d7y.io/dragonfly/v2/scheduler/resource"

type evaluator struct{}

func (e *evaluator) Evaluate(parent *resource.Peer, child *resource.Peer, taskPieceCount int32) float64 {
	return float64(1)
}

func (e *evaluator) IsBadNode(peer *resource.Peer) bool {
	return true
}

func (e *evaluator) EvalType() string {
	return "plugin"
}

func DragonflyPluginInit(option map[string]string) (any, map[string]string, error) {
	return &evaluator{}, map[string]string{"type": "scheduler", "name": "evaluator"}, nil
}
