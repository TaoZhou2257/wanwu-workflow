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

import React, { useState, useEffect } from 'react';

import { I18n } from '@coze-arch/i18n';
import { Input, Button } from '@coze-arch/coze-design';
import { Section } from '@/form';


interface ApiKeyParamsFieldProps {
  defaultValue?: string;
  onChangeValue?: (value: any) => void;
}

export const ApiKeyField = ({
  defaultValue,
  onChangeValue,
}: ApiKeyParamsFieldProps) => {
  const [value, setValue] = useState<any>('')
  useEffect(() => {
    setValue(defaultValue)
  }, [])

  return (
    <Section title={'API KEY'}>
      <div className="w-full flex gap-[4px] items-center">
        <div style={{ flex: 8 }}>
          <Input value={value} onChange={(v) => setValue(v)} />
        </div>
        <div>
          <Button
            size="default"
            color="highlight"
            onClick={() => {
              onChangeValue?.(value)
            }}
          >
            {defaultValue ? I18n.t('Update') : I18n.t('confirm')}
          </Button>
        </div>
      </div>
    </Section>
  )
};
