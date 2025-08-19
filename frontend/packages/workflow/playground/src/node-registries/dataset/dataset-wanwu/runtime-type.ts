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

// eslint-disable-next-line @coze-arch/no-batch-import-or-export
import * as t from 'io-ts';

interface OutputType {
  key: string;
  name: string;
  type: number;
  children?: OutputType[];
}

const nodeMetaType = t.intersection([
  t.type({
    title: t.string,
  }),
  t.partial({
    icon: t.string,
    subtitle: t.string,
    description: t.string,
  }),
]);

const output: t.Type<OutputType> = t.recursion('output', () =>
  t.intersection([
    t.type({
      key: t.string,
      name: t.string,
      type: t.number,
    }),
    t.partial({
      children: t.array(output),
    }),
  ]),
);

const outputsType = t.array(output);

const queryType = t.union([
  t.type({
    type: t.literal('ref'),
    content: t.type({
      keyPath: t.array(t.string),
    }),
  }),
  t.type({
    type: t.literal('literal'),
    content: t.string,
  }),
]);

export const datasetNodeFormDataRuntimeType = t.type({
  nodeMeta: nodeMetaType,
  inputs: t.type({
    datasetParameters: t.type({
      datasetParam: t.array(t.string),
      datasetSetting: t.type({
        topK: t.number,
        threshold: t.union([t.number, t.undefined]),
        semanticsPriority: t.union([t.number, t.undefined]),
        maxHistory: t.union([t.number, t.undefined]),
        matchType: t.string,
        rerankModelId: t.string,
        rewrite: t.boolean,
      }),
    }),
    inputParameters: t.type({
      queryType,
    }),
  }),
  outputs: t.array(output),
});

export const datasetNodeActualDataRuntimeType = t.type({
  nodeMeta: nodeMetaType,
  inputs: t.type({
    datasetParam: t.array(
      t.union([
        t.type({
          name: t.literal('knowledgeList'),
          input: t.type({
            type: t.literal('list'),
            schema: t.type({
              type: t.literal('string'),
            }),
            value: t.type({
              type: t.literal('literal'),
              content: t.array(t.string),
            }),
          }),
        }),
        t.type({
          name: t.literal('topK'),
          input: t.type({
            type: t.literal('integer'),
            value: t.type({
              type: t.literal('literal'),
              content: t.number,
            }),
          }),
        }),
        t.type({
          name: t.literal('threshold'),
          input: t.type({
            type: t.literal('number'),
            value: t.type({
              type: t.literal('literal'),
              content: t.number,
            }),
          }),
        }),
        t.type({
          name: t.literal('semanticsPriority'),
          input: t.type({
            type: t.literal('number'),
            value: t.type({
              type: t.literal('literal'),
              content: t.number,
            }),
          }),
        }),
        t.type({
          name: t.literal('maxHistory'),
          input: t.type({
            type: t.literal('integer'),
            value: t.type({
              type: t.literal('literal'),
              content: t.number,
            }),
          }),
        }),
        t.type({
          name: t.literal('matchType'),
          input: t.type({
            type: t.literal('string'),
            value: t.type({
              type: t.literal('literal'),
              content: t.number,
            }),
          }),
        }),
        t.type({
          name: t.literal('rerankModelId'),
          input: t.type({
            type: t.literal('string'),
            value: t.type({
              type: t.literal('literal'),
              content: t.number,
            }),
          }),
        }),
        t.type({
          name: t.literal('rewrite'),
          input: t.type({
            type: t.literal('boolean'),
            value: t.type({
              type: t.literal('literal'),
              content: t.boolean,
            }),
          }),
        }),
      ]),
    ),
    datasetSetting: t.type({
      topK: t.number,
      threshold: t.number,
      semanticsPriority: t.number,
      maxHistory: t.number
    }),
    inputParameters: t.array(
      t.type({
        name: t.literal('Query'),
        input: queryType,
      }),
    ),
  }),
  outputs: outputsType,
});
