/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { useState, useEffect } from 'react';

import { Modal, Input, Button } from '@coze-arch/coze-design';
import { I18n } from '@coze-arch/i18n';

interface ToolSelectPluginSettingProps {
  visible: boolean;
  value?: string;
  onChange?: (value: string) => void;
  onOk?: () => void;
  onCancel?: () => void;
  readonly?: boolean;
}

export const ToolSelectPluginSetting = ({
  visible,
  value = '',
  onChange,
  onOk,
  onCancel,
  readonly = false,
}: ToolSelectPluginSettingProps) => {
  const [inputValue, setInputValue] = useState<string>(value);

  useEffect(() => {
    setInputValue(value);
  }, [value, visible]);

  const handleUpdate = () => {
    if (onChange) {
      onChange(inputValue);
    }
  };

  const handleOk = () => {
    if (onChange) {
      onChange(inputValue);
    }
    if (onOk) {
      onOk();
    }
  };

  const handleCancel = () => {
    setInputValue(value);
    if (onCancel) {
      onCancel();
    }
  };

  return (
    <Modal
      title={I18n.t('workflow_agnet_tool_mcp_setting_apikey' as any, {}, '编辑API KEY')}
      visible={visible}
      onOk={handleOk}
      onCancel={handleCancel}
      width={480}
      okText={I18n.t('workflow_confirm_modal_ok', {}, '确定')}
      cancelText={I18n.t('workflow_confirm_modal_cancel', {}, '取消')}
    >
      <div className="flex flex-col gap-[16px] py-[8px]">
        <div className="flex flex-col gap-[8px]">
          <label className="text-[14px] leading-[20px] font-medium coz-fg-primary">
            {I18n.t('apikey', {}, 'API KEY')}
          </label>
          <div className="flex gap-[8px] items-center">
            <Input
              type="password"
              value={inputValue}
              onChange={setInputValue}
              placeholder={I18n.t('please_enter_ark_apikey', {}, '请输入 API KEY')}
              disabled={readonly}
              className="flex-1"
            />
          </div>
        </div>
      </div>
    </Modal>
  );
};


