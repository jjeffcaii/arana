/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dal

import (
	"context"
	"fmt"
	"strings"
)

import (
	"github.com/arana-db/arana/pkg/dataset"
	"github.com/arana-db/arana/pkg/mysql/rows"
	"github.com/arana-db/arana/pkg/mysql/thead"
	"github.com/arana-db/arana/pkg/proto"
	"github.com/arana-db/arana/pkg/proto/rule"
	"github.com/arana-db/arana/pkg/resultx"
	"github.com/arana-db/arana/pkg/runtime/ast"
	"github.com/arana-db/arana/pkg/runtime/plan"
)

var _ proto.Plan = (*ShowDatabaseRulesPlan)(nil)

type ShowDatabaseRulesPlan struct {
	plan.BasePlan
	Stmt *ast.ShowDatabaseRule
	rule *rule.Rule
}

func (s *ShowDatabaseRulesPlan) Type() proto.PlanType {
	return proto.PlanTypeQuery
}

func (s *ShowDatabaseRulesPlan) ExecIn(ctx context.Context, _ proto.VConn) (proto.Result, error) {
	_, span := plan.Tracer.Start(ctx, "ShowDatabaseRulesPlan.ExecIn")
	defer span.End()

	fields := thead.DBRule.ToFields()
	ds := &dataset.VirtualDataset{
		Columns: fields,
	}

	vt, ok := s.rule.VTable(s.Stmt.TableName)
	if !ok {
		return resultx.New(resultx.WithDataset(ds)), nil
	}

	for _, vs := range vt.GetVShards() {
		var columns []string
		for i := range vs.DB.ShardColumns {
			columns = append(columns, vs.DB.ShardColumns[i].Name)
		}
		ds.Rows = append(ds.Rows, rows.NewTextVirtualRow(fields, []proto.Value{
			proto.NewValueString(s.Stmt.TableName),
			proto.NewValueString(strings.Join(columns, ",")),
			proto.NewValueString(""),
			proto.NewValueString(fmt.Sprintf("%s", vs.DB.Computer)),
			proto.NewValueInt64(1),
		}))
	}
	return resultx.New(resultx.WithDataset(ds)), nil
}

func (s *ShowDatabaseRulesPlan) SetRule(rule *rule.Rule) {
	s.rule = rule
}

// NewShowDatabaseRulesPlan create ShowDatabaseRules Plan
func NewShowDatabaseRulesPlan(stmt *ast.ShowDatabaseRule) *ShowDatabaseRulesPlan {
	return &ShowDatabaseRulesPlan{
		Stmt: stmt,
	}
}
