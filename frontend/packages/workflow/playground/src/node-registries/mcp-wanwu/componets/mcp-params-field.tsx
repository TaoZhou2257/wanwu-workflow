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

interface McpParamsFieldProps {
  disabledTypes?: ViewVariableType[];
  defaultValue?: RefExpression;
  name?: string;
  params?: any;
  paramsRequiredArr?: any;
}

export const McpParamsField = withFieldArray(({
  paramsRequiredArr,
  params,
  disabledTypes,
}: McpParamsFieldProps) => {
  const { value, remove, append } = useFieldArray<InputValueVO>();
  const { data } = useWorkflowNode();
  const mcpList = data?.inputs?.mcpInfoList || []

  const removeAll = () => {
    const valueArr = JSON.parse(JSON.stringify(value || []))
    valueArr.forEach(() => {
      remove(0)
    })
  }

  const appendAll = () => {
    params.forEach((item: any) => {
      append({
        name: item,
        required: paramsRequiredArr.includes(item),
        input: { type: ValueExpressionType.REF },
      })
    })
  }

  useEffect(() => {
    // 如果 params 非首次进入，通过手动添加的，清之前的数据，加最新数据
    if (params !== null && params) {
      removeAll()
      appendAll()
    }
  }, [params, paramsRequiredArr]);

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
        {value?.map(({name, required}, index) => (
          <ValueExpressionInputField
            key={name + index}
            label={name}
            required={required}
            inputType={ViewVariableType.String}
            disabledTypes={disabledTypes}
            name={`inputs.inputParameters.${index}.input`}
          />
        ))}
      </FieldArrayList>
    </Section>
  ) : null
});
