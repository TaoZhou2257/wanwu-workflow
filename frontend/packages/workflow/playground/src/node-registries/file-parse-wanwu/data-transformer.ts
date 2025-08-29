import { OUTPUTS } from './constants';
import { set } from "lodash-es";

/**
 * 节点后端数据 -> 前端表单数据
 */
export function transformOnInit(value) {
  // New drag-in node initialization
  if (!value) {
    return {
      nodeMeta: undefined,
      inputs: {
        inputParameters: {
          FileUrl: { type: 'ref', content: '' },
        },
      },
      outputs: OUTPUTS,
    };
  }

  const { inputParameters } = value.inputs;
  const formData = {
    ...value,
    inputs: {},
  };

  formData.inputs.inputParameters = inputParameters.reduce(
    (map, obj: { name: string | number; input: unknown }) => {
      map[obj.name] = obj.input;
      return map;
    },
    {},
  );

  return formData;
}

/**
 * 前端表单数据 -> 节点后端数据
 * @param value
 * @returns
 */
export function transformOnSubmit(value) {
  const { nodeMeta, inputs, outputs } = value;
  const { inputParameters = { FileUrl: { type: 'ref' } } } =
  inputs ?? {};
  const actualData = {
    nodeMeta,
    outputs,
    inputs: {},
  };
  set(
    actualData.inputs,
    'inputParameters',
    Object.entries(inputParameters).map(([key, mapValue]) => ({
      name: key,
      input: mapValue,
    })) || [],
  );

  return actualData;
}
