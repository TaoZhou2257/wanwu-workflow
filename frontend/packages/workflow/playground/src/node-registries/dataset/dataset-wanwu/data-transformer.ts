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

import { nanoid } from 'nanoid';
import { isNil, set } from 'lodash-es';
import { BlockInput, ViewVariableType } from '@coze-workflow/base';

export function transformOnInit(value) {
  // New drag-in node initialization
  if (!value) {
    return {
      nodeMeta: undefined,
      inputs: {
        inputParameters: {
          Query: { type: 'ref', content: '' },
        },
        datasetParameters: {
          datasetParam: [],
          datasetSetting: {},
        },
      },
      outputs: [
        {
          key: nanoid(),
          name: 'outputList',
          type: ViewVariableType.ArrayObject,
          children: [
            {
              key: nanoid(),
              name: 'prompt',
              type: ViewVariableType.String,
            },
            {
              key: nanoid(),
              name: 'score',
              type: ViewVariableType.Number,
            },
            {
              key: nanoid(),
              name: 'searchList',
              type: ViewVariableType.ArrayString,
            },
          ],
        },
      ],
    };
  }

  const { inputParameters, datasetParam } = value.inputs;
  const formData = {
    ...value,
    inputs: {
      datasetParameters: {},
    },
  };

  formData.inputs.inputParameters = inputParameters.reduce(
    (map, obj: { name: string | number; input: unknown }) => {
      map[obj.name] = obj.input;
      return map;
    },
    {},
  );
  formData.inputs.datasetParameters.datasetParam = datasetParam[0]?.input.value
    .content as string[];
  // In the case of initial creation/stock data, the topK and threshold are empty, and the initial default value is processed in the dataset-settings component
  formData.inputs.datasetParameters.datasetSetting = {
    topK: datasetParam.find(item => item.name === 'topK')?.input.value
      .content as number,

    threshold: datasetParam.find(item => item.name === 'threshold')?.input.value
      .content as number,

    semanticsPriority: datasetParam.find(item => item.name === 'semanticsPriority')?.input.value
      .content as number,

    maxHistory: datasetParam.find(item => item.name === 'maxHistory')?.input.value
      .content as number,

    matchType: datasetParam.find(item => item.name === 'matchType')?.input.value
      .content as string,

    rerankModelId: datasetParam.find(item => item.name === 'rerankModelId')?.input.value
      .content as string,

    rewrite: datasetParam.find(item => item.name === 'rewrite')?.input
      .value.content as boolean,
  };

  return formData;
}

export function transformOnSubmit(value) {
  const { nodeMeta, inputs, outputs } = value;
  const { inputParameters = { Query: { type: 'ref' } }, datasetParameters } =
    inputs ?? {};
  const { datasetParam, datasetSetting } = datasetParameters ?? {};
  const actualData = {
    nodeMeta,
    outputs,
    inputs: {
      datasetParam: [] as unknown[],
    },
  };

  set(
    actualData.inputs,
    'inputParameters',
    Object.entries(inputParameters).map(([key, mapValue]) => ({
      name: key,
      input: mapValue,
    })) || [],
  );

  set(actualData.inputs, 'datasetParam', [
    {
      name: 'knowledgeList',
      input: {
        type: 'list',
        schema: {
          type: 'string',
        },
        value: {
          type: 'literal',
          content: datasetParam || [],
        },
      },
    },
    {
      name: 'topK',
      input: {
        type: 'integer',
        value: {
          type: 'literal',
          content: datasetSetting?.topK,
        },
      },
    },
    BlockInput.createBoolean('rewrite', datasetSetting?.rewrite),
  ]);

  if (datasetSetting?.threshold !== undefined) {
    actualData.inputs.datasetParam.push({
      name: 'threshold',
      input: {
        type: 'float',
        value: {
          type: 'literal',
          content: datasetSetting?.threshold,
        },
      },
    });
  }

  if (datasetSetting?.semanticsPriority !== undefined) {
    actualData.inputs.datasetParam.push({
      name: 'semanticsPriority',
      input: {
        type: 'float',
        value: {
          type: 'literal',
          content: datasetSetting?.semanticsPriority,
        },
      },
    });
  }

  if (datasetSetting?.maxHistory !== undefined) {
    actualData.inputs.datasetParam.push({
      name: 'maxHistory',
      input: {
        type: 'integer',
        value: {
          type: 'literal',
          content: datasetSetting?.maxHistory,
        },
      },
    });
  }

  // Added search policy configuration, there may be no strategy data not in grey release
  // Strategy may be 0
  if (!isNil(datasetSetting?.matchType)) {
    actualData.inputs.datasetParam.push({
      name: 'matchType',
      input: {
        type: 'string',
        value: {
          type: 'literal',
          content: datasetSetting?.matchType,
        },
      },
    });
  }

  if (!isNil(datasetSetting?.rerankModelId)) {
    actualData.inputs.datasetParam.push({
      name: 'rerankModelId',
      input: {
        type: 'string',
        value: {
          type: 'literal',
          content: datasetSetting?.rerankModelId,
        },
      },
    });
  }

  return actualData;
}
