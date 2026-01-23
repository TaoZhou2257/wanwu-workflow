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

import { useNodeTestId } from '@coze-workflow/base';
import { I18n } from '@coze-arch/i18n';
import { Select } from '@coze-arch/coze-design';

import { MatchType } from '../../type';

import s from './index.module.less';

const optionList = [
  {
    value: MatchType.Semantic,
    label: I18n.t('knowledge_vector_search'),
  },
  {
    value: MatchType.FullText,
    label: I18n.t('knowledge_full_text_search'),
  },
  {
    value: MatchType.Hybird,
    label: I18n.t('knowledge_hybrid_search'),
  },
];

const otherOptionList = [
  {
    value: MatchType.HybirdPriority,
    label: I18n.t('knowledge_hybrid_priority_search'),
  },
]

interface SearchStrategyProps {
  value: MatchType;
  onChange: (v: MatchType) => void;
  style?: React.CSSProperties;
  readonly?: boolean;
  isQaKnowledge?: boolean;
}

export const SearchStrategyWanwu: React.FC<SearchStrategyProps> = props => {
  const { value, onChange, style, readonly, isQaKnowledge } = props;

  const { getNodeSetterId } = useNodeTestId();
  const options = isQaKnowledge ? optionList : [...optionList, ...otherOptionList]

  return (
    <Select
      className={s['strategy-area']}
      dropdownClassName={s['strategy-area-dropdown']}
      size="small"
      value={value}
      style={{
        ...style,
        pointerEvents: readonly ? 'none' : 'auto',
      }}
      onChange={onChange as (v: unknown) => void}
      // defaultValue={MatchType.Semantic}
      data-testid={getNodeSetterId('dataset-search-strategy')}
    >
      {options.map(v => (
        <Select.Option
          value={v.value}
          key={v.value}
          data-testid={getNodeSetterId('dataset-search-strategy-option')}
        >
          {v.label}
        </Select.Option>
      ))}
    </Select>
  );
};
