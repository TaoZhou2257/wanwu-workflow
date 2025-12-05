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

import React from 'react';
import { I18n } from '@coze-arch/i18n';
import { IconCozPlayFill } from '@coze-arch/coze-design/icons';
import { Button } from '@coze-arch/coze-design';

import styles from './index-wanwu.module.less';

interface TestFormSheetFooterV2Props {
  onClick?: (e: React.MouseEvent) => void;
}

export const TestFormSheetFooterV2Wanwu: React.FC<TestFormSheetFooterV2Props> = ({
  onClick,
}) => (
  <div className={styles['test-form-sheet-footer-v2']}>
    <Button
      icon={<IconCozPlayFill />}
      color="green"
      onClick={onClick}
    >
      {I18n.t('workflow_debug_run')}
    </Button>
  </div>
);
