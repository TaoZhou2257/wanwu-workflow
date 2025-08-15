import {
  DEFAULT_NODE_META_PATH,
  DEFAULT_OUTPUTS_PATH,
} from '@coze-workflow/nodes';
import {
  StandardNodeType,
  type WorkflowNodeRegistry,
} from '@coze-workflow/base';

import { MCP_WANWU_FORM_META } from './form-meta';
import { test, type NodeTestMeta } from './node-test';

export const MCP_WANWU_NODE_REGISTRY: WorkflowNodeRegistry<NodeTestMeta> = {
  type: StandardNodeType.McpWanwu,
  meta: {
    nodeDTOType: StandardNodeType.McpWanwu,
    size: { width: 360, height: 130.7 },
    nodeMetaPath: DEFAULT_NODE_META_PATH,
    outputsPath: DEFAULT_OUTPUTS_PATH,
    inputParametersPath: '/inputParameters', //INPUT_PATH,
    test,
  },
  formMeta: MCP_WANWU_FORM_META,
};
