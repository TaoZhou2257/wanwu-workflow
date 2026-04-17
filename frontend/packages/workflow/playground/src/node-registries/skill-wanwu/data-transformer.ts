import { get, camelCase } from 'lodash-es';
import { type NodeFormContext } from '@flowgram-adapter/free-layout-editor';
import { type NodeDataDTO, type InputValueDTO, BlockInput } from '@coze-workflow/base';
import { ModelParamType } from '@coze-arch/bot-api/developer_api';

import { type FormData } from './types';
import { OUTPUTS } from './constants';
import { WorkflowModelsService } from '@/services';
import { getDefaultLLMParams } from '@/nodes-v2/llm-wanwu/utils';
import { type IModelValue } from '@/typing';
import { normalizeInputParameters, parseLLMParam } from './util';

/**
 * 节点后端数据 -> 前端表单数据
 */
export function transformOnInit (value: FormData,context: NodeFormContext) {
  const { node } = context;
  //模型参数
  const modelsService = node.getService<WorkflowModelsService>(
    WorkflowModelsService,
  );
  const models = modelsService.getModels();
  let llmParam = get(value, 'inputs.llmParam') as InputValueDTO[] | undefined;
  if (!llmParam) {
    llmParam = getDefaultLLMParams(models);
  }
  const normalizedLLMParam = normalizeInputParameters(llmParam);
  const modelParams = parseLLMParam(normalizedLLMParam);
  const { systemPrompt, ...model } = modelParams;
  const modelValue = model.modelType ? { ...model } : undefined;

  //输入参数
  const inputParameters = get(value, 'inputs.inputParameters', []) as InputValueDTO[];

  //skill
  const agentSkillParams = (() => {
    const rawData = get(value, 'inputs.agentSkillParams', []);
    if (!Array.isArray(rawData)) return [];
    return rawData.map(item => {
      if (typeof item !== 'object' || item === null) {
        return { kind: 'skill', id: '', name: '', description: '' };
      }
      return {
        ...(item as any),
        kind: 'skill',
        id: String((item as any).skillId || ''),
        name: String((item as any).skillName || ''),
        description: String((item as any).desc || '')
      };
    });
  })() as any[];

  const toolInfoList = [...agentSkillParams];
  const outputs = value?.outputs ?? OUTPUTS;
  return {
    nodeMeta: value?.nodeMeta,
    inputs: {
      inputParameters,
      systemPrompt,
      llmParam: modelValue,
      toolInfoList,
    },
    outputs,
  };
};

/**
 * 前端表单数据 -> 节点后端数据
 */
export function transformOnSubmit (value,context: NodeFormContext){
  const { node, playgroundContext } = context;
  const toolInfoList = get(value, 'inputs.toolInfoList', []) as any[];
  const agentSkillParams: any[] = [];

  toolInfoList.forEach(tool => {
    if (tool.kind === 'skill') {
      agentSkillParams.push({
        ...tool,
        skillId: tool.id,
        skillName: tool.name,
        description: tool.description || tool.desc,
      });
    }
  });

  const llmParam: any[] = [];
  const modelsService = node.getService<WorkflowModelsService>(
    WorkflowModelsService,
  );
  const models = modelsService?.getModels() ?? [];
  const model = get(value, 'inputs.llmParam') as IModelValue | undefined;
  const modelType = model?.modelType;
  const modelMeta = modelType
    ? models.find(m => m.model_type === modelType)
    : undefined;
  
  if (model?.modelType) {
    llmParam.push(
      BlockInput.createInteger('modelType', String(model.modelType)),
    );
  }

  if (model) {
    const excludeKeys = ['modelType', 'modelName', 'generationDiversity', 'responseFormat'];
    Object.keys(model).forEach(k => {
      if (excludeKeys.includes(k)) {
        return;
      }
      const paramValue = model[k];
      if (paramValue === undefined || paramValue === null) {
        return;
      }
      
      const paramDef = modelMeta?.model_params?.find(
        p => camelCase(p.name) === k,
      );
      const paramType = paramDef?.type;
      
      if (ModelParamType.Float === paramType) {
        llmParam.push(BlockInput.createFloat(k, String(paramValue)));
      } else if (ModelParamType.Int === paramType || ['modelType'].includes(k)) {
        llmParam.push(BlockInput.createInteger(k, String(paramValue)));
      } else {
        let _k = k;
        if (_k === 'modelName') {
          _k = 'modleName';
        }
        llmParam.push(BlockInput.createString(_k, String(paramValue)));
      }
    });
  }

  const systemPrompt = get(value, 'inputs.systemPrompt') as string | undefined;
  if (systemPrompt) {
    llmParam.push(
      BlockInput.createString('systemPrompt', String(systemPrompt)),
    );
  }

  // 转换输入参数
  const inputParameters = get(value, 'inputs.inputParameters', []);
  const safeInputParameters = Array.isArray(inputParameters) ? inputParameters : [];
 
  return {
    nodeMeta: value?.nodeMeta,
    inputs: {
      inputParameters: safeInputParameters,
      llmParam,
      agentSkillParams,
    },
    outputs: value.outputs,
  } as NodeDataDTO;
};
