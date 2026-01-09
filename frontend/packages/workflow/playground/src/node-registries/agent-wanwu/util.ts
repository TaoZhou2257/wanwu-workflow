import { type InputValueDTO, VariableTypeDTO } from '@coze-workflow/base';

/**
 * 
 * @param inputParameters 输入参数数组
 * @returns 标准化后的输入参数数组
 */
export function normalizeInputParameters(
  inputParameters: InputValueDTO[],
): InputValueDTO[] {
  return inputParameters.map(input => {
    if (!input.input) {
      return input;
    }
    const inputValue = input.input as any;
    if (inputValue.type && 'content' in inputValue && !('value' in inputValue)) {
      return {
        ...input,
        input: {
          type: inputValue.type as string,
          value: {
            type: inputValue.type,
            content: inputValue.content,
            rawMeta: inputValue.rawMeta,
          },
        } as any,
      };
    }
    return input;
  });
}

export function  transformRefContent (content: any): any  {
  if (!content) return content;
  if (content.source === 'block-output' && content.blockID && content.name) {
    return {
      keyPath: [content.blockID, content.name]
    };
  }
  if (content.keyPath) return content;
  return content;
};

/*
 * 
 * @param llmParam LLM 参数数组
 * @returns 参数对象，键为参数名，值为参数值
 */
export function parseLLMParam(
  llmParam: InputValueDTO[] = [],
): Record<string, unknown> {
  const model: Record<string, unknown> = {};

  llmParam.forEach((d: InputValueDTO) => {
    let name = d?.name || '';
    if (name === 'modleName') {
      name = 'modelName';
    }
    
    if (!d.input?.value) {
      return;
    }
    
    let value = d.input.value.content;
    
    if (
      [VariableTypeDTO.float, VariableTypeDTO.integer].includes(
        d.input.type as VariableTypeDTO,
      )
    ) {
      value = Number(value);
    }

    model[name] = value;
  });

  return model;
}

/**
 * @param knowledgeList 
 * @returns 
 */
export function formatKnowledgeList<T extends { name?: string; dataset_id?: string; id?: string; description?: string }>(
  knowledgeList: (string | T)[],
): (T & { name: string; dataset_id: string; kind: 'database'; id: string; description: string })[] {
  return knowledgeList.map((item, index) => {
    if (typeof item === 'string') {
      return {
        name: item,
        dataset_id: item,
        kind: 'database',
        id: item || `knowledge-${index}`,
        description: ''
      } as T & { name: string; dataset_id: string; kind: 'database'; id: string; description: string };
    } else {
      return {
        ...item,
        kind: 'database',
        id: item.dataset_id || `knowledge-${index}`,
        description: item.description || ''
      } as T & { name: string; dataset_id: string; kind: 'database'; id: string; description: string };
    }
  });
}

/**
 * 将输入参数从嵌套 value 格式转换为扁平格式
 * 转换前: { name: "query", input: { type: "ref", value: { type: "ref", content: { keyPath: [...] }, rawMeta: {...} } } }
 * 转换后: { name: "query", input: { type: "ref", content: { keyPath: [""] }, rawMeta: {...} } }
 * 
 * @param inputItem 输入参数项
 * @returns 转换后的输入参数项
 */
export function transformInputValueToFlatFormat(inputItem: any): any {
  if (!inputItem || !inputItem.input) {
    return inputItem;
  }
  const { input } = inputItem;
  if (input.value && typeof input.value === 'object') {
    const { value } = input;
    
    return {
      ...inputItem,
      input: {
        type: value.type || input.type,
        content: {
          keyPath: ['']
        },
        rawMeta: value.rawMeta || input.rawMeta
      }
    };
  }
  
  if (input.content || input.type) {
    return {
      ...inputItem,
      input: {
        ...input,
        content: {
          keyPath: ['']
        }
      }
    };
  }
  return inputItem;
}

