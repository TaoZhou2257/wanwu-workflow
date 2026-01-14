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

import { useNodeTestId } from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';

import {
  ToolSelect,
  type ToolSelectValue,
} from '@/form-extensions/setters/tool-select-wanwu';
import { Section, Field, type FieldProps, useForm } from '@/form';

type ToolSelectFieldProps = FieldProps<ToolSelectValue> & {
  afterChange?: (value?: ToolSelectValue) => void;
  showDataset?: boolean | undefined;
};

export const ToolSelectField = ({
  name,
  label,
  tooltip,
  afterChange,
  showDataset,
  ...rest
}: ToolSelectFieldProps) => {
  const { getNodeSetterId } = useNodeTestId();
  const form = useForm();
  return (
    <Section
      title={
        showDataset
          ? I18n.t('workflow_detail_knowledge_knowledge')
          : I18n.t('workflow_agnet_tool_select_title' as any)
      }
    >
      <Field<ToolSelectValue> name={name} {...rest}>
        {({ value, onChange, readonly}) => (
          <ToolSelect 
            value={value}
            readonly={readonly}
            onChange={newValue => {
              onChange(newValue);
              afterChange?.(newValue);
            }}
            addButtonTestID={getNodeSetterId(`${name}.addButton`)}
            libraryCardTestID={getNodeSetterId(`${name}.libraryCard`)}
            form={form}
            showDataset={showDataset}
          />
        )}
      </Field>
    </Section>
  );
};
