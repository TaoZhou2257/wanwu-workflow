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

import classNames from 'classnames';
import {
  ValueExpressionType,
  ViewVariableType,
} from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';

import AutoSizeTooltip from '@/ui-components/auto-size-tooltip';
import { DataTypeTag } from '@/node-registries/common/components';
import { ColumnsTitle } from '@/form-extensions/components/columns-title';
import { Section } from '@/form';

import { ValueExpressionInputField } from '../../../common/fields';

interface DatasetParamsFieldProps {
  required?: boolean;
}

export const DatasetParamsFieldWanwu = ({
  required = false,
}: DatasetParamsFieldProps) => {
  return (
    <Section
      title={I18n.t('workflow_detail_node_input')}
      tooltip={I18n.t(
        'workflow_detail_dataset_input_tooltip',
        {},
        '输入需要从知识中匹配的关键信息，Query和Image不能同时为空',
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
      <div className="w-full flex gap-[4px] items-center">
        <div style={{flex: 2}}>
          <div className="flex items-center max-w-[120px]">
            <AutoSizeTooltip
              content={'Query'}
              showArrow
              position="top"
              className="grow-1 truncate"
            >
            <span
              className={classNames(
                'flex-1 grow-1 truncate coz-fg-primary text-[12px] leading-[24px] mr-[2px]',
              )}
            >
              Query
            </span>
            </AutoSizeTooltip>
            {required ? (
              <span
                className="mt-[2px]"
                style={{
                  color: 'var(--light-usage-danger-color-danger,#f93920)',
                }}
              >
              *
            </span>
            ) : null}
            <DataTypeTag type={ViewVariableType.String}/>
          </div>
        </div>
        <div style={{flex: 3}}>
          <ValueExpressionInputField
            testId="/inputs/inputParameters/Query"
            name="inputs.inputParameters.Query"
            required={required}
            inputType={ViewVariableType.String}
            defaultValue={{ type: ValueExpressionType.REF }}
          />
        </div>
      </div>
      <div className="w-full flex gap-[4px] items-center mt-[10px]">
        <div style={{flex: 2}}>
          <div className="flex items-center max-w-[120px]">
            <AutoSizeTooltip
              content={'Image'}
              showArrow
              position="top"
              className="grow-1 truncate"
            >
            <span
              className={classNames(
                'flex-1 grow-1 truncate coz-fg-primary text-[12px] leading-[24px] mr-[2px]',
              )}
            >
              Image
            </span>
            </AutoSizeTooltip>
            <DataTypeTag type={ViewVariableType.Image} />
          </div>
        </div>
        <div style={{flex: 3}}>
          <ValueExpressionInputField
            testId="/inputs/inputParameters/Image"
            name="inputs.inputParameters.Image"
            required={required}
            inputType={ViewVariableType.Image}
            defaultValue={{ type: ValueExpressionType.REF }}
          />
        </div>
      </div>
    </Section>
  );
}
