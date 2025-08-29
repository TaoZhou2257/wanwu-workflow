import {
  DEFAULT_NODE_META_PATH,
  DEFAULT_OUTPUTS_PATH,
} from '@coze-workflow/nodes';
import {
  StandardNodeType,
  type WorkflowNodeRegistry,
} from '@coze-workflow/base';

import { JSON_PARSER_WANWU_FORM_META } from './form-meta';
import { INPUT_PATH } from './constants';
import { test, type NodeTestMeta } from './node-test';

export const JSON_PARSER_WANWU_NODE_REGISTRY: WorkflowNodeRegistry<NodeTestMeta> = {
  type: StandardNodeType.JsonParserWanwu,
  meta: {
    nodeDTOType: StandardNodeType.JsonParserWanwu,
    size: { width: 360, height: 130.7 },
    nodeMetaPath: DEFAULT_NODE_META_PATH,
    outputsPath: DEFAULT_OUTPUTS_PATH,
    inputParametersPath: INPUT_PATH,
    test,
  },
  formMeta: JSON_PARSER_WANWU_FORM_META,
};
