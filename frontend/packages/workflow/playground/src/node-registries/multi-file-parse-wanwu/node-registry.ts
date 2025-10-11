import {
  DEFAULT_NODE_META_PATH,
  DEFAULT_OUTPUTS_PATH,
} from '@coze-workflow/nodes';
import {
  StandardNodeType,
  type WorkflowNodeRegistry,
} from '@coze-workflow/base';

import { MULTI_FILE_PARSE_WANWU_FORM_META } from './form-meta';
import { INPUT_PATH } from './constants';
import { test, type NodeTestMeta } from './node-test';

export const MULTI_FILE_PARSE_WANWU_NODE_REGISTRY: WorkflowNodeRegistry<NodeTestMeta> = {
  type: StandardNodeType.MultiFileParseWanwu,
  meta: {
    nodeDTOType: StandardNodeType.MultiFileParseWanwu,
    size: { width: 360, height: 130.7 },
    nodeMetaPath: DEFAULT_NODE_META_PATH,
    outputsPath: DEFAULT_OUTPUTS_PATH,
    inputParametersPath: INPUT_PATH,
    test,
  },
  formMeta: MULTI_FILE_PARSE_WANWU_FORM_META,
};
