/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { useEffect } from 'react';
import {
  // ViewVariableType,
  ValueExpressionType,
  type RefExpression,
  ViewVariableType,
  type InputValueVO,
  useWorkflowNode,
} from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';

import { ColumnsTitle } from '@/form-extensions/components/columns-title';
import { ValueExpressionInputField } from '@/node-registries/common/fields';
import {
  Section,
  useFieldArray,
  FieldArrayList,
  withFieldArray,
} from '@/form';

const mcpTypeObj = {
  "string": ViewVariableType.String,
  "integer": ViewVariableType.Integer,
  "number": ViewVariableType.Number,
  "boolean": ViewVariableType.Boolean,
  "object": ViewVariableType.Object,
  "time": ViewVariableType.Time,
  "arraystring": ViewVariableType.ArrayString,
  "arrayinteger": ViewVariableType.ArrayInteger,
  "arraynumber": ViewVariableType.ArrayNumber,
  "arrayboolean": ViewVariableType.ArrayBoolean,
  "arrayobject": ViewVariableType.ArrayObject,
  "arraytime": ViewVariableType.ArrayTime,
}

interface McpParamsFieldProps {
  disabledTypes?: ViewVariableType[];
  defaultValue?: RefExpression;
  name?: string;
  params?: any;
}

export const McpParamsField = withFieldArray(({
  params,
  disabledTypes,
}: McpParamsFieldProps) => {
  const { value, remove, append } = useFieldArray<InputValueVO>();
  const { data } = useWorkflowNode();
  const mcpList = data?.inputs?.mcpInfoList || []
  const { properties = {}, required = [] } = mcpList?.[0]?.inputSchema || {}
  const mcpParamsList = Object.keys(properties || {}) || []

  const removeAll = () => {
    const valueArr = JSON.parse(JSON.stringify(value || []))
    valueArr.forEach(() => {
      remove(0)
    })
  }

  const appendValue = (needAddList: any) => {
    needAddList?.forEach((item: any) => {
      append({
        name: item,
        input: { type: ValueExpressionType.REF },
      })
    })
  }

  const initAppend = (params: any) => {
    // 首次进入，但是有选择 mcp，没有 value 或者 value 不全的情况，添加 value
    if (params === null && Boolean(mcpParamsList.length)) {
      if (!value?.length) {
        appendValue(mcpParamsList)
      } else if (Boolean(value?.length) && (value?.length < mcpParamsList.length)) {
        const needAddList = mcpParamsList.filter(item => !value?.map(it => it.name).includes(item))
        appendValue(needAddList)
      }
    }
  }

  useEffect(() => {
    // 如果 params 非首次进入，通过手动添加的，清之前的数据，加最新数据
    if (params !== null && params) {
      removeAll()
      appendValue(mcpParamsList)
    }
    // 首次进入的情况
    initAppend(params)
  }, [params]);

  return mcpList?.length > 0 ? (
    <Section
      title={I18n.t('workflow_detail_node_input')}
      tooltip={I18n.t(
        'node_http_request_params_desc',
        {},
        '输入参数值',
      )}
    >
      <ColumnsTitle
        columns={[
          {
            title: I18n.t('workflow_detail_node_parameter_name'),
            style: {flex: 2},
          },
          {
            title: I18n.t('workflow_detail_end_output_value'),
            style: {flex: 3},
          },
        ]}
        className="mb-[8px]"
      />
      <FieldArrayList>
        {value?.map(({name}, index) => (
          <ValueExpressionInputField
            key={name + index}
            label={name}
            required={required.includes(name)}
            inputType={
              mcpTypeObj[properties[name]?.type + (properties[name]?.items?.type || '')] || ViewVariableType.String
            }
            disabledTypes={disabledTypes}
            name={`inputs.inputParameters.${index}.input`}
          />
        ))}
      </FieldArrayList>
    </Section>
  ) : null
});
