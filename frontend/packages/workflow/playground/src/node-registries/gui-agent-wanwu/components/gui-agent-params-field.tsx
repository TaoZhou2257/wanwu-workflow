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
  type InputValueVO,
} from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';
import { DEFAULT_PARAMS_LIST } from '../constants';

import { ColumnsTitle } from '@/form-extensions/components/columns-title';
import { ValueExpressionInputField } from '@/node-registries/common/fields';
import {
  Section,
  useFieldArray,
  FieldArrayList,
  withFieldArray,
} from '@/form';

export const GuiAgentParamsField = withFieldArray(() => {
  const { value, remove, append } = useFieldArray<InputValueVO>();

  const removeAll = () => {
    const valueArr = JSON.parse(JSON.stringify(value || []))
    valueArr.forEach(() => {
      remove(0)
    })
  }

  const appendValue = (appendList: any) => {
    appendList?.forEach((item: any) => {
      append(item)
    })
  }

  const initAppend = () => {
    // 进入更新数据
    const appendList = DEFAULT_PARAMS_LIST.map(item => {
      if (value?.map(it => it.name).includes(item.name)) {
        return value[value.findIndex(it => it.name === item.name)];
      } else {
        return item
      }
    })
    removeAll()
    appendValue(appendList)
  }

  useEffect(() => {
    initAppend()
  }, []);

  return (
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
        {value?.map(({name}, index) => {
          const { required, type } = DEFAULT_PARAMS_LIST.find(
            item => item.name === name
          ) || {}
          return (
            <ValueExpressionInputField
              key={name + index}
              label={name}
              required={required}
              inputType={type}
              // layout={'vertical'}
              name={`inputs.inputParameters.${index}.input`}
            />
          )
        })}
      </FieldArrayList>
    </Section>
  )
});
