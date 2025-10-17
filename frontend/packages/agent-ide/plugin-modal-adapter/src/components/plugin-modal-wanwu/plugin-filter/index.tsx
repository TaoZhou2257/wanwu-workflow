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

import classNames from 'classnames';
import { I18n } from '@coze-arch/i18n';
import { type Int64 } from '@coze-arch/bot-api/developer_api';

import s from './plugin-filter.module.less';

export interface PluginFilterProps {
  isSearching: boolean;
  type: Int64;
  onChange: (type: Int64) => void;
  isShowStorePlugin?: boolean;
}

export const PluginFilter: React.FC<PluginFilterProps> = ({
  isSearching,
  type,
  onChange,
  isShowStorePlugin = true,
}) => {
  const onChangeAfterDiff = (freshType: typeof type) => {
    // If you are searching, leave the search blank
    if (isSearching) {
      onChange(freshType);
      return;
    }
    if (freshType === type) {
      return;
    }
    onChange(freshType);
  };

  return (
    <div className={s['tool-tag-list']}>
      {isShowStorePlugin ? (
        <>
          <div className={s['tool-content-area']}>
            <div
              className={classNames(s['tool-tag-list-cell'], {
                [s.active]: type === 'builtin',
              })}
              onClick={() => onChangeAfterDiff('builtin')}
            >
              {I18n.t('builtin_tools')}
            </div>
          </div>
          <div className={s['tool-content-area']}>
            <div
              className={classNames(s['tool-tag-list-cell'], {
                [s.active]: type === 'custom',
              })}
              onClick={() => onChangeAfterDiff('custom')}
            >
              {I18n.t('custom_tools')}
            </div>
          </div>
        </>
      ) : null}
    </div>
  );
};
