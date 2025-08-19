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
import { Select } from '@coze-arch/coze-design';

import s from './index.module.less';

const optionList = [
  {
    "modelId": "31",
    "provider": "YuanJing",
    "model": "bge1",
    "modelType": "rerank",
    "displayName": "bge1",
    "avatar": {
      "key": "",
      "path": ""
    },
    "publishDate": "",
    "isActive": true,
    "userId": "1",
    "orgId": "1",
    "createdAt": "",
    "updatedAt": ""
  },
  {
    "modelId": "69",
    "provider": "YuanJing",
    "model": "yuanjingrerank",
    "modelType": "rerank",
    "displayName": "yuanjingrerank",
    "avatar": {
      "key": "",
      "path": ""
    },
    "publishDate": "",
    "isActive": true,
    "userId": "1",
    "orgId": "1",
    "createdAt": "",
    "updatedAt": ""
  },
  {
    "modelId": "67",
    "provider": "Qwen",
    "model": "gte-rerank",
    "modelType": "rerank",
    "displayName": "gte-rerank",
    "avatar": {
      "key": "",
      "path": ""
    },
    "publishDate": "",
    "isActive": true,
    "userId": "1",
    "orgId": "1",
    "createdAt": "",
    "updatedAt": ""
  },
  {
    "modelId": "10",
    "provider": "OpenAI-API-compatible",
    "model": "jina-reranker-m0",
    "modelType": "rerank",
    "displayName": "jina-reranker-m0",
    "avatar": {
      "key": "custom-upload/avatar/42/425d6355-ec08-4f8e-84ef-c921febbb1b7.png",
      "path": "/v1/cache/avatar/42/425d6355-ec08-4f8e-84ef-c921febbb1b7.png"
    },
    "publishDate": "2025-06-25",
    "isActive": true,
    "userId": "1",
    "orgId": "1",
    "createdAt": "",
    "updatedAt": ""
  },
  {
    "modelId": "11",
    "provider": "OpenAI-API-compatible",
    "model": "BAAI/bge-reranker-v2-m3",
    "modelType": "rerank",
    "displayName": "bge-reranker-v2-m3",
    "avatar": {
      "key": "custom-upload/avatar/d0/d04d8ea3-4133-4b4c-842a-f01ff957f36d.png",
      "path": "/v1/cache/avatar/d0/d04d8ea3-4133-4b4c-842a-f01ff957f36d.png"
    },
    "publishDate": "",
    "isActive": true,
    "userId": "1",
    "orgId": "1",
    "createdAt": "",
    "updatedAt": ""
  }
];

interface RerankModelProps {
  value: string;
  onChange: (v: string) => void;
  style?: React.CSSProperties;
  readonly?: boolean;
}

export const RerankModelWanwu: React.FC<RerankModelProps> = props => {
  const { value, onChange, style, readonly } = props;

  const { getNodeSetterId } = useNodeTestId();

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
      data-testid={getNodeSetterId('dataset-rerank-model')}
    >
      {optionList.map(v => (
        <Select.Option
          value={v.modelId}
          key={v.modelId}
          data-testid={getNodeSetterId('dataset-rerank-model-option')}
        >
          {v.displayName || v.model}
        </Select.Option>
      ))}
    </Select>
  );
};
