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

import { useNodeTestId, type WorkflowDatabase } from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';

import {
  McpSelect,
  type McpSelectValue,
} from '@/form-extensions/setters/mcp-select-wanwu';
import { Section, Field, type FieldProps } from '@/form';

type McpSelectFieldProps = FieldProps<McpSelectValue> & {
  afterChange?: (value?: WorkflowDatabase) => void;
};

export const McpSelectField = ({
  name,
  label,
  tooltip,
  afterChange,
  ...rest
}: McpSelectFieldProps) => {
  const { getNodeSetterId } = useNodeTestId();

  return (
    <Section title={I18n.t('mcp_model_title')}>
      <Field<McpSelectValue> name={name} {...rest}>
        {({ value, onChange, readonly }) => (
          <McpSelect
            value={value}
            readonly={readonly}
            onChange={newValue => {
              onChange(newValue);
              afterChange?.(newValue);
            }}
            addButtonTestID={getNodeSetterId(`${name}.addButton`)}
            libraryCardTestID={getNodeSetterId(`${name}.libraryCard`)}
          />
        )}
      </Field>
    </Section>
  );
};
