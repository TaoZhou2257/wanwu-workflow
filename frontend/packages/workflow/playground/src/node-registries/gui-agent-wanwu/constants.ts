import { nanoid } from 'nanoid';
import { ViewVariableType } from '@coze-workflow/variable';

// 入参路径，试运行等功能依赖该路径提取参数
export const INPUT_PATH = 'inputs.inputParameters';

// 定义固定出参
export const OUTPUTS = [
  {
    key: nanoid(),
    name: 'code',
    type: ViewVariableType.Number,
  },
  {
    key: nanoid(),
    name: 'message',
    type: ViewVariableType.String,
  },
  {
    key: nanoid(),
    name: 'content',
    type: ViewVariableType.Object,
  },
  {
    key: nanoid(),
    name: 'usage',
    type: ViewVariableType.Object,
  },
];

export const DEFAULT_PARAMS_LIST = [
  { name: `platform`, input: { type: 'ref' }, required: true, type: ViewVariableType.String },
  { name: `current_screenshot`, input: { type: 'ref' }, required: true, type: ViewVariableType.String },
  { name: `current_screenshot_width`, input: { type: 'ref' }, required: true, type: ViewVariableType.Integer },
  { name: `current_screenshot_height`, input: { type: 'ref' }, required: true, type: ViewVariableType.Integer },
  { name: `task`, input: { type: 'ref' }, required: true, type: ViewVariableType.String },
  { name: `history`, input: { type: 'ref' }, type: ViewVariableType.ArrayString },
]
