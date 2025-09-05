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

import React, { useEffect, useState } from 'react';

import { useNodeTestId } from '@coze-workflow/base';
import { Select } from '@coze-arch/coze-design';
import { KnowledgeApi } from '@coze-arch/bot-api';
import { I18n } from '@coze-arch/i18n';

interface GuiModelProps {
  value?: string;
  onChange?: (v: string) => void;
  style?: React.CSSProperties;
  readonly?: boolean;
}

export const GuiModelWanwu: React.FC<GuiModelProps> = props => {
  const { value, onChange, style, readonly } = props;

  const { getNodeSetterId } = useNodeTestId();
  const [guiModelList, setGuiModelList] = useState<any>([])

  useEffect(() => {
    getGuiModel()
  }, []);

  const getGuiModel = async () => {
    const { data = {} } = await KnowledgeApi.getGuiModel()
    setGuiModelList(data?.list || [])
  }

  return (
    <Select
      className="w-[100%]"
      size="small"
      value={value}
      style={{
        ...style,
        pointerEvents: readonly ? 'none' : 'auto',
      }}
      placeholder={I18n.t('model_placeholder')}
      onChange={onChange as (v: unknown) => void}
      data-testid={getNodeSetterId('gui-model')}
    >
      {guiModelList.map(v => (
        <Select.Option
          value={v.modelId}
          key={v.modelId}
          data-testid={getNodeSetterId('gui-model-option')}
        >
          {v.displayName || v.model}
        </Select.Option>
      ))}
    </Select>
  );
};
