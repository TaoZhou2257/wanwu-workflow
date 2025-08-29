import { nanoid } from 'nanoid';
import { ViewVariableType } from '@coze-workflow/variable';

// 入参路径，试运行等功能依赖该路径提取参数
export const INPUT_PATH = 'inputs.inputParameters';

// 定义固定出参
export const OUTPUTS = [
  {
    key: nanoid(),
    name: 'fileUrl',
    type: ViewVariableType.String,
  },
];

export const DEFAULT_PARAMS_LIST = [
  { name: `Title`, input: { type: 'ref' } },
  { name: `Content`, input: { type: 'ref' } },
  { name: `fileType`, input: { type: 'ref' } }
]
