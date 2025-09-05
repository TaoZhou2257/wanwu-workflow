import {
  ValidateTrigger,
  type FormMetaV2,
} from '@flowgram-adapter/free-layout-editor';
import { I18n } from '@coze-arch/i18n';
import { get } from "lodash-es";

import {
  fireNodeTitleChange,
  provideNodeOutputVariablesEffect,
} from '@/node-registries/common/effects';

import { type FormData } from './types';
import { FormRender } from './form';
import { transformOnInit, transformOnSubmit } from './data-transformer';
import { DEFAULT_PARAMS_LIST } from './constants'

export const GUI_AGENT_WANWU_FORM_META: FormMetaV2<FormData> = {
  // 节点表单渲染
  render: () => <FormRender />,

  // 验证触发时机
  validateTrigger: ValidateTrigger.onChange,

  // 验证规则
  validate: {
    // 必填
    'inputs.inputParameters.*.input': ({ value, formValues, name }) => {
      const currentKey = name.slice(0, name.lastIndexOf(".input")) || ''
      const currentName = get(formValues, currentKey)?.name || ''
      const { required } = DEFAULT_PARAMS_LIST.find(
        item => item.name === currentName
      ) || {}
      return required && !value.content
        ? I18n.t('workflow_detail_node_error_empty', {}, '参数值不可为空')
        : undefined
    }
  },

  // 副作用管理
  effect: {
    nodeMeta: fireNodeTitleChange,
    outputs: provideNodeOutputVariablesEffect,
  },

  // 节点后端数据 -> 前端表单数据
  formatOnInit: transformOnInit,

  // 前端表单数据 -> 节点后端数据
  formatOnSubmit: transformOnSubmit,
};
