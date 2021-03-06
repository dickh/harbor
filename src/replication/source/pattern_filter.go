// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package source

import (
	"regexp"

	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/replication/models"
)

// PatternFilter implements Filter interface for pattern filter
type PatternFilter struct {
	kind      string
	pattern   string
	convertor Convertor
}

// NewPatternFilter returns an instance of PatternFilter
func NewPatternFilter(kind, pattern string, convertor ...Convertor) *PatternFilter {
	filer := &PatternFilter{
		kind:    kind,
		pattern: pattern,
	}

	if len(convertor) > 0 {
		filer.convertor = convertor[0]
	}

	return filer
}

// Init the filter. nil implement for now
func (p *PatternFilter) Init() error {
	return nil
}

// GetConvertor returns the convertor
func (p *PatternFilter) GetConvertor() Convertor {
	return p.convertor
}

// DoFilter filters resources
func (p *PatternFilter) DoFilter(filterItems []models.FilterItem) []models.FilterItem {
	items := []models.FilterItem{}
	for _, item := range filterItems {
		if item.Kind != p.kind {
			log.Warningf("unexpected filter item kind, expected: %s, got: %s, skip",
				p.kind, item.Kind)
			continue
		}

		matched, err := regexp.MatchString(p.pattern, item.Value)
		if err != nil {
			log.Errorf("failed to match pattern %s, value %s: %v, skip",
				p.pattern, item.Value, err)
			continue
		}

		if !matched {
			log.Debugf("%s does not match to the %s filter %s, skip",
				item.Value, p.kind, p.pattern)
			continue
		}

		log.Debugf("add %s to the result of %s filter %s",
			item.Value, p.kind, p.pattern)
		items = append(items, item)
	}

	return items
}
