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

import React from 'react';
import {
  // ViewVariableType,
  ValueExpressionType,
  type RefExpression,
  ViewVariableType,
  type LiteralExpression,
  type InputValueVO,
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
  params?: any
}

export const McpParamsField = withFieldArray(({
  params,
  disabledTypes,
}: McpParamsFieldProps) => {
  const { value, remove, append } = useFieldArray<InputValueVO>();
  console.log(value, '---------------------------456')
  return (
    <Section
      title={I18n.t('workflow_detail_node_input')}
      tooltip={I18n.t(
        'node_http_request_params_desc',
        {},
        '输入参数值',
      )}
    >
       <div onClick={() => remove(0)}>remove</div>
      <div onClick={() => {
        append({
          name: `test` + new Date().getTime().toString(),
          input: { type: ValueExpressionType.REF },
        });
      }}>add</div>
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
            key={index}
            label={name}
            required={false}
            inputType={ViewVariableType.String}
            disabledTypes={disabledTypes}
            name={`inputParameters.${index}.input`}
          />
        ))}
      </FieldArrayList>
    </Section>
  )
});
