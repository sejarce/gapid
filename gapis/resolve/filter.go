// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resolve

import (
	"context"

	"github.com/google/gapid/core/data/id"
	"github.com/google/gapid/gapis/atom"
	"github.com/google/gapid/gapis/gfxapi"
	"github.com/google/gapid/gapis/service/path"
)

type filter func(a atom.Atom, s *gfxapi.State) bool

type filters []filter

func (l filters) pass(a atom.Atom, s *gfxapi.State) bool {
	for _, f := range l {
		if !f(a, s) {
			return false
		}
	}
	return true
}

func (l *filters) add(f filter) { *l = append(*l, f) }

func (l *filters) addContextFilter(ctx context.Context, p *path.Context) error {
	c, err := Context(ctx, p)
	if err != nil {
		return err
	}
	id, err := id.Parse(c.Id)
	if err != nil {
		return err
	}
	ctxID := gfxapi.ContextID(id)
	l.add(func(a atom.Atom, s *gfxapi.State) bool {
		if api := a.API(); api != nil {
			if ctx := api.Context(s); ctx != nil {
				return ctx.ID() == ctxID
			}
		}
		return false
	})
	return nil
}