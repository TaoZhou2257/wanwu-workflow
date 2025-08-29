import {
  DEFAULT_NODE_META_PATH,
  DEFAULT_OUTPUTS_PATH,
} from '@coze-workflow/nodes';
import {
  StandardNodeType,
  type WorkflowNodeRegistry,
} from '@coze-workflow/base';

import { FILE_PARSE_WANWU_FORM_META } from './form-meta';
import { INPUT_PATH } from './constants';
import { test, type NodeTestMeta } from './node-test';

export const FILE_PARSE_WANWU_NODE_REGISTRY: WorkflowNodeRegistry<NodeTestMeta> = {
  type: StandardNodeType.FileParseWanwu,
  meta: {
    nodeDTOType: StandardNodeType.FileParseWanwu,
    size: { width: 360, height: 130.7 },
    nodeMetaPath: DEFAULT_NODE_META_PATH,
    outputsPath: DEFAULT_OUTPUTS_PATH,
    inputParametersPath: INPUT_PATH,
    test,
  },
  formMeta: FILE_PARSE_WANWU_FORM_META,
};
