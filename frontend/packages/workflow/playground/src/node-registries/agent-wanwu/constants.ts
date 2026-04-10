import { nanoid } from 'nanoid';
import { ViewVariableType } from '@coze-workflow/variable';
import {
  ValueExpressionType,
  type InputValueVO,
} from '@coze-workflow/base';

// 入参路径，试运行等功能依赖该路径提取参数
export const INPUT_PATH = 'inputs.inputParameters';

// 定义固定出参
export const OUTPUTS = [
  {
    key: nanoid(),
    name: 'output',
    type: ViewVariableType.Object,
    children: [
      { key: nanoid(), name: 'fullResponse', type: ViewVariableType.Object },
      { key: nanoid(), name: 'response', type: ViewVariableType.String },
      { key: nanoid(), name: 'searchList', type: ViewVariableType.String }
    ]
  },
];

export const FILE_INPUT_AVAILABLE_TYPES = [
  ViewVariableType.File,
  ViewVariableType.Image,
  ViewVariableType.Doc,
  ViewVariableType.Ppt,
  ViewVariableType.Excel,
  ViewVariableType.Txt, 
] as const;

export const DEFAULT_PARAMS_LIST = [
  { name: `query`, input: { type: 'ref', rawMeta: { type: ViewVariableType.String } }, required: true, type: ViewVariableType.String },
  { name: `file`, input: { type: 'ref',  rawMeta: { type: ViewVariableType.File } }, required: false, type: ViewVariableType.File }
]
