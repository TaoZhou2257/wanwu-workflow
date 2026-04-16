import { type InputValueVO } from '@coze-workflow/base';
import { type IModelValue } from '@/typing';
import { type ToolSelectValue } from '@/form-extensions/setters/tool-select-wanwu';

export interface FormData {
  nodeMeta?: any;
  inputs: {
    inputParameters: InputValueVO[];
    toolInfoList?: ToolSelectValue;
    model?: IModelValue;
    systemPrompt: string;
    datasetSetting;
  };
  outputs?: any[];
}
