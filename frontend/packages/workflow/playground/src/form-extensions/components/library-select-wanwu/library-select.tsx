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

import { IconCozPlus, IconCozSetting } from '@coze-arch/coze-design/icons';
import { IconButton, Tooltip } from '@coze-arch/coze-design';
import { I18n } from '@coze-arch/i18n';
import { FieldEmpty } from '@/form';

import { type Library } from './types';
import { LibraryCard } from './library-card';

type DefaultLibraryRender = () => React.ReactNode;
interface RenderLibraryProps {
  defaultLibraryRender: DefaultLibraryRender;
  library: Library;
}
type RenderLibrary = (props: RenderLibraryProps) => React.ReactNode;

interface LibrarySelectProps {
  libraries?: Library[];
  readonly?: boolean;
  onEditLibrary?: (id: string) => void;
  onDeleteLibrary?: (id: string) => void;
  onAddLibrary?: () => void;
  onSettingLibrary?: () => void;
  onClickLibrary?: (id: string) => void;
  renderLibrary?: RenderLibrary;
  settingLibraryTooltip?: string;
  hideSettingButton?: boolean;
  emptyText?: string;
  hideAddButton?: boolean;
  addButtonTestID?: string;
  libraryCardTestID?: string;
}

export const LibrarySelect = ({
  libraries = [],
  readonly,
  onEditLibrary,
  onDeleteLibrary,
  onAddLibrary,
  onClickLibrary,
  renderLibrary,
  onSettingLibrary,
  settingLibraryTooltip,
  emptyText = '',
  hideSettingButton = true,
  hideAddButton = false,
  addButtonTestID = '',
  libraryCardTestID = '',
}: LibrarySelectProps) => (
  <div className="relative">
    {readonly || hideAddButton ? (
      <></>
    ) : (
      <div className="absolute right-[0] top-[-32px]">
        <IconButton
          color="highlight"
          onClick={onAddLibrary}
          theme="borderless"
          icon={<IconCozPlus />}
          size="small"
          data-testid={addButtonTestID}
        />
      </div>
    )}
    {readonly || hideSettingButton ? (
      <></>
    ) : (
      <div className="absolute right-[36px] top-[-32px]">
        <Tooltip content={settingLibraryTooltip || I18n.t('basic_setting')} autoAdjustOverflow>
          <IconButton
            color="highlight"
            onClick={onSettingLibrary}
            theme="borderless"
            icon={<IconCozSetting />}
            size="small"
            data-testid={`${libraryCardTestID}.setting`}
          />
        </Tooltip>
      </div>
    )}
    <div className="flex flex-col gap-[4px]">
      {libraries.length > 0 ? (
        libraries.map(library => {
          const isInvalid = library?.isInvalid;
          const defaultLibraryRender = () => (
            <LibraryCard
              isInvalid={isInvalid}
              readonly={readonly}
              key={library.id}
              library={library}
              showEditBtn={!library.external}
              onEdit={onEditLibrary}
              onDelete={onDeleteLibrary}
              onClick={id => {
                if (isInvalid) {
                  return;
                }
                onClickLibrary?.(id);
              }}
              testID={libraryCardTestID}
            />
          );

          if (renderLibrary) {
            return renderLibrary({ library, defaultLibraryRender });
          }

          return defaultLibraryRender();
        })
      ) : (
        <FieldEmpty text={emptyText} isEmpty={true} />
      )}
    </div>
  </div>
);
